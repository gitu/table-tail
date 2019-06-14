package tail

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
)

// Option to configure Tails
type Option func(c *Config) (*Config, error)

// Config for tails
type Config struct {
	db     *sql.DB
	table  string
	target chan string

	fields      string
	id          string
	format      string
	interval    time.Duration
	placeholder string

	template *template.Template
	ticker   *time.Ticker
	done     chan bool
	last     interface{}
	// TODO: add logger if nothing happens...
	lastChange time.Time
	started    bool
}

// Format for tail - default is {{.ID}}
func Format(format string) Option {
	return func(c *Config) (*Config, error) {
		c.format = format
		return c, nil
	}
}

// ID for tail - default is ID
func ID(id string) Option {
	return func(c *Config) (*Config, error) {
		c.id = id
		return c, nil
	}
}

// Fields for tail - default is ID
func Fields(fields string) Option {
	return func(c *Config) (*Config, error) {
		c.fields = fields
		return c, nil
	}
}

// Interval for tail - default is 200ms
func Interval(d time.Duration) Option {
	return func(c *Config) (*Config, error) {
		c.interval = d
		return c, nil
	}
}

// Placeholder ident to use for your driver, default is :
func Placeholder(p string) Option {
	return func(c *Config) (*Config, error) {
		c.placeholder = p
		return c, nil
	}
}

// Start creates and starts a new tail new Tail
func Start(db *sql.DB, table string, target chan string, opts ...Option) (*Config, error) {
	c := &Config{
		db:     db,
		table:  table,
		target: target,

		id:          "ID",
		fields:      "ID",
		format:      "{{.ID}}",
		interval:    200 * time.Millisecond,
		placeholder: ":",

		done:    make(chan bool, 1),
		started: false,
	}
	var err error

	for _, opt := range opts {
		c, err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.started {
		panic("already started")
	}
	c.started = true

	c.template, err = template.New("").Funcs(sprig.TxtFuncMap()).Parse(c.format)
	if err != nil {
		return nil, err
	}

	err = c.fetchInitial()
	if err != nil {
		return nil, err
	}

	c.ticker = time.NewTicker(c.interval)
	go c.tail()
	return c, nil
}

// Stop the tail and close the channel
func (c *Config) Stop() {
	c.ticker.Stop()
	close(c.done)
	close(c.target)
}

func (c *Config) tail() {
	for {
		select {
		case <-c.ticker.C:
			// TODO: add possibility for error channel
			err := c.fetch()
			if err != nil {
				c.target <- err.Error()
			}
		case <-c.done:
			return
		}
	}
}

func (c *Config) fetchInitial() error {
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC FETCH FIRST ROW ONLY", c.id, c.table, c.id)
	row := c.db.QueryRow(query)
	if row != nil {
		var initialLast interface{}
		err := row.Scan(&initialLast)
		if err != nil && err != sql.ErrNoRows {
			return errors.Wrap(err, "while executing: "+query)
		}
		c.last = initialLast
	}
	return nil
}

// tries to fetch and set last row
func (c *Config) fetch() error {
	query := fmt.Sprintf("SELECT %s FROM %s where %s > %s1 ORDER BY %s ASC", c.fields, c.table, c.id, c.placeholder, c.id)
	params := []interface{}{c.last}
	if c.last == nil {
		query = fmt.Sprintf("SELECT %s FROM %s ORDER BY %s ASC FETCH FIRST ROW ONLY", c.fields, c.table, c.id)
		params = []interface{}{}
	}

	rows, err := c.db.Query(query, params...)
	if err != nil {
		return errors.Wrap(err, "while executing: "+query)
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	var newLast interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		var tpl bytes.Buffer
		if err := c.template.Execute(&tpl, m); err != nil {
			return err
		}
		last, found := m[c.id]
		if !found {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			return errors.Errorf("id: [%v] not found in fields: %+v", c.id, keys)
		}
		newLast = last
		c.target <- tpl.String()
	}
	if newLast != nil {
		c.last = newLast
	}
	return nil
}

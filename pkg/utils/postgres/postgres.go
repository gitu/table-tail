package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gitu/table-tail/pkg/utils"
	"time"
)

func init() {
	utils.Register("postgres", newUtil())
}

type util struct {
}

func newUtil() utils.TailUtil {
	u := util{}
	return &u
}

// Returns Info about connection
// example: connected to HOST/INSTANCE (VERSION) -- [139.584µs]
func (*util) ConnectionInfo(db *sql.DB) (string, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var version, database string
	qry := "SELECT current_version(), database()"
	err := db.QueryRowContext(ctx, qry).Scan(&version, &database)
	if err != nil {
		return "", err
	}
	stop := time.Now()
	return fmt.Sprintf("connected to %s (%s) -- [%s]", database, version, stop.Sub(start).String()), nil
}

func (*util) PlaceHolderMarker() string {
	return "&"
}

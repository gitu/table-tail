package tail

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TailTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *TailTestSuite) SetupSuite() {
	var err error
	datasource := "postgres://root@localhost:5432/test?sslmode=disable"
	s.db, err = sql.Open("postgres", datasource)
	if err != nil {
		s.FailNow("could not connect to db", err, datasource)
	}
}

func (s *TailTestSuite) SetupTest() {
	_, err := s.db.Exec("DROP TABLE IF EXISTS sample_log_table  ")
	if err != nil {
		s.FailNow("could not create table", err)
	}
	_, err = s.db.Exec("CREATE TABLE sample_log_table (id serial PRIMARY KEY,  msg VARCHAR (255) NOT NULL)")
	if err != nil {
		s.FailNow("could not create table", err)
	}
}

func (s *TailTestSuite) TestSimple() {
	cons := make(chan string, 1)
	c, err := Start(s.db, "sample_log_table", cons, Format("now"))
	if err != nil {
		s.FailNow("could not create table", err)
	}

	defer c.Stop()

	select {
	case msg := <-cons:
		s.Fail("expected no channel msg, but got " + msg)
	default:
	}

	s.addMsg("test")

	select {
	case msg := <-cons:
		fmt.Println("Got Response: ", msg)
	case <-time.After(1 * time.Second):
		s.Fail("ran into timeout ")
	}

}

func (s *TailTestSuite) TestMsgContent() {
	cons := make(chan string, 1)

	c, err := Start(s.db, "sample_log_table", cons, ID("id"), Fields("id,msg"), Format("{{.msg}}"), Placeholder("&"))
	if err != nil {
		s.FailNow("could not select table", err)
	}
	defer c.Stop()

	select {
	case msg := <-cons:
		s.Fail("expected no channel msg, but got " + msg)
	default:
	}

	s.addMsg("test ASDF")

	for {
		select {
		case msg := <-cons:
			s.Equal("test ASDF", msg)
			s.Equal(int64(1), c.last)
			return
		case <-time.After(3 * time.Second):
			s.Fail("ran into timeout ")
		}
	}
}

func (s *TailTestSuite) TestMultipleMessages() {
	cons := make(chan string, 1)
	s.addMsg("test XXXX")

	c, err := Start(s.db, "sample_log_table", cons, ID("id"), Fields("id,msg"), Format("{{.msg}}"), Interval(100*time.Millisecond), Placeholder("$"))
	if err != nil {
		s.FailNow("could not create table", err)
	}
	defer c.Stop()

	select {
	case msg := <-cons:
		s.Fail("expected no channel msg, but got " + msg)
	default:
	}

	s.Equal(int64(1), c.last)

	s.addMsg("test 1")
	s.addMsg("test 2")
	s.addMsg("test 3")
	s.addMsg("test 4")

	select {
	case msg := <-cons:
		s.Equal("test 1", msg)
	case <-time.After(1 * time.Second):
		s.Fail("ran into timeout ")
	}
	select {
	case msg := <-cons:
		s.Equal("test 2", msg)
	case <-time.After(1 * time.Second):
		s.Fail("ran into timeout ")
	}
	select {
	case msg := <-cons:
		s.Equal("test 3", msg)
	case <-time.After(1 * time.Second):
		s.Fail("ran into timeout ")
	}
	select {
	case msg := <-cons:
		s.Equal("test 4", msg)
	case <-time.After(1 * time.Second):
		s.Fail("ran into timeout ")
	}

}

func (s *TailTestSuite) addMsg(msg string) {
	res, err := s.db.Exec("insert into sample_log_table (msg) values ($1)", msg)
	if err != nil {
		s.FailNow("could insert message", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		s.FailNow("could insert message", err)
	}
	if rows != 1 {
		s.FailNow("could insert message", err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TailTestSuite))
}

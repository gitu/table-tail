package oracle

import (
	"github.com/gitu/table-tail/pkg/utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestConnectionInfo(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	util, err := utils.Get("postgres")
	if err != nil {
		t.Errorf("error was not expected while getting util: %s", err)
	}

	rows := sqlmock.NewRows([]string{"current_version", "database"}).
		AddRow("VERSION", "DATABASE")

	mock.ExpectQuery("^SELECT current_version\\(\\), database\\(\\)").WillReturnRows(rows)

	info, err := util.ConnectionInfo(db)

	if err != nil {
		t.Errorf("error was not expected while selecting infos: %s", err)
	}

	assert.Regexp(regexp.MustCompile("connected to DATABASE \\(VERSION\\) -- .+"), info)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

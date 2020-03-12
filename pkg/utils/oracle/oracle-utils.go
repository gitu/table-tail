package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gitu/table-tail/pkg/utils"
	_ "github.com/godror/godror"
)

func init() {
	utils.Register("godror", NewOracleUtil())
}

type util struct {
}

// Returns new util for an oracle connection
func NewOracleUtil() utils.TailUtil {
	u := util{}
	return &u
}

// Returns Info about connection
// example: connected to HOST/INSTANCE (VERSION) -- [139.584µs]
func (*util) ConnectionInfo(db *sql.DB) (string, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var version, instanceName, hostName string
	qry := "SELECT VERSION, INSTANCE_NAME, HOST_NAME FROM V$INSTANCE"
	err := db.QueryRowContext(ctx, qry).Scan(&version, &instanceName, &hostName)
	if err != nil {
		return "", err
	}
	stop := time.Now()
	return fmt.Sprintf("connected to %s/%s (%s) -- [%s]", hostName, instanceName, version, stop.Sub(start).String()), nil
}

func (*util) PlaceHolderMarker() string {
	return ":"
}

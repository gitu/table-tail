package utils

import (
	"database/sql"
	"github.com/pkg/errors"
	"sync"
)

// TailUtil provides an interface for additional information specific to the driver
type TailUtil interface {
	// Returns Info about connection
	// example: connected to HOST/INSTANCE (VERSION) -- [139.584µs]
	ConnectionInfo(db *sql.DB) (string, error)
}

var (
	utilsMutex sync.RWMutex
	utils      = make(map[string]TailUtil)
)

// Register new Tail Util - name should be the same as the corresponding driver
func Register(name string, util TailUtil) {
	utilsMutex.Lock()
	defer utilsMutex.Unlock()
	if util == nil {
		panic("tail-utils: Register util is nil")
	}
	if _, dup := utils[name]; dup {
		panic("tail-utils: Register called twice for util " + name)
	}
	utils[name] = util
}

// Get Tail Util
func Get(name string) (TailUtil, error) {
	util, found := utils[name]
	if !found {
		return nil, errors.Errorf("tail-utils: util not found %s", name)
	}
	return util, nil
}

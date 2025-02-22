package redis

import (
	"io"

	red "github.com/go-redis/redis"
	"github.com/melonwool/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(server, pass string, db int) (*red.Client, error) {
	val, err := clientManager.GetResource(server, func() (io.Closer, error) {
		store := red.NewClient(&red.Options{
			Addr:         server,
			Password:     pass,
			DB:           db,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
		})
		store.WrapProcess(process)
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}

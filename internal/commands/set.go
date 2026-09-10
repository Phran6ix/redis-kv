package commands

import (
	"errors"

	"github.com/Phran6ix/redis-kv/internal/store"
)

func set(key string, value any) error {
	done, msg := store.GlobalStore.Set(key, value)
	if !done {
		return errors.New(msg)
	}

	return nil
}

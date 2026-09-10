package commands

import (
	"github.com/Phran6ix/redis-kv/internal/store"
)

func get(key string) any {
	exists, value := store.GlobalStore.Get(key)
	if !exists {
		return nil
	}

	return value
}

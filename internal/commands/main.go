package commands

import (
	"errors"
	"strings"

	"github.com/Phran6ix/redis-kv/internal/resp"
)

func ProcessCommand(redisValue []resp.RedisValue) (bool, any, error) {
	var returnVal any = nil

	command, ok := redisValue[0].(resp.BulkString)

	if !ok {
		return false, "", errors.New("command should be a simple string")
	}

	switch strings.ToUpper(command.Value) {
	case "SET":
		{
			if len(redisValue) < 3 {
				return false, "", errors.New("missing parameter in SET argument")
			}
			var key string
			var value any
			k, ok := redisValue[1].(resp.BulkString)
			if !ok {
				return false, "", errors.New("missing key in SET argument")
			}

			key = k.Value

			if key == "" {
				return false, "", errors.New("key cannot be empty.")
			}

			v := redisValue[2]

			switch v := v.(type) {
			case resp.Error:
				value = v.Message
			case resp.Boolean:
				value = v.Value
			case resp.SimpleString:
				value = v.Value
			case resp.BulkString:
				value = v.Value
			case resp.Integer:
				value = v.Value
			case resp.Double:
				value = v.Value
			default:
				value = nil
			}

			err := set(key, value)
			if err != nil {
				return false, "", err
			}
		}
	case "GET":
		{
			if len(redisValue) < 2 {
				return false, "", errors.New("key is missing in GET command")
			}

			key := redisValue[1].(resp.BulkString)

			if key.Value == "" {
				return false, "", errors.New("key is missing in GET command")
			}
			value := get(key.Value)

			if value == nil {
				return false, "", nil
			}

			returnVal = value
		}
	}

	return true, returnVal, nil
}

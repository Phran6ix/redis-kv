package resp

import (
	"errors"
	"fmt"
)

func Serialize(redisValue RedisValue) (string, error) {
	var result string

	switch r := redisValue.(type) {
	case SimpleString:
		result = fmt.Sprintf("+%s\r\n", r.Value)
	case BulkString:
		stringlen := len(r.Value)
		result = fmt.Sprintf("$%d\r\n%s\r\n", stringlen, r.Value)
	case Integer:
		result = fmt.Sprintf(":%d\r\n", r.Value)
	case Null:
		result = fmt.Sprint("-1\r\n")
	case Error:
		result = fmt.Sprintf("-%s\r\n", r.Message)
	case Double:
		result = fmt.Sprintf(",%f\r\n", r.Value)
	case Boolean:
		var booleanString string = "f"

		if r.Value {
			booleanString = "t"
		}
		result = fmt.Sprintf("#%s", booleanString)
	default:
		return "", errors.New("Unsupported interface")
	}

	return result, nil
}

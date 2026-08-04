package resp

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var CRLF = []byte("\r\n")

func Parse(b []byte) ([]RedisValue, error) {
	//  *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	idx := 0

	if !bytes.HasPrefix(b, []byte("*")) {
		log.Println("Expecting an array")
		return nil, errors.New("-Error expecting an array as first byte")
	}

	// Move the index to the length of the array
	idx += 1
	length := int(b[idx])
	fmt.Println("this is the length of the array => %d", length)

	// move the index to the start of the array
	idx += len(CRLF)

	redisValue := make([]any, length)

	for {
		if idx == length {
			break
		}

		end := bytes.Index(b, CRLF)
		if end == -1 {
			return nil, errors.New("Incomplete command")
		}

		data := b[:end]

		// handle strings

		if bytes.HasPrefix(data, []byte("+")) {
		} else if bytes.HasPrefix(data, []byte(":")) {
		} else if bytes.HasPrefix(data, []byte("-")) {
		} else if bytes.HasPrefix(data, []byte("*")) {
		} else if bytes.HasPrefix(data, []byte("#")) {
		} else if bytes.HasPrefix(data, []byte(",")) {
		} else {
			log.Fatal("Unsupported command")
			return nil, errors.New("-Error Unsupported RESP type")
		}

		// idx += end +
	}

	return
}

func parseString(bs []byte) RedisValue {
	return SimpleString{Value: string(bs)}
}

func parseInt(bi []byte) (RedisValue, error) {
	i, err := strconv.Atoi(string(bi))
	if err != nil {
		return Integer{}, errors.New("-Error Invalid integer value")
	}
	return Integer{Value: int64(i)}, nil
}

func parseError(be []byte) RedisValue {
	return Error{Message: string(be)}
}

func parseDouble(bd []byte) (RedisValue, error) {
	f, err := strconv.ParseFloat(string(bd), 64)
	if err != nil {
		return Double{}, errors.New("-Error invalid double")
	}

	return Double{Value: float64(f)}, nil
}

func parseBoolean(bb []byte) (RedisValue, error) {
	var b bool
	switch bb[0] {
	case 't':
		b = true
	case 'f':
		b = false
	default:
		return Boolean{}, errors.New("-Error Invalid boolean payload")
	}

	return Boolean{Value: b}, nil
}

type RedisValue interface {
	isRedisValue([]byte)
}

type SimpleString struct {
	Value string
}

type Error struct {
	Message string
}

type Boolean struct {
	Value bool
}

type Integer struct {
	Value int64
}

type Double struct {
	Value float64
}

// type Array struct {
// 	elements []RedisValue
// }

func (ss SimpleString) isRedisValue(bb []byte) {}
func (e Error) isRedisValue(bb []byte)         {}
func (b Boolean) isRedisValue(bb []byte)       {}
func (i Integer) isRedisValue(bb []byte)       {}
func (d Double) isRedisValue(bb []byte)        {}

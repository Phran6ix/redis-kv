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
	done := 0

	redisValue := make([]RedisValue, length)

	for {
		if len(b) < idx+1 {
			break
		}
		if done == length {
			break
		}

		end := bytes.Index(b, CRLF)
		if end == -1 {
			return nil, errors.New("Incomplete command")
		}

		data := b[idx : idx+end]
		fmt.Println("DATA => %s", data)

		// handle strings

		if bytes.HasPrefix(data, []byte("+")) {
			// +OK\r\n
			redisValue = append(redisValue, parseString(data[1:]))
		} else if bytes.HasPrefix(data, []byte(":")) {
			i, err := parseInt(data[1:])
			if err != nil {
				return nil, err
			}
			redisValue = append(redisValue, i)
		} else if bytes.HasPrefix(data, []byte("-")) {
			// -Error message\r\n
			redisValue = append(redisValue, parseError(data[1:]))
		} else if bytes.HasPrefix(data, []byte("#")) {
			bv, err := parseBoolean(data[1:])
			if err != nil {
				return nil, err
			}
			redisValue = append(redisValue, bv)
		} else if bytes.HasPrefix(data, []byte(",")) {
			db, err := parseDouble(data[1:])
			if err != nil {
				return nil, err
			}
			redisValue = append(redisValue, db)
		} else if bytes.HasPrefix(data, []byte("$")) {
			// $<length>\r\n<data>\r\n
			bs, read, err := parseBulkString(data[1:])
			if err != nil {
				return nil, err
			}
			idx += read
			redisValue = append(redisValue, bs)
		} else {
			log.Fatal("Unsupported command")
			return nil, errors.New("-Error Unsupported RESP type")
		}

		// idx += 1
		idx += end + len(CRLF)
		done++
	}

	// incomplete data
	if done != length {
		return nil, errors.New("-Error Incomplete Array Elements")
	}

	if len(b[idx:]) > 0 {
		return nil, errors.New("-Error More than expected elements in the array")
	}

	return redisValue, nil
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

func parseBulkString(bb []byte) (val RedisValue, read int, error error) {
	read = 0

	lengthIdx := bytes.Index(bb, CRLF)
	if lengthIdx == -1 {
		return nil, 0, errors.New("-Error Invalid Data Type in bulk string ")
	}

	read += lengthIdx
	lengthStr := string(bb[:lengthIdx])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, read, errors.New("-Error expected base 10 bulk string length")
	}

	rawstrIdx := bytes.Index(bb[lengthIdx+len(CRLF):], CRLF)

	if rawstrIdx == -1 {
		return nil, read, errors.New("-Error Incomplete Bulk String")
	}

	rawStr := bb[length+len(CRLF) : rawstrIdx]
	println("RAW STR %b", rawStr)
	if len(rawStr) != length {
		return nil, read, errors.New("-Error Length of bulk string mismatch")
	}

	read += rawstrIdx + len(CRLF)

	return BulkString{Value: string(rawStr)}, read, nil
}

type RedisValue interface {
	isRedisValue([]byte)
}

type SimpleString struct {
	Value string
}

type BulkString struct {
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
func (bs BulkString) isRedisValue(bb []byte)   {}

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
	// link:  https://redis.io/docs/latest/develop/reference/protocol-spec/#resp-protocol-description
	//  *112\r\n$5\r\nhello\r\n$5\r\nworld\r\n

	idx := 0

	if !bytes.HasPrefix(b, []byte("*")) {
		log.Println("Expecting an array")
		return nil, errors.New("-Error expecting an array as first byte")
	}

	// Move the index to the length of the array
	idx += len([]byte("*"))

	lengthidx := bytes.Index(b[idx:], CRLF)
	if lengthidx == -1 {
		return nil, errors.New("-Error Missing protocol terminator")
	}
	lengthBytes := b[idx : lengthidx+idx]

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return nil, errors.New("-Error Invalid integer provided for array length")
	}

	fmt.Printf("The raw bytes => %q bytes\n", lengthBytes)

	// move the index to the start of the array

	idx += len(lengthBytes) + len(CRLF)
	done := 0

	redisValue := make([]RedisValue, 0, length)

	for done < length {
		if len(b) < idx+1 {
			break
		}

		end := bytes.Index(b[idx:], CRLF)
		if end == -1 {
			return nil, errors.New("-Error Incomplete command")
		}

		data := b[idx : idx+end]
		fmt.Printf("DATA => %q\n", data)

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
			// fmt.Printf("The total byte data => %q \n ", data)
			bs, read, err := parseBulkString(b[idx:])
			if err != nil {
				return nil, err
			}
			idx += read
			done++

			redisValue = append(redisValue, bs)

			continue
		} else {
			fmt.Printf("prefix = %q ", data)
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
	return Integer{Value: int(i)}, nil
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
	fmt.Printf("%q \n ", bb)
	if lengthIdx == -1 {
		return nil, 0, errors.New("-Error Invalid Data Type in bulk string ")
	}

	read += lengthIdx + len(CRLF)

	lengthStr := string(bb[1:lengthIdx])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, read, errors.New("-Error expected base 10 bulk string length")
	}
	// Null Bulk String
	// $-1\r\n
	if length == -1 {
		return Null{}, read, nil
	}

	rest := bb[lengthIdx+len(CRLF):]
	// fmt.Printf("The lengthstr %q\n", rest)
	rawstrIdx := bytes.Index(rest, CRLF)

	if rawstrIdx == -1 {
		return nil, read, errors.New("-Error Incomplete Bulk String")
	}

	rawStr := rest[:rawstrIdx]
	fmt.Printf("\n RAW STR %q \n", rawStr)
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
	Value int
}

type Double struct {
	Value float64
}

type Null struct{}

// type Array struct {
// 	elements []RedisValue
// }

func (ss SimpleString) isRedisValue(bb []byte) {}
func (e Error) isRedisValue(bb []byte)         {}
func (n Null) isRedisValue(bb []byte)          {}
func (b Boolean) isRedisValue(bb []byte)       {}
func (i Integer) isRedisValue(bb []byte)       {}
func (d Double) isRedisValue(bb []byte)        {}
func (bs BulkString) isRedisValue(bb []byte)   {}

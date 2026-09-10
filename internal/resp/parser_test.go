package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	var testCase []byte
	assert := assert.New(t)
	require := require.New(t)

	testCase = []byte("+OK\r\n")

	redisValue, err := Parse(testCase)
	require.Error(err, "-Error expecting an array as first byte")
	require.Nil(redisValue)

	// // Test with simple strings
	testCase = []byte("*3\r\n+Str1\r\n+Str2\r\n+Str3\r\n")
	redisValue, err = Parse(testCase)

	require.NoError(err)
	require.NotNil(redisValue)
	assert.Equal(len(redisValue), 3)

	testCase = []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	redisValue, err = Parse(testCase)
	require.NoError(err)
	require.NotNil(redisValue)
	assert.Equal(2, len(redisValue))

	// Test: Mixed Types (Bulk String and Simple String together)
	testCase = []byte("*2\r\n$4\r\nPING\r\n+PONG\r\n")
	redisValue, err = Parse(testCase)
	require.NoError(err)
	require.NotNil(redisValue)
	assert.Equal(2, len(redisValue))

	// Test: Empty Array (Valid in RESP)
	testCase = []byte("*0\r\n")
	redisValue, err = Parse(testCase)
	require.NoError(err)
	assert.Equal(0, len(redisValue))

	// Test: Null Bulk String inside an array (e.g., a missing key response)
	testCase = []byte("*3\r\n$3\r\nGET\r\n$6\r\nmy_key\r\n$-1\r\n")
	redisValue, err = Parse(testCase)
	require.NoError(err)
	require.NotNil(redisValue)
	assert.Equal(3, len(redisValue))
}

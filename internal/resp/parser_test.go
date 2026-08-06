package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	var testCase []byte
	assert := assert.New(t)

	testCase = []byte("+OK\r\n")

	redisValue, err := Parse(testCase)
	assert.EqualError(err, "-Error expecting an array as first byte")
	assert.Equal(redisValue, nil)

	// Test with simple strings
	testCase = []byte("*3\r\n+Str1\r\n+Str2\r\n+Str3\r\n")
	redisValue, err = Parse(testCase)

	require.NoError(t, err)
	require.NotNil(t, redisValue)
	assert.Equal(len(redisValue), 3)

	// Test with inconsistent length
	testCase = []byte("*3\r\n+Str1\r\n+Str2\r\n")
}

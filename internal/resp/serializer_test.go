package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespSerializer(t *testing.T) {
	simplestring := SimpleString{
		Value: "simple",
	}

	simpleStringSerial, err := Serialize(simplestring)
	require.Nil(t, err)
	require.NotNil(t, simpleStringSerial)
	assert.Equal(t, "+simple\r\n", simpleStringSerial)

	bulkstring := BulkString{
		Value: "This is my town",
	}

	bulkstringSerial, err := Serialize(bulkstring)
	require.Nil(t, err)
	require.NotNil(t, bulkstringSerial)
	assert.Equal(t, "$15\r\nThis is my town\r\n", bulkstringSerial)

	integer := Integer{
		Value: 24,
	}

	integerSerial, err := Serialize(integer)
	require.NotNil(t, integerSerial)
	require.Nil(t, err)
	assert.Equal(t, ":24\r\n", integerSerial)

	boolean := Boolean{
		Value: true,
	}
	booleanSerial, err := Serialize(boolean)

	require.Nil(t, err)
	require.NotNil(t, booleanSerial)
	assert.Equal(t, "#t", booleanSerial)
}

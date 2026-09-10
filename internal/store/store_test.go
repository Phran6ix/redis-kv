package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetFunction(t *testing.T) {
	testStore := Store{
		store: make(map[string]any),
	}

	testStore.mutex.Lock()
	assert.Empty(t, testStore.store)
	testStore.mutex.Unlock()

	testStore.Set("new_key", "value")

	testStore.mutex.Lock()
	expected, exists := testStore.store["new_key"]
	testStore.mutex.Unlock()

	require.NotNil(t, expected)
	assert.True(t, exists)
	assert.Equal(t, "value", expected)
}

func TestGetFunction(t *testing.T) {
	testStore := Store{
		store: make(map[string]any),
	}

	testStore.mutex.Lock()
	testStore.store["key"] = "value"
	testStore.mutex.Unlock()

	found, value := testStore.Get("key")

	assert.True(t, found)
	assert.Equal(t, "value", value)
}

func TestBothTogetherConcurrently(t *testing.T) {
	testStore := Store{
		store: make(map[string]any),
	}

	var wg sync.WaitGroup

	require.Empty(t, testStore.store)

	iterations := 10
	values := make([]string, iterations)
	for i := range values {
		values[i] = fmt.Sprintf("value-%d", i)
	}

	assert.Equal(t, len(values), iterations)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i, value := range values {
			key := fmt.Sprintf("key-%d", i)
			testStore.Set(key, value)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for i, value := range values {
			key := fmt.Sprintf("key-%d", i+10)
			testStore.Set(key, value)
		}
	}()
	wg.Wait()

	testStore.mutex.Lock()
	store_length := len(testStore.store)
	testStore.mutex.Unlock()

	assert.Equal(t, store_length, iterations*2)
}

package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Выталкивание при переполнении
		c := NewCache(3)
		c.Set("one", 1)
		c.Set("two", 2)
		c.Set("three", 3)
		c.Set("four", 4)

		val, exists := c.Get("one")
		require.Equal(t, nil, val)
		require.Equal(t, false, exists)

		c.Clear()

		// Выталкивание самого невостребованного
		c.Set("one", 1)   // [1]
		c.Set("two", 2)   // [2, 1]
		c.Set("three", 3) // [3,2,1]

		c.Set("one", 11) // [11,3,2]
		c.Get("two")     // [2,11,3]
		c.Get("three")   // [3,2,11]
		c.Set("four", 4) // [4,3,2]

		val, exists = c.Get("one")
		require.Equal(t, nil, val)
		require.Equal(t, false, exists)

		val, exists = c.Get("two")
		require.Equal(t, 2, val)
		require.Equal(t, true, exists)

		val, exists = c.Get("three")
		require.Equal(t, 3, val)
		require.Equal(t, true, exists)

		val, exists = c.Get("four")
		require.Equal(t, 4, val)
		require.Equal(t, true, exists)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

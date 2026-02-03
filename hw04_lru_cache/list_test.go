package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("adds one item in front of empty list", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 1, l.Len())
	})

	t.Run("adds one item in front of non-empty list", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)
		l.PushFront(20)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())
		require.NotEqual(t, l.Front(), l.Back())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})

	t.Run("adds one item in back of empty list", func(t *testing.T) {
		l := NewList()
		l.PushBack(10)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 1, l.Len())
	})

	t.Run("adds one item in back of non-empty list", func(t *testing.T) {
		l := NewList()
		l.PushBack(10)
		l.PushBack(20)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())
		require.NotEqual(t, l.Front(), l.Back())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 20, l.Back().Value)
	})

	t.Run("removes item from single item list", func(t *testing.T) {
		l := NewList()
		i := l.PushBack(10)
		l.Remove(i)

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())
	})

	t.Run("removes item from middle of the list", func(t *testing.T) {
		l := NewList()
		f := l.PushBack(10)
		m := l.PushBack(20)
		b := l.PushBack(30)
		l.Remove(m)

		require.Equal(t, f, l.Front())
		require.Equal(t, b, l.Back())
		require.Equal(t, 2, l.Len())
	})

	t.Run("removes item from front of the list", func(t *testing.T) {
		l := NewList()
		f := l.PushBack(10)
		b := l.PushBack(20)
		l.Remove(f)

		require.Equal(t, b, l.Front())
		require.Equal(t, b, l.Back())
		require.Equal(t, 1, l.Len())
	})

	t.Run("removes item from back of the list", func(t *testing.T) {
		l := NewList()
		f := l.PushBack(10)
		b := l.PushBack(20)
		l.Remove(b)

		require.Equal(t, f, l.Front())
		require.Equal(t, f, l.Back())
		require.Equal(t, 1, l.Len())
	})

	t.Run("removes sole item from the list", func(t *testing.T) {
		l := NewList()
		i := l.PushBack(10)
		l.Remove(i)

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		require.Equal(t, 0, l.Len())
	})

	t.Run("correctly moves same element multiple times", func(t *testing.T) {
		l := NewList()
		b := l.PushBack(10)
		i := l.PushBack(20)
		l.MoveToFront(i)
		l.MoveToFront(i)
		l.MoveToFront(i)

		require.Equal(t, i, l.Front())
		require.Equal(t, b, l.Back())
		require.Equal(t, 2, l.Len())
		require.Nil(t, i.Prev)
		require.Equal(t, b, i.Next)
		require.Nil(t, b.Next)
		require.Equal(t, i, b.Prev)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

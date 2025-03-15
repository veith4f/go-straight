package pkg

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestFib(t *testing.T) {
  t.Run("Fib 0", func(t *testing.T) {
    got := Fib(0)
    want := 1

    assert.Equal(t, got, want, "they should be equal")
  })
  t.Run("Fib 1", func(t *testing.T) {
    got := Fib(1)
    want := 1

    assert.Equal(t, got, want, "they should be equal")
  })
  t.Run("Fib 2", func(t *testing.T) {
    got := Fib(2)
    want := 2

    assert.Equal(t, got, want, "they should be equal")
  })
  t.Run("Fib 3", func(t *testing.T) {
    got := Fib(3)
    want := 3

    assert.Equal(t, got, want, "they should be equal")
  })
  t.Run("Fib 4", func(t *testing.T) {
    got := Fib(4)
    want := 5

    assert.Equal(t, got, want, "they should be equal")
  })
  t.Run("Fib 10", func(t *testing.T) {
    got := Fib(10)
    want := 89

    assert.Equal(t, got, want, "they should be equal")
  })
}

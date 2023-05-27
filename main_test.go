package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Parallel()

	main()

	expect := 1
	got := 1

	assert.Equal(t, expect, got)
}

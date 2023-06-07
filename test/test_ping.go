package test

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func testPing(t *testing.T) {
	var res interface{}
	code := Get("/ping", "", map[string][]string{}, &res)
	assert.Equal(t, code, 200)
}

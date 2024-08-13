package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("Ahmed", "Shaban", "hunter1234")
	assert.NotNil(t, err)

	fmt.Printf("Account is %v", acc)
}

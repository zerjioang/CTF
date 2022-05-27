package tinyeth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		acc, err := CreateAccount()
		assert.NoError(t, err)
		assert.NotNil(t, acc)
	})
}

package tinyeth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	localGanachePeerURL = "http://127.0.0.1:8545"
)

func TestClient(t *testing.T) {
	t.Run("connect-with-empty-url", func(t *testing.T) {
		var c Client
		assert.Error(t, c.Connect(""))
	})
	t.Run("connect", func(t *testing.T) {
		var c Client
		assert.NoError(t, c.Connect(localGanachePeerURL))
	})
	t.Run("fetch-latest-block-number", func(t *testing.T) {
		var c Client
		assert.NoError(t, c.Connect(localGanachePeerURL))
		_, b := c.LatestBlock()
		assert.NotEmpty(t, b)
		t.Log("latest block is:", b)
	})
}

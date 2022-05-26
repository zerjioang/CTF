package tinyeth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleHex(t *testing.T) {
	t.Run("simple-hex", func(t *testing.T) {
		assert.Equal(t, SimpleHex("0x4510108d3033173b84919F4941695d4190A0F40F"), "4510108d3033173b84919F4941695d4190A0F40F")
	})
	t.Run("simple-hex-2", func(t *testing.T) {
		addr := BuildAddress("0x4510108d3033173b84919F4941695d4190A0F40F")
		assert.NotNil(t, addr)
		assert.Equal(t, "0x4510108d3033173b84919F4941695d4190A0F40F", addr.Hex())
	})
}

func BenchmarkSimpleHex(b *testing.B) {
	b.Run("simple-hex", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(1)
		b.ResetTimer()
		var addr string
		for i := 0; i < b.N; i++ {
			addr = SimpleHex("0x4510108d3033173b84919F4941695d4190A0F40F")
		}
		if addr != "" {

		}
	})
}

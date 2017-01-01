package restcache

import (
	//	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRouter(t *testing.T) {
	assert := assert.New(t)
	bench := testing.Benchmark(func(b *testing.B) {
		router := NewRouter()
		assert.NotNil(router)
	})
}

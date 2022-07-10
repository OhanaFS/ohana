package kv_test

import (
	"github.com/OhanaFS/ohana/service/kv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	assert := assert.New(t)
	m := kv.NewMemoryKV()
	assert.NotNil(m)

	_, err := m.Get(nil, "key")
	assert.NotNil(err)

	err = m.Set(nil, "foo", "bar", time.Hour)
	assert.Nil(err)

	v, err := m.Get(nil, "foo")
	assert.Nil(err)
	assert.Equal("bar", v)

	err = m.Delete(nil, "foo")
	assert.Nil(err)

	_, err = m.Get(nil, "foo")
	assert.NotNil(err)
}

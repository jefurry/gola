package gluabit32_test

import (
	"testing"

	"github.com/BixData/gluabit32"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func Test_bxor(t *testing.T) {
	assert := assert.New(t)
	L := lua.NewState()
	defer L.Close()
	gluabit32.Preload(L)

	assert.NoError(L.DoString("return require 'bit32'.bxor(5, 3)"))
	assert.Equal(6, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.bxor(3, 6)"))
	assert.Equal(5, L.ToInt(-1))
}

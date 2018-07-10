package gluabit32_test

import (
	"testing"

	"github.com/BixData/gluabit32"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func Test_rshift(t *testing.T) {
	assert := assert.New(t)
	L := lua.NewState()
	defer L.Close()
	gluabit32.Preload(L)

	// sample data at https://homerl.github.io/2016/03/29/golang-bitwise-operators/
	assert.NoError(L.DoString("return require 'bit32'.rshift(5, 1)"))
	assert.Equal(2, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.rshift(0xffffffff, 0)"))
	assert.Equal(4294967295, L.ToInt(-1))
}

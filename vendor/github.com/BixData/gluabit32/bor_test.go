package gluabit32_test

import (
	"testing"

	"github.com/BixData/gluabit32"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func Test_bor(t *testing.T) {
	assert := assert.New(t)
	L := lua.NewState()
	defer L.Close()
	gluabit32.Preload(L)

	// sample data at https://homerl.github.io/2016/03/29/golang-bitwise-operators/
	assert.NoError(L.DoString("return require 'bit32'.bor(1, 5)"))
	assert.Equal(5, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.bor(10, 12)"))
	assert.Equal(14, L.ToInt(-1))
}

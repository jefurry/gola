package gluabit32_test

import (
	"testing"

	"github.com/BixData/gluabit32"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func Test_lshift(t *testing.T) {
	assert := assert.New(t)
	L := lua.NewState()
	defer L.Close()
	gluabit32.Preload(L)

	assert.NoError(L.DoString("return require 'bit32'.lshift(0, 1)"))
	assert.Equal(0, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.lshift(1, 1)"))
	assert.Equal(2, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(1, 2)"))
	assert.Equal(4, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(1, 3)"))
	assert.Equal(8, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(1, 4)"))
	assert.Equal(16, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.lshift(-1, 1)"))
	assert.Equal(-2, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(-1, 2)"))
	assert.Equal(-4, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(-1, 3)"))
	assert.Equal(-8, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(-1, 4)"))
	assert.Equal(-16, L.ToInt(-1))

	assert.NoError(L.DoString("return require 'bit32'.lshift(3, 1)"))
	assert.Equal(6, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(3, 2)"))
	assert.Equal(12, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(3, 3)"))
	assert.Equal(24, L.ToInt(-1))
	assert.NoError(L.DoString("return require 'bit32'.lshift(3, 4)"))
	assert.Equal(48, L.ToInt(-1))

}

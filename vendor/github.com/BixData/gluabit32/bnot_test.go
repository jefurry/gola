package gluabit32_test

import (
	"testing"

	"github.com/BixData/gluabit32"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
)

func Test_bnot(t *testing.T) {
	assert := assert.New(t)
	L := lua.NewState()
	defer L.Close()
	gluabit32.Preload(L)

	assert.NoError(L.DoString("return require 'bit32'.bnot(0)"))
	assert.Equal(-1, L.ToInt(-1))
}

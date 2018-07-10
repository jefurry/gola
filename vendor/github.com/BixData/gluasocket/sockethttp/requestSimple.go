package gluasocket_sockethttp

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/gopher-lua"
)

func requestSimpleFn(L *lua.LState) int {
	httpClient := http.Client{Timeout: time.Second * 15}
	url := L.ToString(1)

	var res *http.Response
	var err error
	if L.Get(2).Type() == lua.LTNil {
		res, err = httpClient.Get(url)
	} else {
		body := L.ToString(2)
		res, err = httpClient.Post(url, "text/plain", strings.NewReader(body))
	}
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(lua.LString(string(body)))
	headers := createHeadersTable(L, res.Header)
	L.Push(headers)
	L.Push(lua.LNumber(res.StatusCode))
	return 3
}

func createHeadersTable(L *lua.LState, header http.Header) *lua.LTable {
	table := L.NewTable()
	for name, value := range header {
		table.RawSetString(strings.ToLower(name), lua.LString(strings.Join(value, "\n")))
	}
	return table
}

// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwtlib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

// Normal
func TestJwtToken_none_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims == nil or msg ~= nil then
		return false
	end

	if claims:valid() ~= true then
		return false
	end

	local token, msg = jwt.newToken(jwt.SIGNING_METHOD_NONE, claims)
	if token == nil or msg ~= nil then
		return false
	end

	local signedString, msg = token:signed()
	if signedString == nil or msg ~= nil then
		return false
	end

	if token:valid() ~= false then
		return false
	end

	local newToken, msg = jwt.parse(signedString)
	if newToken == nil or msg ~= nil then
		return false
	end

	if newToken:valid() == false then
		return false
	end

	local claims = newToken:getClaims()
	if claims == nil then
		return false
	end
	--[[
	for k, v in pairs(claims) do
		print(k, v)
	end
	--]]

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtToken_none_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local clock = os.clock
	function sleep(n)  -- seconds
		local t0 = clock()
		while clock() - t0 <= n do end
	end

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims == nil or msg ~= nil then
		return false
	end

	if claims:valid() ~= true then
		return false
	end

	local token, msg = jwt.newToken(jwt.SIGNING_METHOD_NONE, claims)
	if token == nil or msg ~= nil then
		return false
	end

	local signedString, msg = token:signed()
	if signedString == nil or msg ~= nil then
		return false
	end

	if token:valid() ~= false then
		return false
	end

	sleep(exp + 5)

	local newToken, msg = jwt.parse(signedString)
	if newToken ~= nil or msg == nil then
		return false
	end

	if msg ~= "Token is expired" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtToken_HS(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local testData = {
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_HS256,
			key="jY4n1%v9dhf3",
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_HS384,
			key="FJDL&da,!b6c5",
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_HS512,
			key="lUnfdf29ujs",
		},
	}

	for _, item in ipairs(testData) do
		local claims, msg = jwt.newClaims(item.claims)
		if claims == nil or msg ~= nil then
			return false
		end

		if claims:valid() ~= true then
			return false
		end

		local token, msg = jwt.newToken(item.method, claims)
		if token == nil or msg ~= nil then
			return false
		end

		local signedString, msg = token:signed(item.key)
		if signedString == nil or msg ~= nil then
			return false
		end

		if token:valid() ~= false then
			return false
		end

		local newToken, msg = jwt.parse(signedString, item.key)
		if newToken == nil or msg ~= nil then
			return false
		end
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtToken_RS(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local testData = {
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="JeffChen", sub="Gola is framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_RS256,
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola of Jeff"},
			method=jwt.SIGNING_METHOD_RS384,
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff Chen", sub="framework of Gola", aud="Gola of Jeff.Chen"},
			method=jwt.SIGNING_METHOD_RS512,
		},
	}

	local priKey = [[-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----]]

	local pubKey = [[-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----]]

	for _, item in ipairs(testData) do
		local claims, msg = jwt.newClaims(item.claims)
		if claims == nil or msg ~= nil then
			return false
		end

		if claims:valid() ~= true then
			return false
		end

		local token, msg = jwt.newToken(item.method, claims)
		if token == nil or msg ~= nil then
			return false
		end

		local signedString, msg = token:signed(priKey)
		if signedString == nil or msg ~= nil then
			return false
		end

		if token:valid() ~= false then
			return false
		end

		local newToken, msg = jwt.parse(signedString, pubKey)
		if newToken == nil or msg ~= nil then
			return false
		end

		local claims = newToken:getClaims()

		for k, v in pairs(claims) do
			if claims[k] ~= v then
				return false
			end
		end
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtToken_PS(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local testData = {
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="JeffChen", sub="Gola is framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_PS256,
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola of Jeff"},
			method=jwt.SIGNING_METHOD_PS384,
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff Chen", sub="framework of Gola", aud="Gola of Jeff.Chen"},
			method=jwt.SIGNING_METHOD_PS512,
		},
	}

	local secureKey = [[-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,7487BB8910A3741B

iL7m48mbFSIy1Y5xbXWwPTR07ufxu7o+myGUE+AdDeWWISkd5W6Gl44oX/jgXldS
mL/ntUXoZzQz2WKEYLwssAtSTGF+QgSIMvV5faiP+pLYvWgk0oVr42po00CvADFL
eDAJC7LgagYifS1l4EAK4MY8RGCHyJWEN5JAr0fc/Haa3WfWZ009kOWAp8MDuYxB
hQlCKUmnUpXCp5c6jwbjlyinLj8XwzzjZ/rVRsY+t2Z0Vcd5qzR5BV8IJCqbG5Py
z15/EFgMG2N2eYMsiEKgdXeKW2H5XIoWyun/3pBigWaDnTtiWSt9kz2MplqYfIT7
F+0XE3gdDGalAeN3YwFPHCkxxBmcI+s6lQG9INmf2/gkJQ+MOZBVXKmGLv6Qis3l
0eyUz1yZvNzf0zlcUBjiPulLF3peThHMEzhSsATfPomyg5NJ0X7ttd0ybnq+sPe4
qg2OJ8qNhYrqnx7Xlvj61+B2NAZVHvIioma1FzqX8DxQYrnR5S6DJExDqvzNxEz6
5VPQlH2Ig4hTvNzla84WgJ6USc/2SS4ehCReiNvfeNG9sPZKQnr/Ss8KPIYsKGcC
Pz/vEqbWDmJwHb7KixCQKPt1EbD+/uf0YnhskOWM15YiFbYAOZKJ5rcbz2Zu66vg
GAmqcBsHeFR3s/bObEzjxOmMfSr1vzvr4ActNJWVtfNKZNobSehZiMSHL54AXAZW
Yj48pwTbf7b1sbF0FeCuwTFiYxM+yiZVO5ciYOfmo4HUg53PjknKpcKtEFSj02P1
8JRBSb++V0IeMDyZLl12zgURDsvualbJMMBBR8emIpF13h0qdyah431gDhHGBnnC
J5UDGq21/flFjzz0x/Okjwf7mPK5pcmF+uW7AxtHqws6m93yD5+RFmfZ8cb/8CL8
jmsQslj+OIE64ykkRoJWpNBKyQjL3CnPnLmAB6TQKxegR94C7/hP1FvRW+W0AgZy
g2QczKQU3KBQP18Ui1HTbkOUJT0Lsy4FnmJFCB/STPRo6NlJiATKHq/cqHWQUvZd
d4oTMb1opKfs7AI9wiJBuskpGAECdRnVduml3dT4p//3BiP6K9ImWMSJeFpjFAFs
AbBMKyitMs0Fyn9AJRPl23TKVQ3cYeSTxus4wLmx5ECSsHRV6g06nYjBp4GWEqSX
RVclXF3zmy3b1+O5s2chJN6TrypzYSEYXJb1vvQLK0lNXqwxZAFV7Roi6xSG0fSY
EAtdUifLonu43EkrLh55KEwkXdVV8xneUjh+TF8VgJKMnqDFfeHFdmN53YYh3n3F
kpYSmVLRzQmLbH9dY+7kqvnsQm8y76vjug3p4IbEbHp/fNGf+gv7KDng1HyCl9A+
Ow/Hlr0NqCAIhminScbRsZ4SgbRTRgGEYZXvyOtQa/uL6I8t2NR4W7ynispMs0QL
RD61i3++bQXuTi4i8dg3yqIfe9S22NHSzZY/lAHAmmc3r5NrQ1TM1hsSxXawT5CU
anWFjbH6YQ/QplkkAqZMpropWn6ZdNDg/+BUjukDs0HZrbdGy846WxQUvE7G2bAw
IFQ1SymBZBtfnZXhfAXOHoWh017p6HsIkb2xmFrigMj7Jh10VVhdWg==
-----END RSA PRIVATE KEY-----]]

local pubKey = [[-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----]]

	for _, item in ipairs(testData) do
		local claims, msg = jwt.newClaims(item.claims)
		if claims == nil or msg ~= nil then
			return false
		end

		if claims:valid() ~= true then
			return false
		end

		local token, msg = jwt.newToken(item.method, claims)
		if token == nil or msg ~= nil then
			return false
		end

		local signedString, msg = token:signed(secureKey, "password")
		if signedString == nil or msg ~= nil then
			return false
		end

		if token:valid() ~= false then
			return false
		end

		local newToken, msg = jwt.parse(signedString, pubKey)
		if newToken ~= nil or msg == nil then
			return false
		end

		if msg ~= "crypto/rsa: verification error" then
			return false
		end
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtToken_ES(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local testData = {
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="JeffChen", sub="Gola is framework", aud="Gola"},
			method=jwt.SIGNING_METHOD_ES256,
			keys={pub=[[-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYD54V/vp+54P9DXarYqx4MPcm+HK
RIQzNasYSoRQHQ/6S6Ps8tpMcT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END PUBLIC KEY-----
]], pri=[[-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIAh5qA3rmqQQuu0vbKV/+zouz/y/Iy2pLpIcWUSyImSwoAoGCCqGSM49
AwEHoUQDQgAEYD54V/vp+54P9DXarYqx4MPcm+HKRIQzNasYSoRQHQ/6S6Ps8tpM
cT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END EC PRIVATE KEY-----
]]},
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola of Jeff"},
			method=jwt.SIGNING_METHOD_ES384,
			keys={pub=[[-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE2Gu6yWlO33izKqwMemnWRTUDqWlB/1Mp
W2S64jizjeWBVcLVVKTtRXxF2VCEKabUT7XOYsSD2OufMoQUm+oq3yCVEj/WmE35
SRipP5g5CuLfJlgc4Yg+Qeo4PXBBoRoA
-----END PUBLIC KEY-----
]], pri=[[-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDCaCvMHKhcG/qT7xsNLYnDT7sE/D+TtWIol1ROdaK1a564vx5pHbsRy
SEKcIxISi1igBwYFK4EEACKhZANiAATYa7rJaU7feLMqrAx6adZFNQOpaUH/Uylb
ZLriOLON5YFVwtVUpO1FfEXZUIQpptRPtc5ixIPY658yhBSb6irfIJUSP9aYTflJ
GKk/mDkK4t8mWBzhiD5B6jg9cEGhGgA=
-----END EC PRIVATE KEY-----
]]}
		},
		{
			claims={nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff Chen", sub="framework of Gola", aud="Gola of Jeff.Chen"},
			method=jwt.SIGNING_METHOD_ES512,
			keys={pub=[[-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQACVyLsNdjFM6R4IImvTzgRWF0sWjh
ihmzIyMgyPuqu8IuyzMNx4G2jpoCKhRu9qPCQUMGDeCG1x3/n/OgkWNQANsB82x7
7eiIZAl0zcQRH32tcjFILvJ/XCihdoi4MkCnCqlt9/HxjsP590ZtmHfxAeertq5w
9vakvpzjPXhkvoMt/Tk=
-----END PUBLIC KEY-----
]], pri=[[-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIB0pE4uFaWRx7t03BsYlYvF1YvKaBGyvoakxnodm9ou0R9wC+sJAjH
QZZJikOg4SwNqgQ/hyrOuDK2oAVHhgVGcYmgBwYFK4EEACOhgYkDgYYABAAJXIuw
12MUzpHggia9POBFYXSxaOGKGbMjIyDI+6q7wi7LMw3HgbaOmgIqFG72o8JBQwYN
4IbXHf+f86CRY1AA2wHzbHvt6IhkCXTNxBEffa1yMUgu8n9cKKF2iLgyQKcKqW33
8fGOw/n3Rm2Yd/EB56u2rnD29qS+nOM9eGS+gy39OQ==
-----END EC PRIVATE KEY-----
]]}
		},
	}

	for _, item in ipairs(testData) do
		local claims, msg = jwt.newClaims(item.claims)
		if claims == nil or msg ~= nil then
			return false
		end

		if claims:valid() ~= true then
			return false
		end

		local token, msg = jwt.newToken(item.method, claims)
		if token == nil or msg ~= nil then
			return false
		end

		local signedString, msg = token:signed(item.keys.pri)
		if signedString == nil or msg ~= nil then
			return false
		end

		if token:valid() ~= false then
			return false
		end

		local newToken, msg = jwt.parse(signedString, item.keys.pub)
		if newToken == nil or msg ~= nil then
			return false
		end

		local claims = newToken:getClaims()

		for k, v in pairs(claims) do
			if claims[k] ~= v then
				return false
			end
		end
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

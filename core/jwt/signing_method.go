// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	// INVALID
	SIGNING_METHOD_INVALID_TYPE SigningMethodType = iota
	// NONE
	SIGNING_METHOD_NONE_TYPE
	// HS
	SIGNING_METHOD_HS_TYPE
	// ES
	SIGNING_METHOD_ES_TYPE
	// RS
	SIGNING_METHOD_RS_TYPE
	// PS
	SIGNING_METHOD_PS_TYPE
)

const (
	// INVALID
	SIGNING_METHOD_INVALID SigningMethod = iota
	// NONE
	SIGNING_METHOD_NONE

	// HS
	SIGNING_METHOD_HS256
	SIGNING_METHOD_HS384
	SIGNING_METHOD_HS512

	// ES
	SIGNING_METHOD_ES256
	SIGNING_METHOD_ES384
	SIGNING_METHOD_ES512

	// RS
	SIGNING_METHOD_RS256
	SIGNING_METHOD_RS384
	SIGNING_METHOD_RS512

	// PS
	SIGNING_METHOD_PS256
	SIGNING_METHOD_PS384
	SIGNING_METHOD_PS512
)

var (
	ErrInvalidSigningMethod = errors.New("signing method is invalid")
)

type (
	SigningMethod     int
	SigningMethodType int
)

func signingMethod(method SigningMethod) (jwt.SigningMethod, SigningMethodType, error) {
	switch method {
	case SIGNING_METHOD_NONE:
		return jwt.SigningMethodNone, SIGNING_METHOD_NONE_TYPE, nil
	case SIGNING_METHOD_HS256:
		return jwt.SigningMethodHS256, SIGNING_METHOD_HS_TYPE, nil
	case SIGNING_METHOD_HS384:
		return jwt.SigningMethodHS384, SIGNING_METHOD_HS_TYPE, nil
	case SIGNING_METHOD_HS512:
		return jwt.SigningMethodHS512, SIGNING_METHOD_HS_TYPE, nil
	case SIGNING_METHOD_ES256:
		return jwt.SigningMethodES256, SIGNING_METHOD_ES_TYPE, nil
	case SIGNING_METHOD_ES384:
		return jwt.SigningMethodES384, SIGNING_METHOD_ES_TYPE, nil
	case SIGNING_METHOD_ES512:
		return jwt.SigningMethodES512, SIGNING_METHOD_ES_TYPE, nil
	case SIGNING_METHOD_RS256:
		return jwt.SigningMethodRS256, SIGNING_METHOD_RS_TYPE, nil
	case SIGNING_METHOD_RS384:
		return jwt.SigningMethodRS384, SIGNING_METHOD_RS_TYPE, nil
	case SIGNING_METHOD_RS512:
		return jwt.SigningMethodRS512, SIGNING_METHOD_RS_TYPE, nil
	case SIGNING_METHOD_PS256:
		return jwt.SigningMethodPS256, SIGNING_METHOD_PS_TYPE, nil
	case SIGNING_METHOD_PS384:
		return jwt.SigningMethodPS384, SIGNING_METHOD_PS_TYPE, nil
	case SIGNING_METHOD_PS512:
		return jwt.SigningMethodPS512, SIGNING_METHOD_PS_TYPE, nil
	}

	return nil, SIGNING_METHOD_INVALID_TYPE, ErrInvalidSigningMethod
}

package services

import "errors"

var ErrUnauthorized error = errors.New("user not authorized")

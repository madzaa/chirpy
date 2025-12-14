package services

import "errors"

var ErrUnauthorized = errors.New("user not authorized")
var ErrUserNotFound = errors.New("user not found")

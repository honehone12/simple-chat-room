package main

import "errors"

var (
	ErrorCastFailed   = errors.New("could not cast")
	ErrorPlayerExists = errors.New("player already exists")
)

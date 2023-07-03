package main

import "errors"

var (
	ErrorCastFailed         = errors.New("could not cast")
	ErrorPlayerExists       = errors.New("player already exists")
	ErrorAlreadyMarkedAsBad = errors.New("connection was marked as bad already")
)

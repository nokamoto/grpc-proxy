package main

import (
	"golang.org/x/net/context"
)

type cluster interface {
	unary(context.Context, string, *message) (*message, error)
}

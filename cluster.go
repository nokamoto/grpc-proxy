package main

import (
	"golang.org/x/net/context"
)

type cluster interface {
	invokeUnary(context.Context, *message, string) (*message, error)
}

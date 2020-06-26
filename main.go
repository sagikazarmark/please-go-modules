package main

import (
	"fmt"

	"github.com/pkg/errors"
	"logur.dev/logur"

	"github.com/sagikazarmark/please-go-modules/lib"
)

var _ = errors.New("")

func main() {
	_ = logur.NewNoopLogger()

	fmt.Println(lib.HelloWorld)
}

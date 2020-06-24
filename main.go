package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sagikazarmark/please-go-modules/lib"
)

var _ = errors.New("")

func main() {
	fmt.Println(lib.HelloWorld)
}

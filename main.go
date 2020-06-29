package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"

	"github.com/sagikazarmark/please-go-modules/lib"
)

var _ = errors.New("")

func main() {
	_ = logur.NewNoopLogger()
	_ = logrusadapter.New(logrus.New())

	fmt.Println(lib.HelloWorld)
}

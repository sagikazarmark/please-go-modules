package main

import (
	_ "emperror.dev/errors/match"
	_ "github.com/containerd/containerd/sys"
	_ "github.com/golang/snappy"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/opencontainers/runc/libcontainer/system"
	_ "golang.org/x/sys/unix"
	_ "google.golang.org/grpc"
)

func main() {

}

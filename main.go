package main

import (
	"github.com/kataras/golog"
	"github.com/mateo08c/go-skolengo/examples"
)

func main() {
	golog.SetLevel("debug")
	examples.Start()
}

package main

import "github.com/yusukebe/rj"

var version = ""

func main() {
	rj.Version = version
	rj.Execute()
}

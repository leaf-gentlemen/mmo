package main

import "github.com/leaf-gentlemen/zinx/znet"

func main() {
	srv := znet.NewServe("")
	srv.Start()
}

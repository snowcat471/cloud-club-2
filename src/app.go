package main

import "test-app/server"

func main() {
	s := server.NewGinServer(3000)
	s.Run()
}

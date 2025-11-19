package main

import "insight/cmd"

// main is the program entry point. It delegates initialization and execution of the application server to server.RunServer().
func main() {
	cmd.Execute()
}

package main

import (
	cmd "not_your_fathers_search_engine/cmd"
)

func init() {
	cmd.InitializeApp()
}

func main() {
	cmd.StartApp()
}
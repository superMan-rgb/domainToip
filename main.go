package main

import "domaintoip/cmd"

func main() {
	config := cmd.FlagConfig{}
	config.ParseFlags()
}

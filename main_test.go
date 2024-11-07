package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestName(t *testing.T) {
	parse, err := url.Parse("http://pililivertmp.open.freebuf.com:9090")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parse.Host)
}

package main

import (
	"blockchainparser/src/parser"
	"flag"
)

func main() {
	var dbuser = flag.String("user", "postgres", "postgres database user name")
	flag.Parse()
	parser.Parse(*dbuser)
}

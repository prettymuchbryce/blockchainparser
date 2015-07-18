package main

import (
	"blockchainparser/src/parser"
	"flag"
)

func main() {
	var dbuser = flag.String("user", "postgres", "postgres database user name")
	var datpath = flag.String("path", "./", "path to .dat files")
	flag.Parse()
	parser.Parse(*dbuser, *datpath)
}

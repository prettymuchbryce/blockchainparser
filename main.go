package main

import (
	"blockchainparser/parser"
	"flag"
	"runtime"
)

func main() {
	var dbuser = flag.String("user", "postgres", "postgres database user name")
	var datpath = flag.String("path", "./", "path to .dat files")
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	parser.Parse(*dbuser, *datpath)
}

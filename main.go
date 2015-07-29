package main

import (
	"blockchainparser/parser"
	"flag"
)

func main() {
	// var ba []byte = []byte{76, 104, 76, 85, 75, 69, 45, 74, 82, 32, 73, 83, 32, 65, 32, 80, 69, 68, 79, 80, 72, 73, 76, 69, 33, 32, 79, 104, 44, 32, 97, 110, 100, 32, 103, 111, 100, 32, 105, 115, 110, 39, 116, 32, 114, 101, 97, 108, 44, 32, 115, 117, 99, 107, 97, 46, 32, 83, 116, 111, 112, 32, 112, 111, 108, 108, 117, 116, 105, 110, 103, 32, 116, 104, 101, 32, 98, 108, 111, 99, 107, 99, 104, 97, 105, 110, 32, 119, 105, 116, 104, 32, 121, 111, 117, 114, 32, 110, 111, 110, 115, 101, 110, 115, 101, 46, 172}
	// key, err := utils.ExtractPublicKeyFromOutputScript(ba)
	// fmt.Println(key, err)
	var dbuser = flag.String("user", "postgres", "postgres database user name")
	var datpath = flag.String("path", "./", "path to .dat files")
	flag.Parse()
	parser.Parse(*dbuser, *datpath)
}

package main

import "flag"

func main() {
	flg, err := parseFlags()
	if err != nil {
		flag.PrintDefaults()
		return
	}
	_ = flg
}

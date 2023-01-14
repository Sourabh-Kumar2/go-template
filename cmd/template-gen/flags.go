package main

import (
	"errors"
	"flag"
)

type flags struct {
	path    string
	handler string
}

func parseFlags() (*flags, error) {
	flg := &flags{}
	flag.StringVar(&flg.handler, "handler", flg.handler, "Handler to be created. It's required")
	flag.StringVar(&flg.path, "path", flg.path, "root path. It's required")
	flag.Parse()
	if flg.handler == "" || flg.path == "" {
		return nil, errors.New("flag defaults not set")
	}
	return flg, nil
}

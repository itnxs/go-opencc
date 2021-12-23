package main

import (
	"fmt"
	occ "github.com/itnxs/go-opencc"
)

func main() {
	s2tw, err := occ.New(occ.S2TW)
	if err != nil {
		panic(err)
	}

	txt, err := s2tw.Convert("中国台湾简体转换")
	if err != nil {
		panic(err)
	}

	fmt.Println(txt)
}

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

	txt, err := s2tw.Convert("背包")
	if err != nil {
		panic(err)
	}

	fmt.Println(txt)
}

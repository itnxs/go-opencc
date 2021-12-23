package goopencc

import (
	"fmt"
	"strings"
	"testing"
)

func TestOpenCC_Option(t *testing.T) {
	dirs, err := fs.ReadDir("opencc")
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		name := dir.Name()
		if !strings.Contains(name, ".json") {
			continue
		}
		name = strings.Replace(name, ".json", "", 1)
		fmt.Println(fmt.Sprintf(`%s = "%s"`, strings.ToUpper(name), name))
	}
}

func TestS2T_Convert(t *testing.T) {
	s2t, err := New(S2T)
	if err != nil {
		t.Error(err)
	}

	txt, err := s2t.Convert("中国")
	if err != nil {
		t.Error(err)
	}

	if txt != "中國" {
		t.Error("convert error")
	}
}

func TestT2S_Convert(t *testing.T) {
	s2t, err := New(T2S)
	if err != nil {
		t.Error(err)
	}

	txt, err := s2t.Convert("中國")
	if err != nil {
		t.Error(err)
	}

	if txt != "中国" {
		t.Error("convert error")
	}
}
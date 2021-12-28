package goopencc

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	HK2S = "hk2s"
	S2HK = "s2hk"
	S2T = "s2t"
	S2TW = "s2tw"
	S2TWP = "s2twp"
	T2HK = "t2hk"
	T2S = "t2s"
	T2TW = "t2tw"
	TW2S = "tw2s"
	TW2SP = "tw2sp"
)

var punctuations = []string{
	" ", "\n", "\r", "\t", "-", ",", ".", "?", "!", "*", "　",
	"，", "。", "、", "；", "：", "？", "！", "…", "“", "”", "「",
	"」", "—", "－", "（", "）", "《", "》", "．", "／", "＼"}

// OpenCC OpenCC
type OpenCC struct {
	o *Option
}

// New 新建OpenCC
func New(dict DictType) (*OpenCC, error) {
	body, err := fs.ReadFile(fmt.Sprintf("opencc/%s.json", dict))
	if err != nil {
		return nil, err
	}

	var o *Option
	err = json.Unmarshal(body, &o)
	if err != nil {
		return nil, err
	}

	err = o.init()
	if err != nil {
		return nil, err
	}

	return &OpenCC{o: o}, nil
}

// Convert 转换
func (oc *OpenCC) Convert(txt string) (text string, err error) {
	vs := make([]string, 0, len(txt))
	for i, c := range strings.Split(txt, "") {
		if i > 0 && isPunctuations(c) {
			if len(vs) > 0 {
				tx, err := oc.convert(strings.Join(vs, ""))
				if err != nil {
					return txt, err
				}
				text = text + tx + c
				vs = vs[:0]
			} else {
				text = text + c
			}
			continue
		}
		vs = append(vs, c)
	}

	if len(vs) > 0 {
		tx, err := oc.convert(strings.Join(vs, ""))
		if err != nil {
			return txt, err
		}
		text = text + tx
	}

	return text, nil
}

func (oc *OpenCC) convert(txt string) (string, error) {
	var err error
	txt, err = oc.o.convert(txt)
	if err != nil {
		return txt, err
	}
	return txt, nil
}

func isPunctuations(character string) bool {
	if len([]byte(character)) <= 1 {
		return true
	}
	for _, c := range punctuations {
		if c == character {
			return true
		}
	}
	return false
}

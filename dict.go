package goopencc

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"
)

//go:embed opencc/*
var fs embed.FS

// DictType //
type DictType string

// Option //
type Option struct {
	Name            string             `json:"name"`
	Segmentation    Segmentation       `json:"segmentation"`
	ConversionChain []*ConversionChain `json:"conversion_chain"`
}

// Segmentation //
type Segmentation struct {
	Type string `json:"type"`
	Dict *Dict  `json:"dict"`
}

func (o *Option) init() (err error) {
	for _, cv := range o.ConversionChain {
		err = cv.init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Option) convert(txt string) (string, error) {
	var err error
	for _, cv := range o.ConversionChain {
		txt, err = cv.convert(txt)
		if err != nil {
			return txt, err
		}
	}
	return txt, nil
}

// ConversionChain //
type ConversionChain struct {
	Dict *Dict `json:"dict"`
}

func (chain *ConversionChain) init() error {
	return chain.Dict.init()
}

func (chain *ConversionChain) convert(txt string) (string, error) {
	var err error
	txt, err = chain.Dict.convert(txt)
	if err != nil {
		return txt, err
	}
	return txt, nil
}

// Dict //
type Dict struct {
	Type  string  `json:"type"`
	File  string  `json:"file"`
	Child []*Dict `json:"dicts"`
	data  map[string][]string
	max   int
	min   int
}

func (d *Dict) init() (err error) {
	if d.File != "" && d.data == nil {
		d.data, d.max, d.min, err = d.read(fmt.Sprintf("opencc/%s", d.File))
		if err != nil {
			return err
		}
	}

	if d.Child != nil && len(d.Child) > 0 {
		for _, cd := range d.Child {
			err = cd.init()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Dict) read(name string) (data map[string][]string, max, min int, err error) {
	data = make(map[string][]string, 0)

	fn, err := fs.Open(name)
	if err != nil {
		return data, max, min, err
	}

	r := bufio.NewReader(fn)
	for {
		row, err := r.ReadString('\n')
		if err == io.EOF {
			return data, max, min, nil
		} else if err != nil {
			return data, max, min, err
		}

		fields := strings.Fields(row)
		if len(fields) > 1 {
			if len([]rune(fields[0])) > max {
				max = len([]rune(fields[0]))
			}
			if min <= 0 || len([]rune(fields[0])) < min {
				min = len([]rune(fields[0]))
			}
			data[fields[0]] = fields[1:]
		}
	}
}

func (d *Dict) convert(txt string) (text string, err error) {
	text = txt
	runes := []rune(txt)

	if d.data != nil {
		if len(runes) < d.min {
			return text, nil
		}

		max := d.max
		if max > len(runes) {
			max = len(runes)
		}

		for i := max; i >= d.min; i-- {
			for j := 0; j <= len(runes)-i; j++ {
				if i == 0 || j+i > len(runes) {
					continue
				}
				old := string(runes[j : j+i])
				if newStr, ok := d.data[old]; ok {
					text = strings.Replace(text, old, newStr[0], 1)
					j = j + i - 1
				}
			}
		}
	}

	for _, cd := range d.Child {
		text, err = cd.convert(text)
		if err != nil {
			return text, err
		}
	}

	return text, nil
}

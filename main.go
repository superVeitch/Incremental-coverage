package main

import (
	"fmt"
	"golang.org/x/tools/cover"
	"io/ioutil"
	"strings"
)

func setup() *Diff {
	byt, err := ioutil.ReadFile("example1.diff")
	if err != nil {
		fmt.Println(err)
	}
	diff, err := Parse(string(byt))
	return diff
}

type DiffFileRe struct {
	Name string
	Lines []int
}

func main()  {

	// You now have a slice of files from the diff,
	files := setup()


	profile, err := cover.ParseProfiles("coverage.cov")
	if err != nil {
		fmt.Println("解析cover文件异常")
	}


	var diffs []DiffFileRe
	for _, f := range files.Files {
		var dl = make([]*DiffLine, 0)
		for _, h := range f.Hunks  {
			i := filter(h.NewRange.Lines)
			dl = append(dl, i...)
		}

		lines := coverFilter(f.NewName, dl, profile)
		if len(lines) > 0 {
			re := DiffFileRe{
				Name:  f.NewName,
				Lines: lines,
			}
			diffs = append(diffs, re)
		}
	}


	fmt.Println(diffs)
}

func coverFilter(fileName string, dl []*DiffLine, f []*cover.Profile) []int {
	dd := make([]int, 0)
	for _, d := range dl {
		for _, f1 := range f {
			name := strings.TrimPrefix(f1.FileName, "bible-go/")
			if name == fileName{
				for _, b := range f1.Blocks  {
					if b.EndLine >= d.Number && b.StartLine <= d.Number {
						dd = append(dd, d.Number)
					}
				}
			}
		}
	}
	return dd
}

func filter(dl []*DiffLine) []*DiffLine {
	r := make([]*DiffLine, 0)
	for _, v := range dl {
		if v.Mode == ADDED {
			r = append(r, v)
		}
	}
	return r
}

package main

import (
	"fmt"
	"io/ioutil"
)

func setup() *Diff {
	byt, err := ioutil.ReadFile("example.diff")
	if err != nil {
		fmt.Println(err)
	}
	diff, err := Parse(string(byt))
	return diff
}

type DiffLine1 struct {
	Name string
	Lines []int
}

func main()  {

	// You now have a slice of files from the diff,
	files := setup()


	//profile := cover.ParseProfiles("coverage.cov")
	var diffs []DiffLine
	for _, f := range files.Files {
		for _, h := range f.Hunks  {
			i := h.NewRange.Lines
			filter(i.)
		}
	}

}

func filter(dl []DiffLine) []DiffLine{
	r := make([]DiffLine, 0)
	for _, v := range dl {
		if v.Mode == ADDED {
			r = append(r, v)
		}
	}
	return r
}

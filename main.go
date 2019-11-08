package main

import (
	"fmt"
	"golang.org/x/tools/cover"
	"io/ioutil"
	"os"
	"strings"
)

func parseFile() *Diff {
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
	NoLines []int
}

type DiffFileR struct {
	Name string
	Lines []Line
	NoLines []Line
}

type Line struct {
	Name string
	Num int
	Hit int
}

func main()  {

	// You now have a slice of files from the diff,
	files := parseFile()


	profile, err := cover.ParseProfiles("coverage.cov")
	if err != nil {
		fmt.Println("解析cover文件异常")
	}


	var diffs []DiffFileR
	for _, f := range files.Files {
		var dl = make([]*DiffLine, 0)
		for _, h := range f.Hunks  {
			i := filter(h.NewRange.Lines)
			dl = append(dl, i...)
		}

		lines, noLines := coverFilter(f.NewName, dl, profile)
		if len(lines) > 0 {
			re := DiffFileR{
				Name: f.NewName,
				Lines: lines,
				NoLines: noLines,
			}
			diffs = append(diffs, re)
		}
	}
	var prof []cover.Profile
	if len(diffs) > 0 {
		for _, f := range diffs  {
			var bs []cover.ProfileBlock
			for _, l := range f.Lines {
				cp := cover.ProfileBlock{
					StartLine: l.Num,
					EndLine: l.Num,
					StartCol:1,
					EndCol:200,
					Count: l.Hit,
					NumStmt:1,
				}
				bs = append(bs, cp)
			}

			for _, l := range f.NoLines {
				cp := cover.ProfileBlock{
					StartLine: l.Num,
					EndLine: l.Num,
					StartCol:1,
					EndCol:200,
					Count: 0,
					NumStmt:1,
				}
				bs = append(bs, cp)
			}


			pro := cover.Profile{
				FileName: f.Name,
				Mode: "set",
				Blocks: bs,
			}
			prof = append(prof, pro)
		}
	}

	f, _ := os.Open("cover.out")
	fmt.Fprintln(f, prof)
}

func coverFilter(fileName string, dl []*DiffLine, f []*cover.Profile) ([]Line, []Line) {
	dd := make([]Line, 0)
	da := make([]Line, 0)
	for _, d := range dl {
		for _, f1 := range f {
			name := strings.TrimPrefix(f1.FileName, "go-main/")
			if name == fileName {
				for _, b := range f1.Blocks  {
					if b.EndLine >= d.Number && b.StartLine <= d.Number {
						l := Line {
							Name: name,
							Num: d.Number,
							Hit: b.Count,
						}
						dd = append(dd, l)
					} else {
						l := Line {
							Name: name,
							Num: d.Number,
							Hit: 0,
						}
						da = append(da, l)
					}
				}
			}
		}
	}
	return dd, da
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

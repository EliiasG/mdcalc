package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eliiasg/mdcalc/parse"
	"github.com/eliiasg/mdcalc/unitlib"
)

// Terrible code, I know

func main() {
	if len(os.Args) != 4 {
		fmt.Println("must call with 2 param (dir, prob name, title)")
		return
	}
	var doc strings.Builder
	doc.WriteString(fmt.Sprintf("<span style=\"font-size:0\">\n# %v\n</span>\n\n", os.Args[3]))
	lib, err := unitlib.NewSavedUnitLib(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	n := 1
	for {
		dat, err := os.ReadFile(fmt.Sprintf("%v/%v.mdc", os.Args[1], n))
		if err != nil {
			break
		}
		res, err := parse.Parse(string(dat), fmt.Sprintf("%v %v", os.Args[2], n), fmt.Sprintf("%v.<n>", n), lib)
		if err != nil {
			fmt.Printf("Error in file %v.mdc\n", n)
			fmt.Println(err.Error())
			return
		}
		doc.WriteString(res)
		doc.WriteString("  \n\n")
		n++
	}
	path := os.Args[1] + "/Result.md"
	os.Remove(path)
	f, _ := os.Create(path)
	f.Write([]byte(doc.String()))
}

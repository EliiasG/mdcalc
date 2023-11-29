package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eliiasg/mdcalc/parse"
)

// Terrible code, I know

func main() {
	if len(os.Args) != 3 {
		fmt.Println("must call with 2 param (dir, prob name)")
		return
	}
	var doc strings.Builder
	n := 1
	for {
		dat, err := os.ReadFile(fmt.Sprintf("%v/%v.mdc", os.Args[1], n))
		if err != nil {
			break
		}
		res, err := parse.Parse(string(dat), fmt.Sprintf("%v %v", os.Args[2], n), fmt.Sprintf("%v.<n>", n))
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

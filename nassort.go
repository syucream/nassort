package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v2"
)

type NassortDefinition struct {
	Dst string
	Src []SourceCondition
}

type SourceCondition struct {
	Contains []string
	// TODO support regexp condition
}

func (defs NassortDefinition) judgeConditions(filename string) bool {
	for _, cond := range defs.Src {
		if matched := cond.judge(filename); matched {
			return true
		}
	}

	return false
}

func (cond SourceCondition) judge(filename string) bool {
	for _, substr := range cond.Contains {
		if contained := strings.Contains(filename, substr); contained {
			return true
		}
	}

	// TODO support regexp condition

	return false
}

func main() {
	srcDirPath := flag.String("src", ".", "src directory")
	dstDirPath := flag.String("dst", ".", "dst directory")
	defsPath := flag.String("f", "nassort.yaml", "nassort config file(only yaml supported)")

	flag.Parse()

	srcFilePaths, err := ioutil.ReadDir(*srcDirPath)
	if err != nil {
		fmt.Println("cannot read the src directory: " + *srcDirPath)
		os.Exit(1)
	}

	var defs []NassortDefinition
	defsData, err := ioutil.ReadFile(*defsPath)
	if err != nil {
		fmt.Println("cannot read the nassort config file: " + *defsPath)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(defsData, &defs); err != nil {
		fmt.Println("cannot parse the nassort config file: " + *defsPath)
		os.Exit(1)
	}

	for _, fileInfo := range srcFilePaths {
		for _, def := range defs {
			filename := fileInfo.Name()

			// NOTE denormalize a filename normalized by UTF-8-MAC
			// see also http://blog.sarabande.jp/post/89636452673
			buf := []byte(filename)
			filename = string(norm.NFC.Bytes(buf))

			if matched := def.judgeConditions(filename); matched {
				dstPath := *dstDirPath + def.Dst

				if err := os.MkdirAll(dstPath, 0755); err != nil {
					fmt.Println("cannot make dst directory, " + dstPath)
					os.Exit(1)
				}

				origPath := *srcDirPath + filename
				movedPath := dstPath + "/" + filename
				if err := os.Rename(origPath, movedPath); err != nil {
					fmt.Println("Skip moving %s to %s", origPath, movedPath)
				}
			}
		}
	}

	fmt.Println("assorting done.")
}

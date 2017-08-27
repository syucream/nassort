package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func Scan(srcDirPath string) ([]string, error) {
	var paths []string

	err := filepath.Walk(srcDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip root dir and files
		if path != srcDirPath && info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	return paths, err
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

	// srcFilePaths, err := ioutil.ReadDir(*srcDirPath)
	srcFilePaths, err := Scan(*srcDirPath)
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

	for _, def := range defs {
		for _, path := range srcFilePaths {
			// NOTE denormalize a filename normalized by UTF-8-MAC
			// see also http://blog.sarabande.jp/post/89636452673
			buf := []byte(path)
			path = string(norm.NFC.Bytes(buf))

			if matched := def.judgeConditions(path); matched {
				dstPath := *dstDirPath + def.Dst

				if err := os.MkdirAll(dstPath, 0755); err != nil {
					fmt.Println("cannot make dst directory, " + dstPath)
					os.Exit(1)
				}

				movedPath := strings.Replace(dstPath+"/"+path, *srcDirPath, "", 1)
				if err := os.Rename(path, movedPath); err != nil {
					fmt.Println("Skip moving %s to %s", path, movedPath)
				}
			}
		}
	}

	fmt.Println("assorting done.")
}

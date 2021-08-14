package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func walk(files []string, core func(string)) {
	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if fi.Mode().IsRegular() {
			core(f)
		} else if fi.IsDir() {
			filepath.WalkDir(f, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Println(err)
				} else {
					if d.Type().IsRegular() {
						core(path)
					}
				}
				return nil
			})
		}
	}
}

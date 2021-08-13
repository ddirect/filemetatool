package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ddirect/check"
	"github.com/ddirect/filemeta"
)

func main() {
	var do string
	flag.StringVar(&do, "do", "", "list|refresh|stat|scrub|inspect")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 || do == "" {
		flag.Usage()
		return
	}

	var core, queue func(string)
	var epilogue func()
	workerPool := true
	switch do {
	case "list":
		core = listCore
		workerPool = false
	case "refresh":
		core, epilogue = fetch(filemeta.Refresh)
	case "stat":
		core, epilogue = fetch(filemeta.Get)
	case "inspect":
		core, epilogue = fetch(filemeta.Inspect)
	case "scrub":
		core, epilogue = scrub()
	default:
		fmt.Fprintf(os.Stderr, "unknown operation '%s'\n", do)
		return
	}

	if epilogue != nil {
		defer epilogue()
	}

	if workerPool {
		workers := runtime.NumCPU()
		var wg sync.WaitGroup
		fileChannel := make(chan string, 4000)
		queue = func(path string) {
			fileChannel <- path
		}
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				for fileName := range fileChannel {
					core(fileName)
				}
				wg.Done()
			}()
		}
		defer func() {
			close(fileChannel)
			wg.Wait()
		}()
	} else {
		queue = core
	}

	for _, f := range files {
		fi, err := os.Stat(f)
		check.E(err)
		if fi.Mode().IsRegular() {
			queue(f)
		} else if fi.IsDir() {
			filepath.WalkDir(f, func(path string, d fs.DirEntry, err error) error {
				check.E(err)
				if d.Type().IsRegular() {
					queue(path)
				}
				return nil
			})
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

func ensureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func createLink(src, dest string) error {
	return os.Link(src, dest)
}

func processHashFiles(srcDir, destDir string, fileName string) error {
	prefix1 := fileName[:3]
	prefix2 := fileName[3:6]
	destPath := filepath.Join(destDir, prefix1, prefix2, fileName)

	if err := ensureDir(filepath.Dir(destPath)); err != nil {
		return err
	}

	srcPath := filepath.Join(srcDir, fileName)
	return createLink(srcPath, destPath)
}

func processOrdinalFiles(srcDir, destDir string, fileName string) error {
	ordinal, err := strconv.Atoi(fileName)
	if err != nil {
		return err
	}

	lowerBound := (ordinal / 20000) * 20000

	destPath := filepath.Join(destDir, fmt.Sprintf("%d", lowerBound), fileName)

	if err := ensureDir(filepath.Dir(destPath)); err != nil {
		return err
	}

	srcPath := filepath.Join(srcDir, fileName)
	return createLink(srcPath, destPath)
}

func isHashName(fileName string) bool {
	return len(fileName) == 64
}

func worker(srcDir, hashDestDir, ordinalDestDir string, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for fileName := range jobs {
		if isHashName(fileName) {
			if err := processHashFiles(srcDir, hashDestDir, fileName); err != nil {
				fmt.Printf("Error processing hash file %s: %v\n", fileName, err)
			}
		} else {
			if err := processOrdinalFiles(srcDir, ordinalDestDir, fileName); err != nil {
				fmt.Printf("Error processing ordinal file %s: %v\n", fileName, err)
			}
		}
	}
}

func main() {
  srcDir := flag.String("src", "", "Source directory containing files")

  flag.Parse()

  if *srcDir == "" {
    fmt.Println("Error: srcDir is required.")
    flag.Usage()
    os.Exit(1)
  }

  hashDestDir := filepath.Join(*srcDir, "hash")
  ordinalDestDir := filepath.Join(*srcDir, "ordinal")

  numWorkers := runtime.NumCPU()
  jobs := make(chan string, numWorkers)

  var wg sync.WaitGroup

  for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go worker(*srcDir, hashDestDir, ordinalDestDir, jobs, &wg)
  }

  err := filepath.Walk(*srcDir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if !info.IsDir() {
      fileName := info.Name()
      jobs <- fileName
    }

    return nil
  })

  if err != nil {
    fmt.Printf("Error walking through the directory: %v\n", err)
    close(jobs)
    wg.Wait()
    return
  }

  close(jobs)
  wg.Wait()

  fmt.Println("Snapshots have been successfully migrated!")
}

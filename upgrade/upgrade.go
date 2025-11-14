//go:build !cgo && upgrade && ignore
// +build !cgo,upgrade,ignore

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func mergeFile(src string, dst string) error {
	defer func() error {
		fmt.Printf("Removing: %s\n", src)
		err := os.Remove(src)

		if err != nil {
			return err
		}

		return nil
	}()

	// Open destination
	fdst, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fdst.Close()

	// Read source content
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Add Additional newline
	if _, err := fdst.WriteString("\n"); err != nil {
		return err
	}

	fmt.Printf("Merging: %s into %s\n", src, dst)
	if _, err = fdst.Write(content); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Go-SQLite3 Upgrade Tool")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.Base(wd) != "upgrade" {
		log.Printf("Current directory is %q but should run in upgrade directory", wd)
		os.Exit(1)
	}

	// Use local files instead of downloading
	srcDir := "TODO:local file"

	// Check if source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", srcDir)
	}

	// Create Source Zip Reader
	//rSource, err := zip.NewReader(bytes.NewReader(source), int64(len(source)))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Extract from local directory
	filesToExtract := map[string]string{
		"sqlite3.c":    "../sqlite3-binding.c",
		"sqlite3.h":    "../sqlite3-binding.h",
		"sqlite3ext.h": "../sqlite3ext.h",
	}

	for filename, outputPath := range filesToExtract {
		srcPath := filepath.Join(srcDir, filename)

		// Check if file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", srcPath)
			continue
		}

		f, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}
		zr, err := os.Open(srcPath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.WriteString(f, "#ifndef USE_LIBSQLITE3\n")
		if err != nil {
			zr.Close()
			f.Close()
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(zr)
		for scanner.Scan() {
			text := scanner.Text()
			if text == `#include "sqlite3.h"` {
				text = `#include "sqlite3-binding.h"
#ifdef __clang__
#define assert(condition) ((void)0)
#endif
`
			}
			_, err = fmt.Fprintln(f, text)
			if err != nil {
				break
			}
		}
		err = scanner.Err()
		if err != nil {
			zr.Close()
			f.Close()
			log.Fatal(err)
		}
		_, err = io.WriteString(f, "#else // USE_LIBSQLITE3\n // If users really want to link against the system sqlite3 we\n// need to make this file a noop.\n #endif")
		if err != nil {
			zr.Close()
			f.Close()
			log.Fatal(err)
		}
		zr.Close()
		f.Close()
		fmt.Printf("Extracted: %v\n", filepath.Base(f.Name()))
	}

	//Extract Source
	//for _, zf := range rSource.File {
	//	var f *os.File
	//	switch path.Base(zf.Name) {
	//	case "userauth.c":
	//		f, err = os.Create("../userauth.c")
	//	case "sqlite3userauth.h":
	//		f, err = os.Create("../userauth.h")
	//	default:
	//		continue
	//	}
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	zr, err := zf.Open()
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	//	_, err = io.Copy(f, zr)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	//	zr.Close()
	//	f.Close()
	//	fmt.Printf("extracted %v\n", filepath.Base(f.Name()))
	//}

	// Merge SQLite User Authentication into amalgamation
	//if err := mergeFile("../userauth.c", "../sqlite3-binding.c"); err != nil {
	//	log.Fatal(err)
	//}
	//if err := mergeFile("../userauth.h", "../sqlite3-binding.h"); err != nil {
	//	log.Fatal(err)
	//}

	os.Exit(0)
}

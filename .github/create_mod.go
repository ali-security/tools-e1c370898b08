//go:build ignore
// +build ignore

// create_mod.go builds a Go module zip with golang.org/x/mod/zip.CreateFromDir,
// the same implementation the module proxy uses, so the file set matches the
// public zip by construction. It is run from a separate helper module by the
// seal-build workflow; the build tag keeps it out of the library's own builds.
//
// Usage: go run create_mod.go <module-path> <version> <source-dir> <output-zip>
package main

import (
	"log"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: go run create_mod.go <module-path> <version> <source-dir> <output-zip>")
	}
	m := module.Version{Path: os.Args[1], Version: os.Args[2]}
	f, err := os.Create(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	if err := zip.CreateFromDir(f, m, os.Args[3]); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("created module zip: %s", os.Args[4])
}

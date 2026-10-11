// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestStaticIsUpToDate(t *testing.T) {
	if runtime.GOOS == "android" {
		t.Skip("files not available on android")
	}
	oldBuf, err := ioutil.ReadFile("static.go")
	if err != nil {
		t.Errorf("error while reading static.go: %v\n", err)
	}

	newBuf, err := Generate()
	if err != nil {
		t.Errorf("error while generating static.go: %v\n", err)
	}

	if !bytes.Equal(oldBuf, newBuf) {
		t.Error(`static.go is stale.  Run:
  $ go generate golang.org/x/tools/godoc/static
  $ git diff
to see the differences.`)

	}
}

// TestAppendQuote ensures that AppendQuote produces a valid literal.
func TestAppendQuote(t *testing.T) {
	var in, out bytes.Buffer
	for r := rune(0); r < unicode.MaxRune; r++ {
		in.WriteRune(r)
	}
	appendQuote(&out, in.Bytes())
	in2, err := strconv.Unquote(out.String())
	if err != nil {
		t.Fatalf("AppendQuote produced invalid string literal: %v", err)
	}
	if got, want := in2, in.String(); got != want {
		t.Fatal("AppendQuote modified string") // no point printing got/want: huge
	}
}

// TestJQueryOptionWrapIsGated ensures the served jquery.js does not
// unconditionally wrap HTML containing <option> elements in a
// <select multiple> container before parsing it (CVE-2020-11023).
func TestJQueryOptionWrapIsGated(t *testing.T) {
	js, ok := Files["jquery.js"]
	if !ok {
		t.Fatal("jquery.js missing from Files")
	}
	if strings.Contains(js, `option:[1,"<select multiple='multiple'>","</select>"]`) {
		t.Error("jquery.js unconditionally wraps <option> HTML in <select multiple> (CVE-2020-11023)")
	}
	if !strings.Contains(js, `support.option||(`) {
		t.Error("jquery.js does not gate the <option> wrapMap entry on support.option (CVE-2020-11023)")
	}
}

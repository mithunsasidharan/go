// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Make sure tracebacks from initialization code are reported correctly.

package a

import (
	"fmt"
	"runtime"
	"strings"
)

var x = f() // line 15

func f() int {
	var b [4096]byte
	n := runtime.Stack(b[:], false) // line 19
	s := string(b[:n])
	var pcs [10]uintptr
	n = runtime.Callers(1, pcs[:]) // line 22

	// Check the Stack results.
	if debug {
		println(s)
	}
	if strings.Contains(s, "autogenerated") {
		panic("autogenerated code in traceback")
	}
	if !strings.Contains(s, "a.go:15") {
		panic("missing a.go:15")
	}
	if !strings.Contains(s, "a.go:19") {
		panic("missing a.go:19")
	}
	if !strings.Contains(s, "a.init.ializers") {
		panic("missing a.init.ializers")
	}

	// Check the CallersFrames results.
	if debug {
		iter := runtime.CallersFrames(pcs[:n])
		for {
			f, more := iter.Next()
			fmt.Printf("%s %s:%d\n", f.Function, f.File, f.Line)
			if !more {
				break
			}
		}
	}
	iter := runtime.CallersFrames(pcs[:n])
	f, more := iter.Next()
	if f.Function != "a.f" || !strings.HasSuffix(f.File, "a.go") || f.Line != 22 {
		panic(fmt.Sprintf("bad f %v\n", f))
	}
	if !more {
		panic("traceback truncated after f")
	}
	f, more = iter.Next()
	if f.Function != "a.init.ializers" || !strings.HasSuffix(f.File, "a.go") || f.Line != 15 {
		panic(fmt.Sprintf("bad init.ializers %v\n", f))
	}
	if !more {
		panic("traceback truncated after init.ializers")
	}
	f, _ = iter.Next()
	if f.Function != "runtime.main" {
		panic("runtime.main missing")
	}

	return 0
}

const debug = false

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// NOTE: This is from http://golang.org/src/pkg/net/http/httputil/dump.go 
// 

package sion

import(
	"io"
	"io/ioutil"
	"bytes"
)

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
    var buf bytes.Buffer
    if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
    if err = b.Close(); err != nil {
		return nil, nil, err
	}
    return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
}

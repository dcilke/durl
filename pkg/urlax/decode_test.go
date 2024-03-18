// This is a modification of the "net/url" package from the Go standard library.
// Everything but "Parse()" and supporting functions have been omitted. Further,
// "Parse()" has been modified to ignore any encoding errors.

// The original source code can be found here:
// https://github.com/golang/go/blob/master/src/net/url/url_test.go
// The original copyright is attached below:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package urlax

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	for in, out := range map[string]string{
		"http://www.google.com":                                           "http://www.google.com",
		"http://www.google.com/":                                          "http://www.google.com/",
		"http://www.google.com/file%20one%26two":                          "http://www.google.com/file one&two",
		"http://www.google.com/#file%20one%26two":                         "http://www.google.com/#file one&two",
		"ftp://john%20doe@www.google.com/":                                "ftp://john doe@www.google.com/",
		"http://www.google.com/?q=go%20language":                          "http://www.google.com/?q=go%20language",
		"ftp://webmaster@www.google.com/":                                 "ftp://webmaster@www.google.com/",
		"http://www.google.com/?":                                         "http://www.google.com/?",
		"http://www.google.com/?foo=bar?":                                 "http://www.google.com/?foo=bar?",
		"http://www.google.com/?q=go+language":                            "http://www.google.com/?q=go+language",
		"http://www.google.com/a%20b?q=c+d":                               "http://www.google.com/a b?q=c+d",
		"http:www.google.com/?q=go+language":                              "http:www.google.com/?q=go+language",
		"http:%2f%2fwww.google.com/?q=go+language":                        "http:%2f%2fwww.google.com/?q=go+language",
		"mailto:/webmaster@golang.org":                                    "mailto:/webmaster@golang.org",
		"mailto:webmaster@golang.org":                                     "mailto:webmaster@golang.org",
		"/foo?query=http://bad":                                           "/foo?query=http://bad",
		"//foo":                                                           "//foo",
		"//user@foo/path?a=b":                                             "//user@foo/path?a=b",
		"///threeslashes":                                                 "///threeslashes",
		"http://user:password@google.com":                                 "http://user:password@google.com",
		"http://j@ne:password@google.com":                                 "http://j@ne:password@google.com",
		"http://jane:p@ssword@google.com":                                 "http://jane:p@ssword@google.com",
		"http://j@ne:password@google.com/p@th?q=@go":                      "http://j@ne:password@google.com/p@th?q=@go",
		"http://www.google.com/?q=go+language#foo":                        "http://www.google.com/?q=go+language#foo",
		"http://www.google.com/?q=go+language#foo&bar":                    "http://www.google.com/?q=go+language#foo&bar",
		"http://www.google.com/?q=go+language#foo%26bar":                  "http://www.google.com/?q=go+language#foo&bar",
		"file:///home/adg/rabbits":                                        "file:///home/adg/rabbits",
		"file:///C:/FooBar/Baz.txt":                                       "file:///C:/FooBar/Baz.txt",
		"MaIlTo:webmaster@golang.org":                                     "mailto:webmaster@golang.org",
		"a/b/c":                                                           "a/b/c",
		"http://%3Fam:pa%3Fsword@google.com":                              "http://?am:pa?sword@google.com",
		"http://192.168.0.1/":                                             "http://192.168.0.1/",
		"http://192.168.0.1:8080/":                                        "http://192.168.0.1:8080/",
		"http://[fe80::1]/":                                               "http://[fe80::1]/",
		"http://[fe80::1]:8080/":                                          "http://[fe80::1]:8080/",
		"http://[fe80::1%25en0]/":                                         "http://[fe80::1%en0]/",
		"http://[fe80::1%25en0]:8080/":                                    "http://[fe80::1%en0]:8080/",
		"http://[fe80::1%25%65%6e%301-._~]/":                              "http://[fe80::1%en01-._~]/",
		"http://[fe80::1%25%65%6e%301-._~]:8080/":                         "http://[fe80::1%en01-._~]:8080/",
		"http://rest.rsc.io/foo%2fbar/baz%2Fquux?alt=media":               "http://rest.rsc.io/foo/bar/baz/quux?alt=media",
		"mysql://a,b,c/bar":                                               "mysql://a,b,c/bar",
		"http://host/!$&'()*+,;=:@[hello]":                                "http://host/!$&'()*+,;=:@[hello]",
		"http://example.com/oid/[order_id]":                               "http://example.com/oid/[order_id]",
		"http://192.168.0.2:8080/foo":                                     "http://192.168.0.2:8080/foo",
		"http://192.168.0.2:/foo":                                         "http://192.168.0.2:/foo",
		"http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:8080/foo":          "http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:8080/foo",
		"http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:/foo":              "http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:/foo",
		"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080/foo":        "http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080/foo",
		"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:/foo":            "http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:/foo",
		"http://hello.世界.com/foo":                                         "http://hello.世界.com/foo",
		"http://hello.%e4%b8%96%e7%95%8c.com/foo":                         "http://hello.世界.com/foo",
		"http://hello.%E4%B8%96%E7%95%8C.com/foo":                         "http://hello.世界.com/foo",
		"http://example.com//foo":                                         "http://example.com//foo",
		"tcp://[2020::2020:20:2020:2020%25Windows%20Loves%20Spaces]:2020": "tcp://[2020::2020:20:2020:2020%Windows Loves Spaces]:2020",
		"magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a&dn": "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a&dn",
		"mailto:?subject=hi":                                              "mailto:?subject=hi",
	} {
		t.Run(in, func(t *testing.T) {
			u, err := Parse(in)
			require.NoError(t, err)
			d := Decode(u)
			require.Equal(t, out, d)
		})
	}
}

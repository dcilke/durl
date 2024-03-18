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
	"net/url"
	"reflect"
	"testing"

	"github.com/sanity-io/litter"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		in        string
		out       *url.URL // expected parse
		roundtrip string   // expected result of reserializing the URL; empty means same as "in".
	}{
		// no path
		{
			"http://www.google.com",
			&url.URL{
				Scheme: "http",
				Host:   "www.google.com",
			},
			"",
		},
		// path
		{
			"http://www.google.com/",
			&url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/",
			},
			"",
		},
		// path with hex escaping
		{
			"http://www.google.com/file%20one%26two",
			&url.URL{
				Scheme:  "http",
				Host:    "www.google.com",
				Path:    "/file one&two",
				RawPath: "/file%20one%26two",
			},
			"",
		},
		// fragment with hex escaping
		{
			"http://www.google.com/#file%20one%26two",
			&url.URL{
				Scheme:      "http",
				Host:        "www.google.com",
				Path:        "/",
				Fragment:    "file one&two",
				RawFragment: "file%20one%26two",
			},
			"",
		},
		// user
		{
			"ftp://webmaster@www.google.com/",
			&url.URL{
				Scheme: "ftp",
				User:   url.User("webmaster"),
				Host:   "www.google.com",
				Path:   "/",
			},
			"",
		},
		// escape sequence in username
		{
			"ftp://john%20doe@www.google.com/",
			&url.URL{
				Scheme: "ftp",
				User:   url.User("john doe"),
				Host:   "www.google.com",
				Path:   "/",
			},
			"ftp://john%20doe@www.google.com/",
		},
		// empty query
		{
			"http://www.google.com/?",
			&url.URL{
				Scheme:     "http",
				Host:       "www.google.com",
				Path:       "/",
				ForceQuery: true,
			},
			"",
		},
		// query ending in question mark (Issue 14573)
		{
			"http://www.google.com/?foo=bar?",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/",
				RawQuery: "foo=bar?",
			},
			"",
		},
		// query
		{
			"http://www.google.com/?q=go+language",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/",
				RawQuery: "q=go+language",
			},
			"",
		},
		// query with hex escaping: NOT parsed
		{
			"http://www.google.com/?q=go%20language",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/",
				RawQuery: "q=go%20language",
			},
			"",
		},
		// %20 outside query
		{
			"http://www.google.com/a%20b?q=c+d",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/a b",
				RawQuery: "q=c+d",
			},
			"",
		},
		// path without leading /, so no parsing
		{
			"http:www.google.com/?q=go+language",
			&url.URL{
				Scheme:   "http",
				Opaque:   "www.google.com/",
				RawQuery: "q=go+language",
			},
			"http:www.google.com/?q=go+language",
		},
		// path without leading /, so no parsing
		{
			"http:%2f%2fwww.google.com/?q=go+language",
			&url.URL{
				Scheme:   "http",
				Opaque:   "%2f%2fwww.google.com/",
				RawQuery: "q=go+language",
			},
			"http:%2f%2fwww.google.com/?q=go+language",
		},
		// non-authority with path; see golang.org/issue/46059
		{
			"mailto:/webmaster@golang.org",
			&url.URL{
				Scheme:   "mailto",
				Path:     "/webmaster@golang.org",
				OmitHost: true,
			},
			"",
		},
		// non-authority
		{
			"mailto:webmaster@golang.org",
			&url.URL{
				Scheme: "mailto",
				Opaque: "webmaster@golang.org",
			},
			"",
		},
		// unescaped :// in query should not create a scheme
		{
			"/foo?query=http://bad",
			&url.URL{
				Path:     "/foo",
				RawQuery: "query=http://bad",
			},
			"",
		},
		// leading // without scheme should create an authority
		{
			"//foo",
			&url.URL{
				Host: "foo",
			},
			"",
		},
		// leading // without scheme, with userinfo, path, and query
		{
			"//user@foo/path?a=b",
			&url.URL{
				User:     url.User("user"),
				Host:     "foo",
				Path:     "/path",
				RawQuery: "a=b",
			},
			"",
		},
		// Three leading slashes isn't an authority, but doesn't return an error.
		// (We can't return an error, as this code is also used via
		// ServeHTTP -> ReadRequest -> Parse, which is arguably a
		// different URL parsing context, but currently shares the
		// same codepath)
		{
			"///threeslashes",
			&url.URL{
				Path: "///threeslashes",
			},
			"",
		},
		{
			"http://user:password@google.com",
			&url.URL{
				Scheme: "http",
				User:   url.UserPassword("user", "password"),
				Host:   "google.com",
			},
			"http://user:password@google.com",
		},
		// unescaped @ in username should not confuse host
		{
			"http://j@ne:password@google.com",
			&url.URL{
				Scheme: "http",
				User:   url.UserPassword("j@ne", "password"),
				Host:   "google.com",
			},
			"http://j%40ne:password@google.com",
		},
		// unescaped @ in password should not confuse host
		{
			"http://jane:p@ssword@google.com",
			&url.URL{
				Scheme: "http",
				User:   url.UserPassword("jane", "p@ssword"),
				Host:   "google.com",
			},
			"http://jane:p%40ssword@google.com",
		},
		{
			"http://j@ne:password@google.com/p@th?q=@go",
			&url.URL{
				Scheme:   "http",
				User:     url.UserPassword("j@ne", "password"),
				Host:     "google.com",
				Path:     "/p@th",
				RawQuery: "q=@go",
			},
			"http://j%40ne:password@google.com/p@th?q=@go",
		},
		{
			"http://www.google.com/?q=go+language#foo",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/",
				RawQuery: "q=go+language",
				Fragment: "foo",
			},
			"",
		},
		{
			"http://www.google.com/?q=go+language#foo&bar",
			&url.URL{
				Scheme:   "http",
				Host:     "www.google.com",
				Path:     "/",
				RawQuery: "q=go+language",
				Fragment: "foo&bar",
			},
			"http://www.google.com/?q=go+language#foo&bar",
		},
		{
			"http://www.google.com/?q=go+language#foo%26bar",
			&url.URL{
				Scheme:      "http",
				Host:        "www.google.com",
				Path:        "/",
				RawQuery:    "q=go+language",
				Fragment:    "foo&bar",
				RawFragment: "foo%26bar",
			},
			"http://www.google.com/?q=go+language#foo%26bar",
		},
		{
			"file:///home/adg/rabbits",
			&url.URL{
				Scheme: "file",
				Host:   "",
				Path:   "/home/adg/rabbits",
			},
			"file:///home/adg/rabbits",
		},
		// "Windows" paths are no exception to the rule.
		// See golang.org/issue/6027, especially comment #9.
		{
			"file:///C:/FooBar/Baz.txt",
			&url.URL{
				Scheme: "file",
				Host:   "",
				Path:   "/C:/FooBar/Baz.txt",
			},
			"file:///C:/FooBar/Baz.txt",
		},
		// case-insensitive scheme
		{
			"MaIlTo:webmaster@golang.org",
			&url.URL{
				Scheme: "mailto",
				Opaque: "webmaster@golang.org",
			},
			"mailto:webmaster@golang.org",
		},
		// Relative path
		{
			"a/b/c",
			&url.URL{
				Path: "a/b/c",
			},
			"a/b/c",
		},
		// escaped '?' in username and password
		{
			"http://%3Fam:pa%3Fsword@google.com",
			&url.URL{
				Scheme: "http",
				User:   url.UserPassword("?am", "pa?sword"),
				Host:   "google.com",
			},
			"",
		},
		// host subcomponent; IPv4 address in RFC 3986
		{
			"http://192.168.0.1/",
			&url.URL{
				Scheme: "http",
				Host:   "192.168.0.1",
				Path:   "/",
			},
			"",
		},
		// host and port subcomponents; IPv4 address in RFC 3986
		{
			"http://192.168.0.1:8080/",
			&url.URL{
				Scheme: "http",
				Host:   "192.168.0.1:8080",
				Path:   "/",
			},
			"",
		},
		// host subcomponent; IPv6 address in RFC 3986
		{
			"http://[fe80::1]/",
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1]",
				Path:   "/",
			},
			"",
		},
		// host and port subcomponents; IPv6 address in RFC 3986
		{
			"http://[fe80::1]:8080/",
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1]:8080",
				Path:   "/",
			},
			"",
		},
		// host subcomponent; IPv6 address with zone identifier in RFC 6874
		{
			"http://[fe80::1%25en0]/", // alphanum zone identifier
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1%en0]",
				Path:   "/",
			},
			"",
		},
		// host and port subcomponents; IPv6 address with zone identifier in RFC 6874
		{
			"http://[fe80::1%25en0]:8080/", // alphanum zone identifier
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1%en0]:8080",
				Path:   "/",
			},
			"",
		},
		// host subcomponent; IPv6 address with zone identifier in RFC 6874
		{
			"http://[fe80::1%25%65%6e%301-._~]/", // percent-encoded+unreserved zone identifier
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1%en01-._~]",
				Path:   "/",
			},
			"http://[fe80::1%25en01-._~]/",
		},
		// host and port subcomponents; IPv6 address with zone identifier in RFC 6874
		{
			"http://[fe80::1%25%65%6e%301-._~]:8080/", // percent-encoded+unreserved zone identifier
			&url.URL{
				Scheme: "http",
				Host:   "[fe80::1%en01-._~]:8080",
				Path:   "/",
			},
			"http://[fe80::1%25en01-._~]:8080/",
		},
		// alternate escapings of path survive round trip
		{
			"http://rest.rsc.io/foo%2fbar/baz%2Fquux?alt=media",
			&url.URL{
				Scheme:   "http",
				Host:     "rest.rsc.io",
				Path:     "/foo/bar/baz/quux",
				RawPath:  "/foo%2fbar/baz%2Fquux",
				RawQuery: "alt=media",
			},
			"",
		},
		// issue 12036
		{
			"mysql://a,b,c/bar",
			&url.URL{
				Scheme: "mysql",
				Host:   "a,b,c",
				Path:   "/bar",
			},
			"",
		},
		// worst case host, still round trips
		{
			"scheme://!$&'()*+,;=hello!:1/path",
			&url.URL{
				Scheme: "scheme",
				Host:   "!$&'()*+,;=hello!:1",
				Path:   "/path",
			},
			"",
		},
		// worst case path, still round trips
		{
			"http://host/!$&'()*+,;=:@[hello]",
			&url.URL{
				Scheme:  "http",
				Host:    "host",
				Path:    "/!$&'()*+,;=:@[hello]",
				RawPath: "/!$&'()*+,;=:@[hello]",
			},
			"",
		},
		// golang.org/issue/5684
		{
			"http://example.com/oid/[order_id]",
			&url.URL{
				Scheme:  "http",
				Host:    "example.com",
				Path:    "/oid/[order_id]",
				RawPath: "/oid/[order_id]",
			},
			"",
		},
		// golang.org/issue/12200 (colon with empty port)
		{
			"http://192.168.0.2:8080/foo",
			&url.URL{
				Scheme: "http",
				Host:   "192.168.0.2:8080",
				Path:   "/foo",
			},
			"",
		},
		{
			"http://192.168.0.2:/foo",
			&url.URL{
				Scheme: "http",
				Host:   "192.168.0.2:",
				Path:   "/foo",
			},
			"",
		},
		{
			// Malformed IPv6 but still accepted.
			"http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:8080/foo",
			&url.URL{
				Scheme: "http",
				Host:   "2b01:e34:ef40:7730:8e70:5aff:fefe:edac:8080",
				Path:   "/foo",
			},
			"",
		},
		{
			// Malformed IPv6 but still accepted.
			"http://2b01:e34:ef40:7730:8e70:5aff:fefe:edac:/foo",
			&url.URL{
				Scheme: "http",
				Host:   "2b01:e34:ef40:7730:8e70:5aff:fefe:edac:",
				Path:   "/foo",
			},
			"",
		},
		{
			"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080/foo",
			&url.URL{
				Scheme: "http",
				Host:   "[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080",
				Path:   "/foo",
			},
			"",
		},
		{
			"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:/foo",
			&url.URL{
				Scheme: "http",
				Host:   "[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:",
				Path:   "/foo",
			},
			"",
		},
		// golang.org/issue/7991 and golang.org/issue/12719 (non-ascii %-encoded in host)
		{
			"http://hello.世界.com/foo",
			&url.URL{
				Scheme: "http",
				Host:   "hello.世界.com",
				Path:   "/foo",
			},
			"http://hello.%E4%B8%96%E7%95%8C.com/foo",
		},
		{
			"http://hello.%e4%b8%96%e7%95%8c.com/foo",
			&url.URL{
				Scheme: "http",
				Host:   "hello.世界.com",
				Path:   "/foo",
			},
			"http://hello.%E4%B8%96%E7%95%8C.com/foo",
		},
		{
			"http://hello.%E4%B8%96%E7%95%8C.com/foo",
			&url.URL{
				Scheme: "http",
				Host:   "hello.世界.com",
				Path:   "/foo",
			},
			"",
		},
		// golang.org/issue/10433 (path beginning with //)
		{
			"http://example.com//foo",
			&url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "//foo",
			},
			"",
		},
		// test that we can reparse the host names we accept.
		{
			"myscheme://authority<\"hi\">/foo",
			&url.URL{
				Scheme: "myscheme",
				Host:   "authority<\"hi\">",
				Path:   "/foo",
			},
			"",
		},
		// spaces in hosts are disallowed but escaped spaces in IPv6 scope IDs are grudgingly OK.
		// This happens on Windows.
		// golang.org/issue/14002
		{
			"tcp://[2020::2020:20:2020:2020%25Windows%20Loves%20Spaces]:2020",
			&url.URL{
				Scheme: "tcp",
				Host:   "[2020::2020:20:2020:2020%Windows Loves Spaces]:2020",
			},
			"",
		},
		// test we can roundtrip magnet url
		// fix issue https://golang.org/issue/20054
		{
			"magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a&dn",
			&url.URL{
				Scheme:   "magnet",
				Host:     "",
				Path:     "",
				RawQuery: "xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a&dn",
			},
			"magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a&dn",
		},
		{
			"mailto:?subject=hi",
			&url.URL{
				Scheme:   "mailto",
				Host:     "",
				Path:     "",
				RawQuery: "subject=hi",
			},
			"mailto:?subject=hi",
		},
	} {
		u, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) returned error %v", tt.in, err)
			continue
		}
		if !reflect.DeepEqual(u, tt.out) {
			t.Errorf("Parse(%q):\n\tgot  %v\n\twant %v\n", tt.in, litter.Sdump(u), litter.Sdump(tt.out))
		}
	}
}

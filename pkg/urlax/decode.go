package urlax

import (
	"net/url"
	"strings"
)

// String reassembles the URL into a valid URL string.
// The general form of the result is one of:
//
//	scheme:opaque?query#fragment
//	scheme://userinfo@host/path?query#fragment
//
// If u.Opaque is non-empty, String uses the first form;
// otherwise it uses the second form.
// Any non-ASCII characters in host are escaped.
// To obtain the path, String uses u.EscapedPath().
//
// In the second form, the following rules apply:
//   - if u.Scheme is empty, scheme: is omitted.
//   - if u.User is nil, userinfo@ is omitted.
//   - if u.Host is empty, host/ is omitted.
//   - if u.Scheme and u.Host are empty and u.User is nil,
//     the entire scheme://userinfo@host/ is omitted.
//   - if u.Host is non-empty and u.Path begins with a /,
//     the form host/path does not add its own /.
//   - if u.RawQuery is empty, ?query is omitted.
//   - if u.Fragment is empty, #fragment is omitted.
func Decode(u *url.URL) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Opaque != "" {
		buf.WriteString(u.Opaque)
	} else {
		if u.Scheme != "" || u.Host != "" || u.User != nil {
			if u.OmitHost && u.Host == "" && u.User == nil {
				// omit empty host
			} else {
				if u.Host != "" || u.Path != "" || u.User != nil {
					buf.WriteString("//")
				}
				if ui := u.User; ui != nil {
					buf.WriteString(ui.Username())
					if p, ok := ui.Password(); ok {
						buf.WriteByte(':')
						buf.WriteString(p)
					}
					buf.WriteByte('@')
				}
				if h := u.Host; h != "" {
					buf.WriteString(h)
				}
			}
		}
		if u.Path != "" && u.Path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		// if buf.Len() == 0 {
		// 	// RFC 3986 §4.2
		// 	// A path segment that contains a colon character (e.g., "this:that")
		// 	// cannot be used as the first segment of a relative-path reference, as
		// 	// it would be mistaken for a scheme name. Such a segment must be
		// 	// preceded by a dot-segment (e.g., "./this:that") to make a relative-
		// 	// path reference.
		// 	if segment, _, _ := strings.Cut(u.Path, "/"); strings.Contains(segment, ":") {
		// 		buf.WriteString("./")
		// 	}
		// }
		buf.WriteString(u.Path)
	}
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}

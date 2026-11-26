package link

import (
	"net/http"
	"strings"
)

// Group returned by Parse, contains multiple links indexed by "rel".
type Group map[string]*Link

// Link contains a Link item with URI, Rel, and other non-URI components in Extra.
type Link struct {
	URI   string
	Rel   string
	Extra map[string]string
}

// String returns the URI.
func (l *Link) String() string {
	return l.URI
}

// ParseRequest parses the provided *http.Request into a Group.
func ParseRequest(req *http.Request) Group {
	if req == nil {
		return nil
	}

	return ParseHeader(req.Header)
}

// ParseResponse parses the provided *http.Response into a Group.
func ParseResponse(resp *http.Response) Group {
	if resp == nil {
		return nil
	}

	return ParseHeader(resp.Header)
}

// ParseHeader retrieves the Link header from the provided http.Header and parses it into a Group.
func ParseHeader(h http.Header) Group {
	var group Group

	for _, s := range h["Link"] {
		group = parse(group, s)
	}

	return group
}

// Parse parses the provided string into a Group.
//
// Malformed link-values are skipped. Parse returns nil if no link with a rel parameter was found.
func Parse(s string) Group {
	return parse(nil, s)
}

// parse adds the links found in s to group, allocating it on first use.
func parse(group Group, s string) Group {
	for s != "" {
		var link *Link

		link, s = parseLink(s)

		if link == nil || link.Rel == "" {
			continue
		}

		if group == nil {
			group = Group{}
		}

		group[link.Rel] = link
	}

	return group
}

// parseLink parses one link-value from the start of s,
// returning it (nil if malformed) and the rest of s.
func parseLink(s string) (*Link, string) {
	s = strings.TrimLeft(s, ", \t")

	if s == "" {
		return nil, ""
	}

	if s[0] != '<' {
		return nil, skipLink(s)
	}

	end := strings.IndexByte(s, '>')
	if end < 0 {
		return nil, ""
	}

	// A '<' before the '>' means this link is missing its closing '>'.
	if strings.IndexByte(s[1:end], '<') >= 0 {
		return nil, skipLink(s)
	}

	link := &Link{URI: strings.TrimSpace(s[1:end]), Extra: map[string]string{}}

	s = s[end+1:]

	for {
		s = strings.TrimLeft(s, " \t")

		if s == "" || s[0] == ',' {
			return link, s
		}

		if s[0] != ';' {
			return nil, skipLink(s)
		}

		var key, val string

		key, val, s = parseParam(s[1:])

		switch key {
		case "":
		case "rel":
			// Occurrences after the first must be ignored (RFC 8288, section 3.3).
			if link.Rel == "" {
				link.Rel = parseRel(val)
			}
		default:
			link.Extra[key] = val
		}
	}
}

// parseParam parses `key=value` or `key="value"` from the start of s.
// The returned key is empty for a parameter without a value.
func parseParam(s string) (key, val, rest string) {
	i := strings.IndexAny(s, "=;,")
	if i < 0 {
		return "", "", ""
	}

	if s[i] != '=' {
		return "", "", s[i:]
	}

	key = strings.ToLower(strings.TrimSpace(s[:i]))
	s = strings.TrimLeft(s[i+1:], " \t")

	if s != "" && s[0] == '"' {
		val, s = parseQuoted(s[1:])

		return key, val, s
	}

	j := strings.IndexAny(s, ";,")
	if j < 0 {
		j = len(s)
	}

	return key, strings.TrimSpace(s[:j]), s[j:]
}

// parseQuoted parses a quoted-string whose opening quote has already been consumed,
// returning the unescaped value and the rest of s after the closing quote.
//
// Values without backslash escapes are returned as substrings of s, without allocating.
func parseQuoted(s string) (string, string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			return unescapeQuoted(s, i)
		}
	}

	return s, ""
}

// unescapeQuoted continues parseQuoted from the first backslash at index i.
func unescapeQuoted(s string, i int) (string, string) {
	var b strings.Builder

	b.Grow(len(s))
	b.WriteString(s[:i])

	for ; i < len(s); i++ {
		switch s[i] {
		case '\\':
			if i+1 < len(s) {
				i++
				b.WriteByte(s[i])
			}
		case '"':
			return b.String(), s[i+1:]
		default:
			b.WriteByte(s[i])
		}
	}

	return b.String(), ""
}

// parseRel normalizes a rel value. Registered relation types are lowercased,
// and extension relation types (URIs) after the first are dropped.
func parseRel(val string) string {
	if !strings.ContainsAny(val, " \t") {
		return lowerRelType(val)
	}

	var b strings.Builder

	b.Grow(len(val))

	for first := true; val != ""; {
		val = strings.TrimLeft(val, " \t")

		end := strings.IndexAny(val, " \t")
		if end < 0 {
			end = len(val)
		}

		v := val[:end]
		val = val[end:]

		if v == "" || !first && strings.HasPrefix(v, "http") {
			continue
		}

		if !first {
			b.WriteByte(' ')
		}

		b.WriteString(lowerRelType(v))

		first = false
	}

	return b.String()
}

// lowerRelType lowercases registered relation types, while extension relation types (URIs) keep their case.
func lowerRelType(v string) string {
	if strings.IndexByte(v, ':') >= 0 {
		return v
	}

	return strings.ToLower(v)
}

// skipLink skips to the next link-value separator, ignoring commas inside quoted strings.
func skipLink(s string) string {
	quoted := false

	for i := 0; i < len(s); i++ {
		switch {
		case quoted && s[i] == '\\':
			i++
		case s[i] == '"':
			quoted = !quoted
		case !quoted && s[i] == ',':
			return s[i:]
		}
	}

	return ""
}

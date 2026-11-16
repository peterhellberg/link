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
	if headers, found := h["Link"]; found {
		return Parse(strings.Join(headers, ", "))
	}

	return nil
}

// Parse parses the provided string into a Group.
func Parse(s string) Group {
	group := Group{}

	for s != "" {
		var l *Link

		l, s = parseLink(s)

		if l != nil && l.Rel != "" {
			group[l.Rel] = l
		}
	}

	if len(group) == 0 {
		return nil
	}

	return group
}

// parseLink parses one link-value from the start of s, returning it (nil if malformed) and the rest of s.
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
			link.Rel = parseRel(val)
		default:
			link.Extra[key] = val
		}
	}
}

// parseParam parses `key=value` or `key="value"` from the start of s; key is empty for a param without a value.
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

// parseQuoted reads a quoted-string body (opening quote already consumed), handling backslash escapes.
func parseQuoted(s string) (string, string) {
	var b strings.Builder

	for i := 0; i < len(s); i++ {
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

// parseRel lowercases registered relation types and drops extension (URI) types after the first.
func parseRel(val string) string {
	var rels []string

	for i, v := range strings.Fields(val) {
		if i > 0 && strings.HasPrefix(v, "http") {
			continue
		}

		if !strings.Contains(v, ":") {
			v = strings.ToLower(v)
		}

		rels = append(rels, v)
	}

	return strings.Join(rels, " ")
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

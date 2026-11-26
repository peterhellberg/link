package link

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLinkString(t *testing.T) {
	t.Parallel()

	l := Parse(`<https://example.com/?page=2>; rel="next"; title="foo"`)["next"]

	if got, want := l.String(), "https://example.com/?page=2"; got != want {
		t.Fatalf(`l.String() = %q, want %q`, got, want)
	}
}

func TestParseRequest(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "", nil)
	req.Header.Set("Link", `<https://example.com/?page=2>; rel="next"`)

	g := ParseRequest(req)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["next"] == nil {
		t.Fatalf(`g["next"] == nil`)
	}

	if got, want := g["next"].URI, "https://example.com/?page=2"; got != want {
		t.Fatalf(`g["next"].URI = %q, want %q`, got, want)
	}

	if got, want := g["next"].Rel, "next"; got != want {
		t.Fatalf(`g["next"].Rel = %q, want %q`, got, want)
	}

	if got, want := len(ParseRequest(nil)), 0; got != want {
		t.Fatalf(`len(ParseRequest(nil)) = %d, want %d`, got, want)
	}
}

func TestParseResponse(t *testing.T) {
	t.Parallel()

	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Link", `<https://example.com/?page=2>; rel="next"`)

	g := ParseResponse(resp)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["next"] == nil {
		t.Fatalf(`g["next"] == nil`)
	}

	if got, want := g["next"].URI, "https://example.com/?page=2"; got != want {
		t.Fatalf(`g["next"].URI = %q, want %q`, got, want)
	}

	if got, want := g["next"].Rel, "next"; got != want {
		t.Fatalf(`g["next"].Rel = %q, want %q`, got, want)
	}

	if got, want := len(ParseResponse(nil)), 0; got != want {
		t.Fatalf(`len(ParseResponse(nil)) = %d, want %d`, got, want)
	}
}

func TestParseHeader_single(t *testing.T) {
	t.Parallel()

	h := http.Header{}
	h.Set("Link", `<https://example.com/?page=2>; rel="next"`)

	g := ParseHeader(h)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["next"] == nil {
		t.Fatalf(`g["next"] == nil`)
	}

	if got, want := g["next"].URI, "https://example.com/?page=2"; got != want {
		t.Fatalf(`g["next"].URI = %q, want %q`, got, want)
	}

	if got, want := g["next"].Rel, "next"; got != want {
		t.Fatalf(`g["next"].Rel = %q, want %q`, got, want)
	}
}

func TestParseHeader_multiple(t *testing.T) {
	t.Parallel()

	h := http.Header{}
	h.Add("Link", `<https://example.com/?page=2>; rel="next",<https://example.com/?page=34>; rel="last"`)

	g := ParseHeader(h)

	if got, want := len(g), 2; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["next"] == nil {
		t.Fatalf(`g["next"] == nil`)
	}

	if got, want := g["next"].URI, "https://example.com/?page=2"; got != want {
		t.Fatalf(`g["next"].URI = %q, want %q`, got, want)
	}

	if got, want := g["next"].Rel, "next"; got != want {
		t.Fatalf(`g["next"].Rel = %q, want %q`, got, want)
	}

	if g["last"] == nil {
		t.Fatalf(`g["last"] == nil`)
	}

	if got, want := g["last"].URI, "https://example.com/?page=34"; got != want {
		t.Fatalf(`g["last"].URI = %q, want %q`, got, want)
	}

	if got, want := g["last"].Rel, "last"; got != want {
		t.Fatalf(`g["last"].Rel = %q, want %q`, got, want)
	}
}

func TestParseHeader_multiple_headers(t *testing.T) {
	t.Parallel()

	h := http.Header{}
	h.Add("Link", `<https://example.com/?page=2>; rel="next",<https://example.com/?page=34>; rel="last"`)
	h.Add("Link", `<https://example.com/?page=foo>; rel="foo",<https://example.com/?page=bar>; rel="bar"`)

	g := ParseHeader(h)

	if got, want := len(g), 4; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["foo"] == nil {
		t.Fatalf(`g["foo"] == nil`)
	}

	if got, want := g["bar"].URI, "https://example.com/?page=bar"; got != want {
		t.Fatalf(`g["bar"].URI = %q, want %q`, got, want)
	}

	if got, want := g["next"].Rel, "next"; got != want {
		t.Fatalf(`g["next"].Rel = %q, want %q`, got, want)
	}

	if g["last"] == nil {
		t.Fatalf(`g["last"] == nil`)
	}

	if got, want := g["last"].URI, "https://example.com/?page=34"; got != want {
		t.Fatalf(`g["last"].URI = %q, want %q`, got, want)
	}

	if got, want := g["last"].Rel, "last"; got != want {
		t.Fatalf(`g["last"].Rel = %q, want %q`, got, want)
	}
}

func TestParseHeader_extra(t *testing.T) {
	t.Parallel()

	h := http.Header{}
	h.Add("Link", `<https://example.com/?page=2>; rel="next"; title="foo"`)

	g := ParseHeader(h)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["next"] == nil {
		t.Fatalf(`g["next"] == nil`)
	}

	if got, want := g["next"].Extra["title"], "foo"; got != want {
		t.Fatalf(`g["next"].Extra["title"] = %q, want %q`, got, want)
	}
}

func TestParseHeader_noLink(t *testing.T) {
	t.Parallel()

	if ParseHeader(http.Header{}) != nil {
		t.Fatalf(`Parse(http.Header{}) != nil`)
	}
}

func TestParseHeader_nilHeader(t *testing.T) {
	t.Parallel()

	if ParseHeader(nil) != nil {
		t.Fatalf(`ParseHeader(nil) != nil`)
	}
}

func TestParse_emptyString(t *testing.T) {
	t.Parallel()

	if Parse("") != nil {
		t.Fatalf(`Parse("") != nil`)
	}
}

func TestParse_valuesWithComma(t *testing.T) {
	t.Parallel()

	g := Parse(`<//www.w3.org/wiki/LinkHeader>; rel="original latest-version",<//www.w3.org/wiki/Special:TimeGate/LinkHeader>; rel="timegate",<//www.w3.org/wiki/Special:TimeMap/LinkHeader>; rel="timemap"; type="application/link-format"; from="Mon, 03 Sep 2007 14:52:48 GMT"; until="Tue, 16 Jun 2015 22:59:23 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152>; rel="first memento"; datetime="Mon, 03 Sep 2007 14:52:48 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=84697>; rel="last memento"; datetime="Tue, 16 Jun 2015 22:59:23 GMT"`)

	if got, want := len(g), 5; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if got, want := g["original latest-version"].URI, "//www.w3.org/wiki/LinkHeader"; got != want {
		t.Fatalf(`g["original latest-version"].URI = %q, want %q`, got, want)
	}

	if got, want := g["last memento"].Extra["datetime"], "Tue, 16 Jun 2015 22:59:23 GMT"; got != want {
		t.Fatalf(`g["last memento"].Extra["datetime"] = %q, want %q`, got, want)
	}
}

func TestParse_rfc5988Example1(t *testing.T) {
	t.Parallel()

	g := Parse(`<http://example.com/TheBook/chapter2>; rel="previous"; title="previous chapter"`)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["previous"] == nil {
		t.Fatalf(`g["previous"] == nil`)
	}

	if got, want := g["previous"].Extra["title"], "previous chapter"; got != want {
		t.Fatalf(`g["previous"].Extra["title"] = %q, want %q`, got, want)
	}
}

func TestParse_rfc5988Example2(t *testing.T) {
	t.Parallel()

	g := Parse(`</>; rel="http://example.net/foo"`)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["http://example.net/foo"] == nil {
		t.Fatalf(`g["http://example.net/foo"] == nil`)
	}

	l := g["http://example.net/foo"]

	if got, want := l.URI, "/"; got != want {
		t.Fatalf(`l.URI = %q, want %q`, got, want)
	}
}

func TestParse_rfc5988Example3(t *testing.T) {
	t.Parallel()

	// Extended notation is not supported yet
	// g := Parse(`</TheBook/chapter2>; rel="previous"; title*=UTF-8'de'letztes%20Kapitel, </TheBook/chapter4>; rel="next"; title*=UTF-8'de'n%c3%a4chstes%20Kapitel`)
}

func TestParse_rfc5988Example4(t *testing.T) {
	t.Parallel()

	// Extension relation types are ignored for now
	g := Parse(`<http://example.org/>; rel="start http://example.net/relation/other"`)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["start"] == nil {
		t.Fatalf(`g["start"] == nil`)
	}

	if got, want := g["start"].URI, "http://example.org/"; got != want {
		t.Fatalf(`g["start"].URI = %q, want %q`, got, want)
	}
}

func TestParse_noQuotesExample(t *testing.T) {
	t.Parallel()

	g := Parse(`</example.png>;rel=preload;as=image`)

	if got, want := len(g), 1; got != want {
		t.Fatalf(`len(g) = %d, want %d`, got, want)
	}

	if g["preload"] == nil {
		t.Fatalf(`g["preload"] == nil`)
	}

	l := *g["preload"]

	if got, want := l.URI, "/example.png"; got != want {
		t.Fatalf("l.URI = %q, want = %q", got, want)
	}

	if got, want := l.Rel, "preload"; got != want {
		t.Fatalf("l.Rel = %q, want = %q", got, want)
	}

	if l.Extra == nil {
		t.Fatalf("l.Extra == nil")
	}

	if got, want := l.Extra["as"], "image"; got != want {
		t.Fatalf("l.Extra[\"as\"] = %q, want %q", got, want)
	}
}

func TestParse_edgeCases(t *testing.T) {
	t.Parallel()

	type want struct{ rel, uri, key, val string }

	for _, tt := range []struct {
		name  string
		input string
		want  []want
	}{
		{
			name:  "multiple links with unquoted values",
			input: `<https://example.com/?page=2>; rel=next, <https://example.com/?page=34>; rel=last`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2"}, {rel: "last", uri: "https://example.com/?page=34"}},
		},
		{
			name:  "comma in URI",
			input: `<https://example.com/?ids=1,2>; rel="next"`,
			want:  []want{{rel: "next", uri: "https://example.com/?ids=1,2"}},
		},
		{
			name:  "comma in quoted value",
			input: `<https://example.com/>; rel="memento"; datetime="Mon, 03 Sep 2007 14:52:48 GMT"`,
			want:  []want{{rel: "memento", uri: "https://example.com/", key: "datetime", val: "Mon, 03 Sep 2007 14:52:48 GMT"}},
		},
		{
			name:  "malformed link does not drop the others",
			input: `garbage "with, comma", <https://example.com/?page=2>; rel="next"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2"}},
		},
		{
			name:  "trailing garbage after param drops only that link",
			input: `<https://example.com/a>; rel="a" junk, <https://example.com/b>; rel="b"`,
			want:  []want{{rel: "b", uri: "https://example.com/b"}},
		},
		{
			name:  "missing closing angle bracket drops only that link",
			input: `<https://example.com/a; rel="a", <https://example.com/b>; rel="b"`,
			want:  []want{{rel: "b", uri: "https://example.com/b"}},
		},
		{
			name:  "link without params next to one with params",
			input: `<https://example.com/a>, <https://example.com/?page=2>; rel="next"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2"}},
		},
		{
			name:  "whitespace before semicolon",
			input: `<https://example.com/?page=2> ; rel="next"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2"}},
		},
		{
			name:  "uppercase param names",
			input: `<https://example.com/?page=2>; REL="next"; Title="foo"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2", key: "title", val: "foo"}},
		},
		{
			name:  "uppercase rel value",
			input: `<https://example.com/?page=2>; rel="Next"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2"}},
		},
		{
			name:  "extension rel value keeps its case",
			input: `</>; rel="http://Example.net/Foo"`,
			want:  []want{{rel: "http://Example.net/Foo", uri: "/"}},
		},
		{
			name:  "multiple rel values are normalized",
			input: `</>; rel="  First   Memento "`,
			want:  []want{{rel: "first memento", uri: "/"}},
		},
		{
			name:  "only the first rel param is used",
			input: `</>; rel="next"; rel="last"`,
			want:  []want{{rel: "next", uri: "/"}},
		},
		{
			name:  "escaped quote in quoted value",
			input: `<https://example.com/?page=2>; rel="next"; title="say \"hi\"; ok"`,
			want:  []want{{rel: "next", uri: "https://example.com/?page=2", key: "title", val: `say "hi"; ok`}},
		},
		{
			name:  "param without value is ignored",
			input: `</>; crossorigin; rel=preload`,
			want:  []want{{rel: "preload", uri: "/"}},
		},
		{
			name:  "no usable links",
			input: `garbage, <https://example.com/>; title="no rel"`,
		},
	} {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := Parse(tt.input)

			if got, want := len(g), len(tt.want); got != want {
				t.Fatalf(`len(g) = %d, want %d (g = %v)`, got, want, g)
			}

			for _, w := range tt.want {
				l := g[w.rel]

				if l == nil {
					t.Fatalf(`g[%q] == nil (g = %v)`, w.rel, g)
				}

				if l.URI != w.uri {
					t.Fatalf(`g[%q].URI = %q, want %q`, w.rel, l.URI, w.uri)
				}

				if l.Rel != w.rel {
					t.Fatalf(`g[%q].Rel = %q, want %q`, w.rel, l.Rel, w.rel)
				}

				if w.key != "" && l.Extra[w.key] != w.val {
					t.Fatalf(`g[%q].Extra[%q] = %q, want %q`, w.rel, w.key, l.Extra[w.key], w.val)
				}
			}
		})
	}
}

func TestParse_fuzzCrashers(t *testing.T) {
	t.Parallel()

	Parse("0")
}

func ExampleParse() {
	l := Parse(`<https://example.com/?page=2>; rel="next"; title="foo"`)["next"]

	fmt.Printf("URI: %q, Rel: %q, Extra: %+v\n", l.URI, l.Rel, l.Extra)
	// Output: URI: "https://example.com/?page=2", Rel: "next", Extra: map[title:foo]
}

func FuzzParse(f *testing.F) {
	f.Add(`<https://example.com/?page=2>; rel="next"; title="foo"`)
	f.Add(`<http://example.org/>; rel="start http://example.net/relation/other"`)

	f.Fuzz(func(t *testing.T, s string) {
		Parse(s)
	})
}

func BenchmarkParse(b *testing.B) {
	for _, bb := range []struct {
		name  string
		input string
	}{
		{"single", `<https://example.com/?page=2>; rel="next"`},
		{"extra", `<https://example.com/?page=2>; rel="next"; title="foo"`},
		{"multiple", `<https://example.com/?page=2>; rel="next", <https://example.com/?page=34>; rel="last"`},
		{"unquoted", `</example.png>;rel=preload;as=image`},
		{"extension", `<http://example.org/>; rel="start http://example.net/relation/other"`},
		{"valuesWithComma", `<//www.w3.org/wiki/LinkHeader>; rel="original latest-version",<//www.w3.org/wiki/Special:TimeGate/LinkHeader>; rel="timegate",<//www.w3.org/wiki/Special:TimeMap/LinkHeader>; rel="timemap"; type="application/link-format"; from="Mon, 03 Sep 2007 14:52:48 GMT"; until="Tue, 16 Jun 2015 22:59:23 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=10152>; rel="first memento"; datetime="Mon, 03 Sep 2007 14:52:48 GMT",<//www.w3.org/wiki/index.php?title=LinkHeader&oldid=84697>; rel="last memento"; datetime="Tue, 16 Jun 2015 22:59:23 GMT"`},
	} {
		b.Run(bb.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				Parse(bb.input)
			}
		})
	}
}

func BenchmarkParseHeader(b *testing.B) {
	h := http.Header{}
	h.Add("Link", `<https://example.com/?page=2>; rel="next",<https://example.com/?page=34>; rel="last"`)
	h.Add("Link", `<https://example.com/?page=foo>; rel="foo",<https://example.com/?page=bar>; rel="bar"`)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseHeader(h)
	}
}

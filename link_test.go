package link

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLinkString(t *testing.T) {
	l := Parse(`<https://example.com/?page=2>; rel="next"; title="foo"`)["next"]

	if got, want := l.String(), "https://example.com/?page=2"; got != want {
		t.Fatalf(`l.String() = %q, want %q`, got, want)
	}
}

func TestParseRequest(t *testing.T) {
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
	h := http.Header{}
	h.Set("Link", "<https://example.com/?page=2>; rel=\"next\"")

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
	h := http.Header{}
	h.Add("Link", "<https://example.com/?page=2>; rel=\"next\", <https://example.com/?page=34>; rel=\"last\"")

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

func TestParseHeader_extra(t *testing.T) {
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
	if ParseHeader(http.Header{}) != nil {
		t.Fatalf(`Parse(http.Header{}) != nil`)
	}
}

func ExampleParse() {
	l := Parse(`<https://example.com/?page=2>; rel="next"; title="foo"`)["next"]

	fmt.Printf("URI: %q, Rel: %q, Extra: %+v\n", l.URI, l.Rel, l.Extra)
	// Output: URI: "https://example.com/?page=2", Rel: "next", Extra: map[title:foo]
}

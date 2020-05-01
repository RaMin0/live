package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractText(t *testing.T) {
	cases := []struct {
		name string
		a    string
		text string
	}{
		{
			name: "valid",
			a:    `<a href="/login">Login</a>`,
			text: "Login",
		},
		{
			name: "valid: nested",
			a:    `<a href="/login">Login <span>as <strong>Admin<strong></span></a>`,
			text: "Login as Admin",
		},
		{
			name: "valid: comments",
			a:    `<a href="/login">Login <!-- This is a comment --></a>`,
			text: "Login",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := parse(t, c.a)
			text := extractText(a)
			if text != c.text {
				t.Fatalf("expected %q, got %q", c.text, text)
			}
		})
	}
}

func TestExtractHref(t *testing.T) {
	cases := []struct {
		name string
		a    string
		href string
	}{
		{
			name: "valid",
			a:    `<a href="/login">Login</a>`,
			href: "/login",
		},
		{
			name: "missing href",
			a:    `<a>Login</a>`,
			href: "",
		},
		{
			name: "other attrs",
			a:    `<a class="link">Login</a>`,
			href: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := parse(t, c.a)
			href := extractHref(a)
			if href != c.href {
				t.Fatalf("expected %q, got %q", c.href, href)
			}
		})
	}
}

func parse(t *testing.T, a string) *html.Node {
	n, err := html.Parse(strings.NewReader(a))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
		return nil
	}
	return n.FirstChild.FirstChild.NextSibling.FirstChild
}

package main

import (
	"github.com/pintjuk/faas/function"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

func md2html(md string) (html string, err error) {
	unsafe := blackfriday.Run([]byte(md))
	html = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	return
}

func main() {
	function.RunFunc([]string{"md"}, md2html, ":8080")
}

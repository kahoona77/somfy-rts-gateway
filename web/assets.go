package web

import (
	"embed"
	"io/fs"
)

//go:embed tmpl/*.html static/*
var embeddedFiles embed.FS

func StaticFS() fs.FS {
	return mustSub("static")
}

func mustSub(path string) fs.FS {
	sub, err := fs.Sub(embeddedFiles, path)
	if err != nil {
		panic(err)
	}
	return sub
}

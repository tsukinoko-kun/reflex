package public

import (
	"embed"
	"io"
	"mime"
	"net/http"
	unixpath "path"
)

//go:embed content/*
var efs embed.FS

const sniffLen = 512

func Handler(w http.ResponseWriter, r *http.Request) bool {
	name := unixpath.Join("content", r.URL.Path)
	f, err := efs.Open(name)
	if err != nil {
		return false
	}

	ctype := mime.TypeByExtension(unixpath.Ext(r.URL.Path))
	if ctype == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [sniffLen]byte
		n, _ := io.ReadFull(f, buf[:])
		ctype = http.DetectContentType(buf[:n])
		_ = f.Close()
		f, err = efs.Open(name)
		if err != nil {
			return false
		}
	}
	w.Header().Set("Content-Type", ctype)

	// cache for 4 hours
	w.Header().Set("Cache-Control", "public, max-age=14400")

	if _, err := io.Copy(w, f); err != nil {
		_ = f.Close()
		return false
	}

	_ = f.Close()
	return true
}

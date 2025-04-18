package static

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	//go:embed content/*
	content embed.FS

	//go:embed frontend/*
	frontend embed.FS

	//go:embed backend/*
	backend embed.FS
)

func WriteContentTo(out string) error {
	return writeTo(content, "content", out)
}

func WriteFrontendTo(out string) error {
	return writeTo(frontend, "frontend", out)
}

func WriteBackendTo(out string) error {
	return writeTo(backend, "backend", out)
}

func writeTo(efs fs.FS, root string, out string) error {
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}

	return fs.WalkDir(efs, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fullPath := filepath.Join(out, path[len(root)+1:])

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer f.Close()

		ef, err := efs.Open(path)
		if err != nil {
			return err
		}
		defer ef.Close()

		if _, err := io.Copy(f, ef); err != nil {
			return err
		}

		return nil
	})
}

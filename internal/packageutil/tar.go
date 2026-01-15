package packageutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func CreateTarGz(outPath string, workDir string, include []string) error {

	file, err := os.Create(outPath)

	if err != nil {
		return err
	}

	defer file.Close()

	gw := gzip.NewWriter(file)

	defer gw.Close()

	tw := tar.NewWriter(gw)

	defer tw.Close()

	for _, p := range include {
		abs := filepath.Join(workDir, p)
		if err := addPath(tw, abs, p); err != nil {
			return err
		}
	}
	return nil
}

func addPath(tw *tar.Writer, absPath string, relPath string) error {
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return filepath.Walk(absPath, func(path string, fi os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if fi.IsDir() {
				return nil
			}
			subRel, err := filepath.Rel(absPath, path)
			if err != nil {
				return err
			}
			entry := filepath.ToSlash(filepath.Join(relPath, subRel))
			return addFile(tw, path, entry, fi)
		})
	}

	return addFile(tw, absPath, filepath.ToSlash(relPath), info)
}

func addFile(tw *tar.Writer, filePath, entryName string, fi os.FileInfo) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	hdr.Name = entryName

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	_, err = io.Copy(tw, f)
	return err
}

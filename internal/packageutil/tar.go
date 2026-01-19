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

	gzipWriter := gzip.NewWriter(file)

	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)

	defer tarWriter.Close()

	for _, p := range include {
		filePath := filepath.Join(workDir, p)
		if err := addPath(tarWriter, filePath, p); err != nil {
			return err
		}
	}
	return nil
}

func addPath(tarWriter *tar.Writer, basePath string, relPath string) error {
	fileInfo, err := os.Stat(basePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return filepath.Walk(basePath, func(path string, fileInfo os.FileInfo, walkError error) error {
			if walkError != nil {
				return walkError
			}
			if fileInfo.IsDir() {
				return nil
			}
			subRel, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}
			entry := filepath.ToSlash(filepath.Join(relPath, subRel))
			return addFile(tarWriter, path, entry, fileInfo)
		})
	}

	return addFile(tarWriter, basePath, filepath.ToSlash(relPath), fileInfo)
}

func addFile(tarWriter *tar.Writer, filePath, entryName string, fi os.FileInfo) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	header.Name = entryName

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, f)
	return err
}

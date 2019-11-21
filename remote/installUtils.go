package remote

// Transferred from cmd/internal/install/*

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"
	"path"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Expects "installing" dir. to be present!!!
func download(url string) (*os.File, error) {
	file, err := ioutil.TempFile("installing", "mod")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return file, err
}

func copyDir(dst, src string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcStat.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFile := path.Join(src, file.Name())
		dstFile := path.Join(dst, file.Name())

		if file.IsDir() {
			err = copyDir(dstFile, srcFile)
		} else {
			err = copyFile(dstFile, srcFile)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcStat.Mode())
}

func extract(file *os.File) (string, error) {
	dir, err := ioutil.TempDir("installing", "mod")
	if err != nil {
		return "", err
	}

	reader, err := zip.OpenReader(file.Name())
	if err != nil {
		return "", err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dir, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dir)+string(os.PathSeparator)) {
			return dir, fmt.Errorf("%s: illegal file path", fpath)
		}

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return dir, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return dir, err
		}

		tmpFile, err := file.Open()
		if err != nil {
			return dir, err
		}

		_, err = io.Copy(outFile, tmpFile)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		tmpFile.Close()

		if err != nil {
			return dir, err
		}
	}
	return dir, nil
}

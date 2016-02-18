package lang

import (
	"io"
	"os"
	"strings"
)

func fileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func dirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func ensureDir(path string) error {
	if !dirExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func copy(dest io.Writer, src io.Reader, progress func(p int64)) (written int64, err error) {

	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)

		if nr > 0 {
			nw, ew := dest.Write(buf[0:nr])

			if nw > 0 {
				written += int64(nw)
				if progress != nil {
					progress(written)
				}

			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}

		if er != nil {
			err = er
			break
		}

	}

	return written, err
}

func copyFile(source string, dest string, progressCB func(progress int64)) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()
	//sourceinfo, _ := sourcefile.Stat()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = copy(destfile, sourcefile, progressCB)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func copyDir(source string, dest string, progressCB func(source, dest string)) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(sourcefilepointer, destinationfilepointer, progressCB)
			if err != nil {
				return err
			}
		} else {
			// perform copy

			err = copyFile(sourcefilepointer, destinationfilepointer, nil)
			if err != nil {
				return err
			}
			if progressCB != nil {
				progressCB(sourcefilepointer, destinationfilepointer)
			}

		}

	}
	return
}

func analyzeDir(source string) (number_of_files int64, total_size int64) {

	// get properties of source dir
	_, err := os.Stat(source)
	if err != nil {
		return -1, -1
	}

	// create dest dir

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		if obj.IsDir() {
			n, t := analyzeDir(sourcefilepointer)
			number_of_files += n
			total_size += t
		} else {
			number_of_files++
			total_size += obj.Size()
		}

	}
	return
}

type StrSlice []string

func (self StrSlice) Contains(str string) bool {
	for _, s := range self {
		if s == str {
			return true
		}
	}
	return false
}

func (self StrSlice) Join(sep string) string {
	return strings.Join(self, sep)
}

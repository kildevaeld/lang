package lang

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var gzipReg = regexp.MustCompile(".*\\.tar\\.gz$")
var zipReg = regexp.MustCompile(".*\\.zip$")
var bzipReg = regexp.MustCompile(".*\\.tar\\.bz2$")
var pkgReg = regexp.MustCompile(".*\\.pkg$")

func unpack(source, target string, progressCB func(progress, total int64)) error {

	file, ferr := os.Open(source)

	if ferr != nil {
		return ferr
	}
	defer file.Close()

	if gzipReg.Match([]byte(source)) {
		reader, rerr := gzip.NewReader(file)

		if rerr != nil {
			return rerr
		}
		defer reader.Close()

		total := analyze_tar(reader)

		file.Seek(0, 0)
		reader.Reset(file)

		return unpack_tar(reader, target, func(progress int64) {
			progressCB(progress, total)
		})

	} else if bzipReg.Match([]byte(source)) {
		reader := bzip2.NewReader(file)

		return unpack_tar(reader, target, nil)
	} else if pkgReg.Match([]byte(source)) {

	}

	/*file, ferr := os.Open(localFile)

	if ferr != nil {
		return ferr
	}
	defer file.Close()

	reader, rerr := gzip.NewReader(file)
	if rerr != nil {
		return rerr
	}
	defer reader.Close()
	progressCB(Unpack)

	target := self.config.Source

	if version.Source {
		target = self.config.Temp
	}

	err = UnpackFile(reader, target, 0)

	if err != nil {
		return err
	}*/
	return nil
}

func analyze_tar(reader io.Reader) int64 {
	var out int64
	tarball := tar.NewReader(reader)

	for {
		_, err := tarball.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0
		}
		out++

	}
	return out
}

func unpack_tar(reader io.Reader, target string, progressCB func(progress int64)) error {
	tarball := tar.NewReader(reader)
	var number int64 = 0
	for {
		header, err := tarball.Next()

		if err == io.EOF {

			break
		}

		if err != nil {
			return err
		}

		filename := header.Name
		filename = filepath.Join(target, filename)
		filename, _ = filepath.Abs(filename)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory

			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				return err
			}

		case tar.TypeReg:
			// handle normal file

			writer, err := os.Create(filename)

			if err != nil {
				return err
			}

			io.Copy(writer, tarball)
			writer.Close()
			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

		case tar.TypeSymlink:

			dirname := filepath.Dir(filename)

			basename := filepath.Base(filename)

			cur, err := os.Getwd()

			if err != nil {
				return err
			}

			os.Chdir(dirname)

			err = os.Symlink(header.Linkname, basename)
			if err != nil {
				return err
			}

			os.Chdir(cur)

		default:
			return fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
		number++
		if progressCB != nil {
			progressCB(number)
		}

	}

	return nil
}

package lang

import (
	"net/http"
	"os"
)

func download(url, target string, progressCB func(step Step, progress, total int64)) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	file, err := os.Create(target)

	if err != nil {
		return err
	}

	defer file.Close()
	defer res.Body.Close()

	_, err = copy(file, res.Body, func(progress int64) {
		if progressCB != nil {
			progressCB(Download, progress, res.ContentLength)
		}
	})

	return err
}

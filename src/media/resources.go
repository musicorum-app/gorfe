package media

import (
	"github.com/fogleman/gg"
	"gorfe/utils"
	"image"
	"io"
	"net/http"
	"os"
)

func GetImage (url string) (image.Image, error) {
	file := utils.Hash(url) + ".jpg"
	path := utils.Config.MediaPath + file

	img, err := gg.LoadImage(path)

	if err == nil {
		return img, nil
	}

	err = DownloadFile(path, url)
	if err != nil {
		return nil, err
	}

	img, err = gg.LoadImage(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func DownloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

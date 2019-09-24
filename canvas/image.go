package canvas

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

var fontCacheMap = map[string]*font.Face{}
var imageCacheMap = map[string]*image.Image{}

// FetchRemoteImage 远程图片缓存在 boltDB 中
func FetchRemoteImage(link string, width uint, height uint, expired time.Duration) (image.Image, error) {
	key := fmt.Sprintf("remote:%s", link)
	if cache := GetCache(key); cache != nil {
		img, _, err := image.Decode(bytes.NewReader(cache))
		if err == nil {
			return ImageResize(img, width, height), nil
		}

		log.Printf("warning: cache-image-decode-error: %s\n", err)
	}

	reader, _, err := Download(link)
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}

	if err := SaveCache(key, buffer, expired); err != nil {
		log.Printf("warning: save-image-cache-error: %s\n", err)
	}

	return ImageResize(img, width, height), nil
}

// LoadLocalImage 本地图片缓存在内存中
func LoadLocalImage(imagePath string) (image.Image, error) {
	if img, cached := imageCacheMap[imagePath]; cached {
		return *img, nil
	}

	img, err := gg.LoadImage(imagePath)
	if err != nil {
		return nil, err
	}

	imageCacheMap[imagePath] = &img
	return img, nil
}

// LoadLocalFont 字体缓存在内存中
func LoadLocalFont(fontPath string, fontSize float64) (font.Face, error) {
	fontKey := fmt.Sprintf("%s:%f", fontPath, fontSize)
	if font, cached := fontCacheMap[fontKey]; cached {
		return *font, nil
	}

	fontPath = path.Join(fontPath)
	font, err := gg.LoadFontFace(fontPath, fontSize)
	if err != nil {
		return nil, err
	}

	fontCacheMap[fontKey] = &font
	return font, nil
}

// ImageResize 图片缩放
func ImageResize(input image.Image, width uint, height uint) image.Image {
	bounds := input.Bounds()
	if bounds.Dx() != int(width) || bounds.Dy() != int(height) {
		return resize.Resize(width, height, input, resize.Lanczos3)
	}

	return input
}

// ImageRound 图片变圆
func ImageRound(input image.Image) image.Image {
	size := (input.Bounds().Dx() + input.Bounds().Dy()) / 2
	ctx := gg.NewContext(size, size)
	ctx.DrawRoundedRectangle(0, 0, float64(size), float64(size), float64(size/2))
	ctx.Clip()
	ctx.DrawImage(input, 0, 0)
	return ctx.Image()
}

package img

import (
	"github.com/maczh/logs"
	"github.com/nfnt/resize"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"image/color"
	"image/draw"
)

func StitchImage(width, height int, backgroundUrl string, template []image.Rectangle, imageUrls []string, borderSize int, borderColor color.Color) (image.Image, error) {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	if backgroundUrl != "" {
		bg, err := impexp.NewURLImporter(backgroundUrl).Import()
		if err != nil {
			logs.Error("背景图下载错误:{}", err.Error())
			return nil, err
		}
		bgResized := resize.Thumbnail(uint(width), uint(height), bg, resize.Lanczos3)
		draw.Draw(canvas, canvas.Bounds(), bgResized, image.ZP, draw.Src)
	} else {
		draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.RGBA{255, 255, 255, 255}), image.ZP, draw.Src)
	}
	for i, temp := range template {
		//如果有边框，则绘制边框底图
		if borderSize > 0 {
			borderRect := image.Rect(temp.Min.X-borderSize, temp.Min.Y-borderSize, temp.Max.X+borderSize, temp.Max.Y+borderSize)
			draw.Draw(canvas, borderRect, image.NewUniform(borderColor), image.ZP, draw.Src)
		}
		//填充图片
		img, err := impexp.NewURLImporter(imageUrls[i]).Import()
		if err != nil {
			logs.Error("第{}张图片{}下载错误", i, imageUrls[i])
			continue
		}
		xr := float64(img.Bounds().Dx()) / float64(temp.Dx())
		yr := float64(img.Bounds().Dy()) / float64(temp.Dy())
		maxWidth := 0
		maxHeight := 0
		logs.Debug("第{}张图片原始尺寸:{}x{}", i, img.Bounds().Dx(), img.Bounds().Dy())
		if xr == yr {
			maxWidth = temp.Dx()
			maxHeight = temp.Dy()
		}
		cropRect := image.Rect(0, 0, temp.Dx(), temp.Dy())
		var imgResized, imgCroped image.Image
		if xr < 1.0 || yr < 1.0 {
			maxWidth = img.Bounds().Dx()
			maxHeight = img.Bounds().Dy()
			imgResized = img
			if img.Bounds().Dx() > temp.Dx() {
				outx := (maxWidth - temp.Dx()) / 2
				cropRect.Min.X = outx
				cropRect.Max.X = cropRect.Min.X + temp.Dx()
			}
			if img.Bounds().Dy() > temp.Dy() {
				outy := (maxHeight - temp.Dy()) / 2
				cropRect.Min.Y = outy
				cropRect.Max.Y = cropRect.Min.Y + temp.Dy()
			}
			if img.Bounds().Dx() < temp.Dx() {
				cropRect.Max.X = img.Bounds().Max.X
			}
			if img.Bounds().Dy() < temp.Dy() {
				cropRect.Max.Y = img.Bounds().Max.Y
			}
		} else {
			if xr > yr {
				maxHeight = temp.Dy()
				maxWidth = int(float64(img.Bounds().Dx()) / yr)
				outx := (maxWidth - temp.Dx()) / 2
				cropRect.Min.X = outx
				cropRect.Max.X = cropRect.Min.X + temp.Dx()
			} else {
				maxWidth = temp.Dx()
				maxHeight = int(float64(img.Bounds().Dy()) / xr)
				outy := (maxHeight - temp.Dy()) / 2
				cropRect.Min.Y = outy
				cropRect.Max.Y = cropRect.Min.Y + temp.Dy()
			}
			imgResized = resize.Thumbnail(uint(maxWidth), uint(maxHeight), img, resize.Lanczos3)
		}
		imgCroped, _ = mergi.Crop(imgResized, cropRect.Min, cropRect.Max)
		rect := temp
		if temp.Dy() > imgCroped.Bounds().Dy() {
			rect.Min.Y = temp.Min.Y + (temp.Dy()-imgCroped.Bounds().Dy())/2
			rect.Max.Y = rect.Min.Y + img.Bounds().Dy()
		}
		if temp.Dx() > imgCroped.Bounds().Dx() {
			rect.Min.X = temp.Min.X + (temp.Dx()-imgCroped.Bounds().Dx())/2
			rect.Max.X = rect.Min.X + img.Bounds().Dx()
		}
		logs.Debug("第{}张图，居中后模板矩形位置:{}", i, rect)
		draw.Draw(canvas, rect, imgCroped, image.ZP, draw.Src)
	}
	return canvas, nil
}

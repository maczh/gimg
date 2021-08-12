package service

import (
	"github.com/maczh/gimg/img"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"github.com/noelyahan/impexp"
	"image"
	"image/color"
	"strconv"
	"strings"
)

var IMAGE_PATH, BASE_URL string

func StitchImage(width, height int, backgroundUrl string, template [][]int, imageUrls []string, borderSize int, borderColor string) mgresult.Result {
	canvas, err := img.StitchImage(width, height, backgroundUrl, generateRects(template), imageUrls, borderSize, hex2color(borderColor))
	if err != nil {
		logs.Error("生成图片错误:{}", err.Error())
		return *mgresult.Error(-1, "生成图片错误:"+err.Error())
	}
	if IMAGE_PATH == "" {
		IMAGE_PATH = mgconfig.GetConfigString("path.img")
	}
	if BASE_URL == "" {
		BASE_URL = mgconfig.GetConfigString("path.url")
	}
	outFileName := utils.GetUUIDString() + ".png"
	err = impexp.NewFileExporter(canvas, IMAGE_PATH+outFileName).Export()
	if err != nil {
		return *mgresult.Error(-1, "图片保存出错:"+err.Error())
	}
	return *mgresult.Success(map[string]string{"url": BASE_URL + outFileName})
}

func hex2color(colorHex string) color.Color {
	if colorHex == "" {
		return color.RGBA{255, 255, 255, 255}
	}
	colorHex = strings.ToLower(colorHex)
	if colorHex[:2] == "0x" {
		colorHex = colorHex[2:]
	}
	r, _ := strconv.ParseInt(colorHex[:2], 16, 16)
	g, _ := strconv.ParseInt(colorHex[2:4], 16, 16)
	b, _ := strconv.ParseInt(colorHex[4:], 16, 16)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func generateRects(template [][]int) []image.Rectangle {
	rects := make([]image.Rectangle, len(template))
	for i, temp := range template {
		rect := image.Rect(temp[0], temp[1], temp[2], temp[3])
		rects[i] = rect
	}
	return rects
}

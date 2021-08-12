package controller

import (
	"github.com/maczh/gimg/service"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/utils"
	"strconv"
)

// StitchImage	godoc
// @Summary		自由拼图
// @Description 通过传入拼图模板和拼图的图片进行自由拼图
// @Tags	拼图
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	width formData int true "结果图片宽度"
// @Param	height formData int true "结果图片高度"
// @Param	backgroundUrl formData string false "背景图片url，若为空则白色背景"
// @Param	template formData string true "拼图模板，二维整数数组，JSON数组格式，每一行数据为一个矩形的左上角和右下角坐标"
// @Param	imageUrls formData string true "所有拼图图片的url地址，字符串数组，JSON数组格式"
// @Param	borderSize formData int false "图片边框宽度，若为0则无边框"
// @Param	borderColor formData string false "边框颜色值，16进制字符串，格式 0x000000或 a3e67f"
// @Success 200 {string} string	"ok"
// @Router	/stitch [post]
func StitchImage(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "width") {
		return *mgresult.Error(-1, "宽度不可为空")
	}
	width, _ := strconv.Atoi(params["width"])
	if !utils.Exists(params, "height") {
		return *mgresult.Error(-1, "高度不可为空")
	}
	height, _ := strconv.Atoi(params["height"])
	if !utils.Exists(params, "template") {
		return *mgresult.Error(-1, "拼图模板数据不可为空")
	}
	template := make([][]int, 0)
	utils.FromJSON(params["template"], &template)
	if !utils.Exists(params, "imageUrls") {
		return *mgresult.Error(-1, "图片地址列表不可为空")
	}
	imageUrls := make([]string, 0)
	utils.FromJSON(params["imageUrls"], &imageUrls)
	borderSize := 0
	if utils.Exists(params, "borderSize") {
		borderSize, _ = strconv.Atoi(params["borderSize"])
	}
	return service.StitchImage(width, height, params["backgroundUrl"], template, imageUrls, borderSize, params["borderColor"])
}

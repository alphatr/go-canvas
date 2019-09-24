package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"time"

	"github.com/alphatr/go-canvas/canvas"
)

func main() {
	image, err := testImage()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("./output.png", image, 0777)
}

func testImage() ([]byte, error) {
	// 缓存获取
	key := fmt.Sprintf("test:%d", time.Now().UnixNano())
	if cache := canvas.GetCache(key); cache != nil {
		return cache, nil
	}

	// 画布初始化
	paint, err := canvas.NewCanvas("assets/background.png")
	if err != nil {
		return nil, err
	}

	// 加载远程图片
	var remote *image.Image
	remoteHandle := func() error {
		res, err := canvas.FetchRemoteImage("https://p2.ssl.qhimg.com/t011cec0ab27b873ef1.jpg", 660, 317, 5*24*time.Hour)
		if err != nil {
			return err
		}

		remote = &res
		return nil
	}

	// 加载本地图片
	var local *image.Image
	localHandle := func() error {
		res, err := canvas.LoadLocalImage("assets/local.png")
		if err != nil {
			return err
		}

		local = &res
		return nil
	}

	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := canvas.NewParallel(timeout, remoteHandle, localHandle); err != nil {
		return nil, err
	}

	paint.DrawImage(*remote, 30, 250)
	paint.DrawImage(*local, 535, 1020)

	// 绘制标题
	title := &canvas.TextConfig{
		Text:     "北国的槐树，也是一种能使人联想起秋来的点缀",
		FontName: "assets/source-han-sans-sc/medium.ttf",
		FontSize: 48,
		Color:    color.Black,
		OffsetX:  375,
		OffsetY:  780,
		AlignX:   0.5,
		MaxWidth: 680,
	}
	if err := paint.DrawText(title); err != nil {
		return nil, err
	}

	// 绘制副标题
	subTitle := &canvas.TextConfig{
		Text:     "象花而又不是花的那一种落蕊，早晨起来，会铺得满地",
		FontName: "assets/source-han-sans-sc/regular.ttf",
		FontSize: 28,
		Color:    color.RGBA{R: 0, G: 0, B: 0, A: 100},
		OffsetX:  375,
		OffsetY:  841,
		AlignX:   0.5,
	}
	if err := paint.DrawText(subTitle); err != nil {
		return nil, err
	}

	// 输出 PNG
	output, err := paint.Output()
	if err != nil {
		return nil, err
	}

	// 缓存
	if err := canvas.SaveCache(key, output, 2*24*time.Hour); err != nil {
		log.Print(err.Error())
	}

	return output, nil
}

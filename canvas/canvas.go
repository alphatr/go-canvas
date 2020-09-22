package canvas

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"time"

	"github.com/fogleman/gg"
)

// Canvas 画布
type Canvas struct {
	context *gg.Context
	Width   uint
	Height  uint
}

// TextConfig 绘制文字配置
type TextConfig struct {
	Text     string
	FontName string
	FontSize float64
	Color    color.Color
	OffsetX  float64
	OffsetY  float64
	AlignX   float64
	AlignY   float64
	MaxWidth float64
}

// NewCanvasWithLocal 从本地图片创建新的画布
func NewCanvasWithLocal(background string) (*Canvas, error) {
	img, err := LoadLocalImage(background)
	if err != nil {
		return nil, err
	}

	return NewCanvas(img), nil
}

// NewCanvasWithRemote 从远程图片创建新的画布
func NewCanvasWithRemote(source string, width uint, height uint, expired time.Duration) (*Canvas, error) {
	img, err := FetchRemoteImage(source, width, height, expired)
	if err != nil {
		return nil, err
	}

	return NewCanvas(img), nil
}

// NewCanvas 创建新的画布
func NewCanvas(background image.Image) *Canvas {
	ins := &Canvas{
		context: gg.NewContext(background.Bounds().Dx(), background.Bounds().Dy()),
		Width:   uint(background.Bounds().Dx()),
		Height:  uint(background.Bounds().Dy()),
	}

	ins.context.DrawImage(background, 0, 0)
	return ins
}

// DrawImage 绘制图片
func (ins *Canvas) DrawImage(img image.Image, offsetX int, offsetY int) {
	ins.context.DrawImage(img, offsetX, offsetY)
}

// DrawLine 绘制横线
func (ins *Canvas) DrawLine(color color.Color, offsetX, offsetY, length float64) {
	ins.context.DrawLine(offsetX, offsetY, offsetX+length, offsetY)
	ins.context.ClosePath()
	ins.context.SetLineWidth(1)
	ins.context.SetColor(color)
	ins.context.StrokePreserve()
	ins.context.Stroke()
}

// MeasureString 测量文字
func (ins *Canvas) MeasureString(opt *TextConfig) float64 {
	font, err := LoadLocalFont(opt.FontName, opt.FontSize)
	if err != nil {
		return 0
	}

	ins.context.SetFontFace(font)
	width, _ := ins.context.MeasureString(opt.Text)
	return width
}

// DrawText 绘制图片
func (ins *Canvas) DrawText(opt *TextConfig) error {
	font, err := LoadLocalFont(opt.FontName, opt.FontSize)
	if err != nil {
		return err
	}

	ins.context.SetColor(opt.Color)

	ins.context.SetFontFace(font)
	width := ins.MeasureString(opt)
	if opt.MaxWidth > 0 && width > opt.MaxWidth {
		font, err = LoadLocalFont(opt.FontName, opt.FontSize*opt.MaxWidth/width)
		if err != nil {
			return err
		}
	}

	ins.context.SetFontFace(font)
	ins.context.DrawStringAnchored(opt.Text, opt.OffsetX, opt.OffsetY, opt.AlignX, opt.AlignY)
	return nil
}

// Output 输出 PNG 图片
func (ins *Canvas) Output() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, ins.context.Image()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

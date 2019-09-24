package canvas

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

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

// NewCanvas 返回新的画布
func NewCanvas(background string) (*Canvas, error) {
	img, err := LoadLocalImage(background)
	if err != nil {
		return nil, err
	}

	ins := &Canvas{
		context: gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy()),
		Width:   uint(img.Bounds().Dx()),
		Height:  uint(img.Bounds().Dy()),
	}

	ins.context.DrawImage(img, 0, 0)
	return ins, nil
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

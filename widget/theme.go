package widget

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Theme struct {
	DefaultFace      *text.Face
	DefaultTextColor color.Color
	ButtonTheme      *ButtonParams
	PanelTheme       *PanelParams
	LabelTheme       *LabelParams
	TextTheme        *TextParams
}
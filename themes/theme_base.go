package themes

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("%w", err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

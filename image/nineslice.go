package image

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

// A NineSlice is an image that can be drawn with any width and height. It is basically a 3x3 grid of image tiles:
// The corner tiles be drawn as-is, while the center columns and rows of tiles will be stretched to fit the desired
// width and height.
type NineSlice struct {
	image       *ebiten.Image
	widths      []int
	heights     []int
	transparent bool

	init      sync.Once
	subImages []*ebiten.Image
}

// A DrawImageOptionsFunc is responsible for setting DrawImageOptions when drawing an image.
// This is usually used to translate the image.
type DrawImageOptionsFunc func(opts *ebiten.DrawImageOptions)

var colorImages map[color.Color]*ebiten.Image = map[color.Color]*ebiten.Image{}

var colorNineSlices map[color.Color]*NineSlice = map[color.Color]*NineSlice{}

// NewNineSliceSimple constructs a new NineSlice from image. borderWidthHeight specifies the width of the
// left and right column and the height of the top and bottom row. centerWidthHeight specifies the width
// of the center column and row.
func NewNineSliceSimple(image *ebiten.Image, borderWidthHeight int, centerWidthHeight int) *NineSlice {
	return &NineSlice{
		image:   image,
		widths:  []int{borderWidthHeight, centerWidthHeight, borderWidthHeight},
		heights: []int{borderWidthHeight, centerWidthHeight, borderWidthHeight},
	}
}

// NewNineSliceColor constructs a new NineSlice that when drawn fills with color c.
func NewNineSliceColor(c color.Color) *NineSlice {
	if n, ok := colorNineSlices[c]; ok {
		return n
	}

	var n *NineSlice
	if c == color.Transparent {
		n = &NineSlice{
			transparent: true,
		}
	} else {
		n = &NineSlice{
			image:   NewImageColor(c),
			widths:  []int{0, 1, 0},
			heights: []int{0, 1, 0},
		}
	}
	colorNineSlices[c] = n
	return n
}

// NewImageColor constructs a new Image that when drawn fills with color c.
func NewImageColor(c color.Color) *ebiten.Image {
	if i, ok := colorImages[c]; ok {
		return i
	}

	i, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
	_ = i.Fill(c)
	colorImages[c] = i
	return i
}

// Draw draws n onto screen, with the size specified by width and height. If optsFunc is not nil, it is used to set
// DrawImageOptions for each tile drawn.
func (n *NineSlice) Draw(screen *ebiten.Image, width int, height int, optsFunc DrawImageOptionsFunc) {
	if n.transparent {
		return
	}

	n.init.Do(n.createSubImages)

	sy := 0
	ty := 0
	for r, sh := range n.heights {
		if sh <= 0 {
			continue
		}

		sx := 0
		tx := 0

		var th int
		if r == 1 {
			th = height - n.heights[0] - n.heights[2]
		} else {
			th = sh
		}

		for c, sw := range n.widths {
			if sw <= 0 {
				continue
			}

			var tw int
			if c == 1 {
				tw = width - n.widths[0] - n.widths[2]
			} else {
				tw = sw
			}

			opts := ebiten.DrawImageOptions{
				Filter: ebiten.FilterNearest,
			}

			if tw != sw || th != sh {
				opts.GeoM.Scale(float64(tw)/float64(sw), float64(th)/float64(sh))
			}

			opts.GeoM.Translate(float64(tx), float64(ty))

			if optsFunc != nil {
				optsFunc(&opts)
			}

			_ = screen.DrawImage(n.subImages[r*3+c], &opts)

			sx += sw
			tx += tw
		}

		sy += sh
		ty += th
	}
}

func (n *NineSlice) createSubImages() {
	defer func() {
		n.image = nil
	}()

	n.subImages = make([]*ebiten.Image, 9)

	// short-circuit if only the center tile is used
	if n.widths[0] <= 0 && n.widths[2] <= 0 && n.heights[0] <= 0 && n.heights[2] <= 0 {
		w, h := n.image.Size()
		if n.widths[1] == w && n.heights[1] == h {
			n.subImages[1*3+1] = n.image
			return
		}
	}

	sy := 0
	for r, sh := range n.heights {
		sx := 0
		for c, sw := range n.widths {
			if sh > 0 && sw > 0 {
				n.subImages[r*3+c] = n.image.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image)
			}
			sx += sw
		}
		sy += sh
	}
}

package widget

import (
	"image"
	"math"
	"testing"

	"github.com/matryer/is"
)

func TestGridLayout_WeightsResize(t *testing.T) {
	is := is.New(t)
	g := NewGridLayout(GridLayoutOpts.Columns(3),
		GridLayoutOpts.Weights([]float64{20, 50, 30}, []float64{1, 3}),
		GridLayoutOpts.Padding(NewInsetsSimple(10)),
		GridLayoutOpts.Spacing(5, 10))
	widgets := make([]PreferredSizeLocateableWidget, 6)
	for i := range widgets {
		widgets[i] = newSimpleWidget(10, 10, nil)
	}
	// Both rows share column widths, even after resizing. Padding, spacing,
	// and a non-zero container origin must all be respected.
	for _, size := range []int{100, 200} {
		rect := image.Rect(30, 40, 30+size+30, 40+size+30)
		g.Layout(widgets, rect)
		widths := []int{size * 20 / 100, size * 50 / 100, size * 30 / 100}
		heights := []int{size / 4, size * 3 / 4}
		y := 50
		for r, height := range heights {
			x := 40
			for c, width := range widths {
				is.Equal(widgets[r*3+c].GetWidget().Rect, image.Rect(x, y, x+width, y+height))
				x += width + 5
			}
			y += height + 10
		}
	}
	w, h := g.PreferredSize(widgets)
	is.Equal(w, 60)
	is.Equal(h, 50)
}

func TestGridLayout_WeightedCellSizes(t *testing.T) {
	tests := []struct {
		name           string
		weights        []float64
		stretch        []bool
		defaultStretch bool
		available      int
		want           []int
	}{
		{"relative fractions", []float64{.2, .5, .3}, nil, false, 110, []int{20, 50, 30}},
		{"rounding", []float64{2, 5, 3}, nil, false, 111, []int{21, 50, 30}},
		{"fixed cell", []float64{0, 1, 3}, nil, true, 110, []int{10, 23, 67}},
		{"partial weights", []float64{3}, nil, true, 110, []int{60, 20, 20}},
		{"zero overrides stretch", []float64{0, 0, 0}, nil, true, 110, []int{10, 10, 10}},
		{"legacy stretch", nil, []bool{true, false, true}, false, 111, []int{46, 10, 45}},
		{"default stretch", nil, nil, true, 110, []int{34, 33, 33}},
		{"undersized container", []float64{0, 1, 3}, nil, false, 15, []int{10, 0, 0}},
		{"large weights", []float64{math.MaxFloat64, math.MaxFloat64, 0}, nil, false, 110, []int{45, 45, 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			g := NewGridLayout(GridLayoutOpts.Columns(3),
				GridLayoutOpts.Weights(tt.weights, nil),
				GridLayoutOpts.Stretch(tt.stretch, nil),
				GridLayoutOpts.DefaultStretch(tt.defaultStretch, false))
			is.Equal(g.weightedCellSizes([]int{10, 10, 10}, tt.available, 5, g.columnWeights, g.columnStretched), tt.want)
		})
	}
}

func TestGridLayout_WeightsMaxSizeAndPosition(t *testing.T) {
	is := is.New(t)
	g := NewGridLayout(GridLayoutOpts.Columns(1), GridLayoutOpts.Weights([]float64{1}, []float64{1}))
	w := newSimpleWidget(10, 10, GridLayoutData{
		MaxWidth: 20, MaxHeight: 30,
		HorizontalPosition: GridLayoutPositionCenter,
		VerticalPosition:   GridLayoutPositionEnd,
	})
	g.Layout([]PreferredSizeLocateableWidget{w}, image.Rect(10, 20, 110, 120))
	is.Equal(w.GetWidget().Rect, image.Rect(50, 90, 70, 120))
}

func TestGridLayout_InvalidWeights(t *testing.T) {
	for _, weight := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		for _, row := range []bool{false, true} {
			func() {
				defer func() {
					if recover() == nil {
						t.Errorf("expected panic for weight %v (row=%v)", weight, row)
					}
				}()
				c, r := []float64{weight}, []float64(nil)
				if row {
					c, r = r, c
				}
				NewGridLayout(GridLayoutOpts.Columns(1), GridLayoutOpts.Weights(c, r))
			}()
		}
	}
}

func TestGridLayout_EmptyWeightedLayout(t *testing.T) {
	g := NewGridLayout(GridLayoutOpts.Columns(3), GridLayoutOpts.Weights([]float64{20, 50, 30}, nil))
	g.Layout(nil, image.Rect(0, 0, 100, 100))
}

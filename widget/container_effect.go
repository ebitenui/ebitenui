package widget

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ContainerShader draws a container through a Kage shader: the shader's
// first image is the container and its children as they would be drawn,
// and Uniforms are its uniform variables, which may change between frames.
type ContainerShader struct {
	Shader   *ebiten.Shader
	Uniforms map[string]any
}

// Opacity draws the container and its children at an opacity from 0,
// invisible, to 1, opaque, the default. They are drawn as one image, so
// overlapping children don't show through each other.
func (o ContainerOptions) Opacity(opacity float32) ContainerOpt {
	return func(c *Container) {
		c.opacity = opacity
	}
}

// Shader draws the container and its children through a shader.
func (o ContainerOptions) Shader(shader *ContainerShader) ContainerOpt {
	return func(c *Container) {
		c.shader = shader
	}
}

// Opacity returns the opacity the container is drawn at.
func (c *Container) Opacity() float32 {
	return c.opacity
}

// SetOpacity sets the opacity the container and its children are drawn at,
// from 0, invisible, to 1, opaque.
func (c *Container) SetOpacity(opacity float32) {
	c.opacity = opacity
	c.releaseEffectBuffer()
}

// Shader returns the shader the container is drawn through, or nil.
func (c *Container) Shader() *ContainerShader {
	return c.shader
}

// SetShader sets the shader the container and its children are drawn
// through; nil draws them directly.
func (c *Container) SetShader(shader *ContainerShader) {
	c.shader = shader
	c.releaseEffectBuffer()
}

// hasEffect reports whether the container is drawn through an off-screen
// image.
func (c *Container) hasEffect() bool {
	return c.opacity < 1 || c.shader != nil
}

// effectMargin is how far, in pixels, the off-screen image reaches past the
// container on each side, so that a container sliding a little, such as a
// window easing in, keeps its image.
const effectMargin = 64

// releaseEffectBuffer frees the off-screen image once the container no
// longer needs it: drawn directly again, or invisible, as a window faded out
// as it closes.
func (c *Container) releaseEffectBuffer() {
	if c.effectBuffer != nil && (!c.hasEffect() || c.opacity <= 0) {
		c.effectBuffer.Deallocate()
		c.effectBuffer = nil
	}
}

// renderWithEffect renders the container and its children into an
// off-screen image, then draws that image onto screen at the container's
// opacity, through its shader if it has one. The image covers the container
// and effectMargin pixels around it, within screen, in screen's coordinates,
// and is kept while the container stays inside it. It is a texture of its
// own, outside Ebitengine's atlases, so freeing it gives its memory back;
// on an atlas it would grow pages that stay grown.
func (c *Container) renderWithEffect(screen *ebiten.Image) {
	r := c.widget.Rect.Intersect(screen.Bounds())
	if r.Empty() || c.opacity <= 0 {
		c.doLayout()
		return
	}
	if c.effectBuffer == nil || !r.In(c.effectBuffer.Bounds()) {
		if c.effectBuffer != nil {
			c.effectBuffer.Deallocate()
		}
		bounds := r.Inset(-effectMargin).Intersect(screen.Bounds())
		c.effectBuffer = ebiten.NewImageWithOptions(bounds, &ebiten.NewImageOptions{Unmanaged: true})
	}
	target, ok := c.effectBuffer.SubImage(r).(*ebiten.Image)
	if !ok {
		return
	}
	target.Clear()
	c.renderContent(target)

	if c.shader == nil {
		op := &c.effectImageOp
		op.GeoM.Reset()
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
		op.ColorScale.Reset()
		op.ColorScale.ScaleAlpha(c.opacity)
		screen.DrawImage(target, op)
		return
	}
	op := &c.effectShaderOp
	op.GeoM.Reset()
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	op.ColorScale.Reset()
	op.ColorScale.ScaleAlpha(c.opacity)
	op.Uniforms = c.shader.Uniforms
	op.Images[0] = target
	screen.DrawRectShader(r.Dx(), r.Dy(), c.shader.Shader, op)
	op.Images[0] = nil
}

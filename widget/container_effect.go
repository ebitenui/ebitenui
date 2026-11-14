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

// releaseEffectBuffer frees the off-screen image once the container no
// longer needs it.
func (c *Container) releaseEffectBuffer() {
	if c.effectBuffer != nil && !c.hasEffect() {
		c.effectBuffer.Deallocate()
		c.effectBuffer = nil
	}
}

// renderWithEffect renders the container and its children into an
// off-screen image, then draws that image onto screen at the container's
// opacity, through its shader if it has one. The image covers screen and is
// kept until screen's bounds change, so a container that moves, such as a
// window sliding in, does not allocate a new one.
func (c *Container) renderWithEffect(screen *ebiten.Image) {
	r := c.widget.Rect.Intersect(screen.Bounds())
	if r.Empty() || c.opacity <= 0 {
		c.doLayout()
		return
	}
	if c.effectBuffer == nil || c.effectBuffer.Bounds() != screen.Bounds() {
		if c.effectBuffer != nil {
			c.effectBuffer.Deallocate()
		}
		c.effectBuffer = ebiten.NewImageWithOptions(screen.Bounds(), nil)
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

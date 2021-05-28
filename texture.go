package xui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"image"
)

type Texture struct {
	Rgba    *image.RGBA
	Width   int
	Height  int
	IsAlpha bool
	texture uint32
}

func NewTexture(rgba *image.RGBA) (*Texture) {
	p := rgba.Bounds().Size()
	t := &Texture{Rgba: rgba, Width: p.X, Height: p.Y,IsAlpha:true}

	t.init()
	return t
}

func (t *Texture) init() {
	if t.texture != 0 {
		return
	}
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(t.Rgba.Rect.Size().X),
		int32(t.Rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(t.Rgba.Pix))

	t.texture = texture


}
func (t *Texture) Recycle() {
	gl.DeleteTextures(1, &t.texture)
	t.texture = 0
}



func (t *Texture) Draw(c XCanvas, x, y float64,winWidth, winHeight int) {
	tx, ty := c.GetTranslate()


	x = x + tx
	y = y + ty


	x, y = AppCoordinate2OpenGL(winWidth, winHeight, x, y)
	w, h := AppWidthHeight2OpenGL(winWidth, winHeight, (float64(t.Width)), (float64(t.Height)))

	if t.IsAlpha {
		gl.Enable(gl.BLEND);
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
	}

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.BindTexture(gl.TEXTURE_2D, t.texture)
	gl.LineWidth(0)
	gl.PointSize(0)
	gl.Begin(gl.QUADS)

	//gl.Normal3f(float32(x), float32(y-h), 0) //
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(float32(x+w), float32(y-h), 0) //

	gl.TexCoord2f(1, 0)
	gl.Vertex3f(float32(x+w), float32(y), 0) //
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(float32(x), float32(y), 0) //
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(float32(x), float32(y-h), 0) //

	gl.End()

	gl.PopMatrix()
	t.Recycle()

}
func (t *Texture) DrawInRetangle(c XCanvas, x, y ,rx,ry, rw,rh float64,winWidth, winHeight int) {
	rcx,rcy:=rx-x,ry-y
	rcw,rch:=rw,rh


	x, y = AppCoordinate2OpenGL(winWidth, winHeight, rx, ry)
	w, h := AppWidthHeight2OpenGL(winWidth, winHeight, (float64(rw)), (float64(rh)))

	///rx, ry = AppCoordinate2OpenGL(winWidth, winHeight, rx, ry)
	rw, rh = float64(t.Width),float64(t.Height)


	if t.IsAlpha {
		gl.Enable(gl.BLEND);
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
	}

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.BindTexture(gl.TEXTURE_2D, t.texture)
	gl.LineWidth(0)
	gl.PointSize(0)
	gl.Begin(gl.QUADS)


	gl.TexCoord2f(AppCoordinate2Texture(rw,rh,rcx+rcw,rcy+rch))//1,1
	//gl.TexCoord2f(1,1)

	gl.Vertex3f(float32(x+w), float32(y-h), 0) //

	gl.TexCoord2f(AppCoordinate2Texture(rw,rh,rcx+rcw,rcy))//1,0
	//gl.TexCoord2f(1,0)

	gl.Vertex3f(float32(x+w), float32(y), 0) //

	gl.TexCoord2f(AppCoordinate2Texture(rw,rh,rcx,rcy))//0,0
	//gl.TexCoord2f(0,0)

	gl.Vertex3f(float32(x), float32(y), 0) //

	gl.TexCoord2f(AppCoordinate2Texture(rw,rh,rcx,rcy+rch))//0,1
	//gl.TexCoord2f(0,1)

	gl.Vertex3f(float32(x), float32(y-h), 0) //

	gl.End()

	gl.PopMatrix()
	t.Recycle()

}
func AppCoordinate2OpenGL(w, h int, x, y float64) (float64, float64) {

	return x*6/float64(w) - 3, -y*6/float64(h) + 3
}
func AppWidthHeight2OpenGL(w, h int, x, y float64) (float64, float64) {
	return x * 6 / float64(w), y * 6 / float64(h)
}
func AppWidthHeight2Texture(w, h int, x, y float64)(float64, float64){
	return x/float64(w),y/float64(h)
}
func AppCoordinate2Texture(w, h , x, y float64) (float32, float32) {

	return float32(x/float64(w)), float32((y)/float64(h))
}
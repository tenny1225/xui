package xui

import (
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type XPaint struct {
	Color       color.Color
	Size        float64
	StrokeWidth float64
}

type XCanvas interface {
	DrawCircle(dx, dy, radius float64, p XPaint)
	DrawLine(x1, y1, x2, y2 float64, p XPaint)
	DrawText(str string,x, y ,ax,ay ,mw,mh float64, p XPaint)
	DrawRect(x1, y1, x2, y2, fillet float64, p XPaint)
	DrawImage(x, y float64, img image.Image)

	SetTranslate(x, y float64)
	SetScale(sx, sy float64)
	SetRotate(dx, dy, angle float64)

	Save()
	Restore()

	SetAlpha(a float64)

	GetTranslate() (float64, float64)
	GetScale() (float64, float64)
	GetRotate() (float64, float64, float64)
	GetAlpha() float64
	GetWindow() XWindow
}

type xcanvas struct {
	translateX float64
	translateY float64
	scaleX     float64
	scaleY     float64
	rotateX    float64
	rotateY    float64
	rotate     float64
	alpha      float64

	backups []xcanvas
	window  XWindow
}

func (z *xcanvas)MeasureString(str string)(float64,float64)  {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	c := gg.NewContextForRGBA(img)
	return c.MeasureString(str)
}
func (z *xcanvas) DrawText(str string,x, y ,ax,ay,mw,mh float64, p XPaint){
	w,h:=z.MeasureString(str)
	if mw>0&&w>mw{
		w =mw

	}
	if mh>0&&h>mh{
		h=mh
	}
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	c := gg.NewContextForRGBA(img)
	c.SetColor(p.Color)
	c.DrawStringAnchored(str, w*ax, h*ay, ax, ay)

	c.Fill()
	z.DrawImage(x-w*ax, y-h*ax, img)
}
func (z *xcanvas) DrawCircle(dx, dy, radius float64, p XPaint) {
	rgba := image.NewRGBA(image.Rect(0, 0, int(2*radius), int(2*radius)))
	dc := gg.NewContextForRGBA(rgba)
	dc.DrawCircle(radius, radius, radius)
	dc.SetColor(p.Color)
	dc.Fill()
	t := NewTexture(rgba)
	w,h:=z.window.GetSize()
	t.Draw(z, dx-radius, dy-radius,w,h);
}

func (z *xcanvas) DrawLine(x1, y1, x2, y2 float64, p XPaint) {
	w := int(math.Abs(x2 - x1))
	if w == 0 {
		w = int(p.StrokeWidth)
	}
	h := int(math.Abs(y2 - y1))
	if h == 0 {
		h = int(p.StrokeWidth)
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	dc := gg.NewContextForRGBA(rgba)
	dc.DrawLine(0, 0, math.Abs(x2-x1), math.Abs(y2-y1))
	dc.SetLineWidth(p.StrokeWidth)
	dc.SetColor(p.Color)
	dc.Stroke()
	t := NewTexture(rgba)
	ww,wh:=z.window.GetSize()
	t.Draw(z, math.Min(x1, x2), math.Min(y1, y2),ww,wh);
}

func (z *xcanvas) DrawRect(x1, y1, x2, y2, fillet float64, p XPaint) {

	rgba := image.NewRGBA(image.Rect(0, 0, int(math.Abs(x2-x1)), int(math.Abs(y2-y1))))
	dc := gg.NewContextForRGBA(rgba)
	dc.DrawRoundedRectangle(0, 0, math.Abs(x2-x1), math.Abs(y2-y1), fillet)
	dc.SetColor(p.Color)
	dc.Fill()
	t := NewTexture(rgba)
	w,h:=z.window.GetSize()

	t.Draw(z, math.Min(x1, x2), math.Min(y1, y2),w,h)
}

func (z *xcanvas) DrawImage(x, y float64, img image.Image) {
	defer func() {
		if e:=recover();e!=nil{
			//fmt.Println(e)
		}
	}()
	f:= func() {
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
		t := NewTexture(rgba)
		t.IsAlpha = true
		w,h:=z.window.GetSize()
		t.Draw(z, x, y,w,h);
	}

	f()

}

func (z *xcanvas) SetTranslate(x, y float64) {
	z.translateX = x
	z.translateY = y
}

func (z *xcanvas) SetScale(sx, sy float64) {
	z.scaleX = sx
	z.scaleY = sy
}

func (z *xcanvas) SetRotate(dx, dy, angle float64) {
	z.rotateX = dx
	z.rotateY = dy
	z.rotate = angle
}

func (z *xcanvas) Save() {
	z.backups = append(z.backups, xcanvas{
		translateX: z.translateX,
		translateY: z.translateY,
		scaleX:     z.scaleX,
		scaleY:     z.scaleY,
		rotateX:    z.rotateX,
		rotateY:    z.rotateY,
		rotate:     z.rotate,
		alpha:      z.alpha,
	})

	z.translateX = 0
	z.translateY = 0
	z.scaleX = 0
	z.scaleY = 0
	z.rotateX = 0
	z.rotateY = 0
	z.rotate = 0
	z.alpha = 0
}

func (z *xcanvas) Restore() {
	if len(z.backups) == 0 {
		return
	}
	last := z.backups[len(z.backups)-1]
	z.backups = z.backups[:len(z.backups)-1]
	z.translateX = last.translateX
	z.translateY = last.translateY
	z.scaleX = last.scaleX
	z.scaleY = last.scaleY
	z.rotateX = last.rotateX
	z.rotateY = last.rotateY
	z.rotate = last.rotate
	z.alpha = last.alpha
}

func (z *xcanvas) SetAlpha(a float64) {
	z.alpha = a
}

func (z *xcanvas) GetTranslate() (float64, float64) {
	x:=z.translateX
	y:=z.translateY

	for _,l:=range z.backups{
		x+=l.translateX
		y+=l.translateY
	}
	return x, y
}

func (z *xcanvas) GetScale() (float64, float64) {
	x:=z.scaleX
	y:=z.scaleY
	for _,l:=range z.backups{
		x=x*l.scaleX
		y=y*l.scaleY
	}
	return x, y
}

func (z *xcanvas) GetRotate() (float64, float64, float64) {
	x:=z.rotateX
	y:=z.rotateY
	r:=z.rotate
	for _,l:=range z.backups{
		x+=l.rotateX
		y+=l.rotateY
		r+=z.rotate
	}
	return x, y, r
}

func (z *xcanvas) GetAlpha() float64 {
	a:=z.alpha
	for _,l:=range z.backups{
		a=a*l.alpha
	}
	return a
}
func (z *xcanvas) GetWindow() XWindow {
	return z.window
}

func NewCanvas(window XWindow) XCanvas {
	return &xcanvas{
		backups: make([]xcanvas, 0),
		window:  window,
	}

}

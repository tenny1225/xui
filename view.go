package xui

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type CursorType int
type Direction int

const (
	Down  CursorType = 1
	Up    CursorType = 2
	Out   CursorType = 3
	Hover CursorType = 4

	None       Direction = 0
	Vertical   Direction = 1
	Horizontal Direction = 2
)
const (
	AUTO_SCOLL_PARENT = -2
	FULL_PARENT       = -1
	AUTO              = 0
)
const ScollBarMinWidth float64 = 20

type Viewer interface {
	render(canvas *gg.Context)

	measure(w, h float64) (float64, float64)
	setSize(w, h float64)

	getSettingSize() (float64, float64)
	getPosition() (float64, float64)
	setCurrentPosition(x,y float64)
	GetCurrentPosition()(float64, float64)
	getEvent() func(x, y float64, t CursorType) bool
	setPosition(x, y float64)
	getTitle() string
	isShouldMeasure() bool
	RequestLayout()
	setParent(view *View)
	init(w XWindow)
	getRGBA() *image.RGBA
	pushString(s string)bool
	backspace()bool
	addCursorIndex(i int)bool

	setRGBA(rgba *image.RGBA)
	getBackground() *image.RGBA
	setBackground(rgba *image.RGBA)

	Event(x, y float64, t CursorType) bool
	Scroll(x, y float64, distance float64) bool
	GetMeasureSize() (float64, float64)
	GetFont() font.Face
	SetScroll(s float64)
	GetSize() (float64, float64)
	GetDirection() Direction
	SetChildCount(c int)
	SetFocus(x,y float64,b bool) bool
}

type View struct {
	Parent          *View
	BackgroundColor color.Color
	PrimaryColor   color.Color
	AccentColor color.Color
	Children        []Viewer
	Width           float64
	Height          float64
	PaddingLeft     float64
	BorderWidth      float64
	BorderRoundWidth float64
	PaddingTop      float64
	PaddingRight    float64
	PaddingBottom   float64
	TextColor       color.Color
	Direction       Direction
	ScrollLength    float64
	MaxScrollLength float64
	LineCount int
	FontSize int
	FontPath string
	ScaleType    ImageScaleType
	Src  string
	currentLeft float64
	currentTop float64

	Left          float64
	Top           float64
	Title         string
	ShouldMeasure bool
	Drawer        func(canvas *gg.Context)
	Clicker       func(v *View, x, y float64)
	Measurer      func(w, h float64) (float64, float64)
	Eventer       func(x, y float64, t CursorType) bool
	Scroller      func(x, y float64, distance float64) bool
	Focuser       func(x,y float64,b bool)bool
	GetView       func(i int) Viewer
	GetChildCount func() int

	Backgrounder func(canvas *gg.Context)

	Font font.Face

	rgba            *image.RGBA
	background      *image.RGBA
	currentWidth    float64
	currentHeight   float64
	Window          XWindow
	cursorView      Viewer
	cursorIndex     int
	isFocus         bool
	currentPoint    []float64
	downPoint       []float64
	isDownScollBar  bool
	isDynamicRender bool
	childCount      int
}
func (v *View)setCurrentPosition(x,y float64){
	v.currentLeft,v.currentTop=x,y
}
func (v *View)GetCurrentPosition()(float64, float64){
	return v.currentLeft,v.currentTop
}
func (v *View) SetChildCount(c int) {
	v.childCount = c
	v.isDynamicRender = c > 0

}
func (v *View) getRGBA() *image.RGBA {
	return v.rgba
}
func (v *View) setRGBA(rgba *image.RGBA) {
	v.rgba = rgba
}
func (v *View) getBackground() *image.RGBA {
	return v.background
}
func (v *View) setBackground(rgba *image.RGBA) {
	v.background = rgba
}
func (v *View) measure(mw, mh float64) (float64, float64) {

	if v.Children != nil {
		for _, l := range v.Children {
			if l != nil && l.isShouldMeasure() {
				l.setSize(l.measure(l.getSettingSize()))
			}

		}
	}
	if v.Measurer != nil && v.ShouldMeasure {
		v.ShouldMeasure = false
		return v.Measurer(v.getSettingSize())
	} else {
		winWidth, winHeight := v.GetMeasureSize()

		width := 0.0
		height := 0.0
		max := 0.0
		if len(v.Children) > 0 {
			for _, l := range v.Children {
				if l != nil {
					x, y := l.getPosition()
					w, h := l.GetSize()

					width = math.Min(float64(math.Max(float64(width), float64(x+w))), float64(winWidth))
					height = math.Min(float64(math.Max(float64(height), float64(y+h))), float64(winHeight))

					if v.Direction == Vertical {
						max += float64(y + h)
					} else if v.Direction == Horizontal {
						max += float64(x + w)
					}
				}

			}
		}

		if v.Width > 0 {
			width = v.Width
		}
		if v.Height > 0 {
			height = v.Height
		}
		if v.Width == FULL_PARENT {
			width = winWidth
		}
		if v.Height == FULL_PARENT {
			height = winHeight
		}
		if height > 0 && v.Direction == Vertical {
			v.MaxScrollLength = height - max
		}
		if width > 0 && v.Direction == Horizontal {
			v.MaxScrollLength = width - max
		}
		if v.MaxScrollLength > 0 {
			v.MaxScrollLength = 0
		}
		return width, height
	}
	return v.Width, v.Height

}
func (v *View) setParent(view *View) {
	v.Parent = view

}
func (v *View) pushString(s string)bool{
	for _,l:=range v.Children{
		if l.pushString(s){
			return true
		}
	}
	if v.isFocus{
		v.Title = string([]rune(v.Title)[:v.cursorIndex])+s+string([]rune(v.Title)[v.cursorIndex:])
		v.cursorIndex++
		return true
	}
	return false
}
func (v *View)backspace()bool{
	for _,l:=range v.Children{
		if l.backspace(){
			return true
		}
	}
	if v.isFocus{
		if v.Title!=""{
			list:=make([]rune,0)
			for i,l:=range []rune(v.Title){
				if i+1!=v.cursorIndex{
					list = append(list,l)
				}
			}
			v.Title = string(list)
			v.cursorIndex--
		}

		if v.cursorIndex<0{
			v.cursorIndex = 0
		}
		return true
	}
	return false
}
func (v *View)addCursorIndex(i int)bool{
	for _,l:=range v.Children{
		if l.addCursorIndex(i){
			return true
		}
	}
	if v.isFocus{
		if v.Title!=""{
			v.cursorIndex+=i
		}

		if v.cursorIndex<0{
			v.cursorIndex = 0
		}else if v.cursorIndex>len([]rune(v.Title)){
			v.cursorIndex = len([]rune(v.Title))
		}
		return true
	}
	return false
}
func (v *View) SetFocus(x,y float64,b bool) bool{
	for _,l:=range v.Children{
		if l.SetFocus(x,y,b){
			v.cursorView=l
			return true
		}
	}
	if !b{
		v.isFocus=b
	}
	if v.Focuser!=nil{
		return v.Focuser(x,y,b)
	}
	return false
}
func (v *View) RequestLayout() {
	v.Window.RequestLayout()
}
func (v *View) GetFont() font.Face {
	return v.Font
}
func (v *View) GetDirection() Direction {
	return v.Direction
}
func (v *View) isShouldMeasure() bool {
	return v.ShouldMeasure
}
func (v *View) getSettingSize() (float64, float64) {

	w := v.Width
	h := v.Height
	return w, h
}
func (v *View) SetScroll(s float64) {
	v.ScrollLength += s
	if v.ScrollLength < v.MaxScrollLength {
		v.ScrollLength = v.MaxScrollLength
	} else if v.ScrollLength > 0 {
		v.ScrollLength = 0
	}
	v.RequestLayout()
}
func (v *View) GetMeasureSize() (float64, float64) {

	w := v.Width
	h := v.Height

	if (w == AUTO || w == FULL_PARENT) && v.Parent != nil {
		if v.Parent.Direction == Horizontal {
			w = AUTO_SCOLL_PARENT
		} else {
			w, _ = v.Parent.GetMeasureSize()
		}

	}
	if (h == AUTO || h == FULL_PARENT) && v.Parent != nil {
		if v.Parent.Direction == Vertical {
			h = AUTO_SCOLL_PARENT
		} else {
			_, h = v.Parent.GetMeasureSize()
		}

	}

	return w, h
}
func (v *View) GetSize() (float64, float64) {

	return v.currentWidth, v.currentHeight

}
func (v *View) setSize(w, h float64) {

	v.currentWidth = w
	v.currentHeight = h

}
func (v *View) getTitle() string {
	return v.Title
}
func (v *View) getPosition() (float64, float64) {

	return v.Left, v.Top

}
func (v *View) setPosition(x, y float64) {

	v.Left = x
	v.Top = y

}

func (v *View) getEvent() func(x, y float64, t CursorType) bool {
	return v.Eventer
}
func (v *View) init(w XWindow) {
	v.currentWidth, v.currentHeight = v.Width, v.Height
	v.Window = w
	v.ShouldMeasure = true
	v.rgba = nil
	v.background = nil
	if v.Children != nil {
		for _, l := range v.Children {
			if l == nil {
				continue
			}
			l.setParent(v)
			l.init(w)
		}
	}

}
func (v *View) render(canvas *gg.Context) {

	if v.Children == nil {
		if v.Direction != None && v.GetChildCount != nil {
			v.SetChildCount(v.GetChildCount())
		}
		if v.isDynamicRender {
			v.Children = make([]Viewer, v.childCount)
		} else {
			v.Children = make([]Viewer, 0)
		}

	}
	if v.ShouldMeasure {
		v.setSize(v.measure(v.getSettingSize()))
	}
	w, h := v.GetSize()
	if v.background == nil  {
		if v.Backgrounder == nil {
			img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
			dx, dy := float64(w), float64(h)
			ctx := gg.NewContextForRGBA(img)
			if v.GetFont() != nil {
				//ctx.SetFontFace(v.GetFont())
			}
			if v.BackgroundColor != nil{
				ctx.SetColor(v.BackgroundColor)

				ctx.DrawRoundedRectangle(0, 0, dx, dy,v.BorderRoundWidth)
				ctx.Fill()
			}



			if v.PrimaryColor!=nil&&v.BorderWidth>0{
				ctx.SetLineWidth(v.BorderWidth)
				ctx.SetColor(v.PrimaryColor)
				ctx.DrawRoundedRectangle(0, 0, dx, dy,v.BorderRoundWidth)
				ctx.Stroke()
			}

			v.setBackground(img)
		} else {
			img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
			ctx := gg.NewContextForRGBA(img)
			if v.GetFont() != nil {
				ctx.SetFontFace(v.GetFont())
			}
			v.Backgrounder(ctx)
			v.setBackground(img)
		}

	}
	if v.getBackground() != nil {
		canvas.DrawImage(v.getBackground(), 0, 0)
	}

	if len(v.Children) > 0 {

		pw, ph := v.GetSize()
		offsetX, offsetY := 0.0, 0.0
		for i, l := range v.Children {
			//fmt.Println("i---",i,len(v.Children))
			if l == nil && v.GetView != nil {
				l = v.GetView(i)
				l.setParent(v)
				w, h := l.measure(l.getSettingSize())

				l.setSize(w, h)
				v.Children[i] = l
				v.measure(l.GetMeasureSize())
			}
			if l == nil {
				continue
			}
			x, y := l.getPosition()
			dx := int(x + offsetX)
			dy := int(y + offsetY)
			if v.Direction == Vertical {
				dy += int(v.ScrollLength)
			} else if v.Direction == Horizontal {
				dx += int(v.ScrollLength)
			}
			w, h := l.GetSize()

			if float64(dx)+w < -20 || float64(dx) > pw+20 || float64(dy)+h < -20 || float64(dy) > ph+20 {
				l.setRGBA(nil)
				//v.Children[i] = nil
			} else {
				if l.getRGBA() == nil {
					img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
					ctx := gg.NewContextForRGBA(img)
					if l.GetFont() != nil {
						ctx.SetFontFace(l.GetFont())
					}

					l.render(ctx)
					l.setRGBA(img)
				}

				l.setCurrentPosition(float64(dx),float64(dy))

				canvas.DrawImage(l.getRGBA(), dx, dy)
			}

			if v.Direction == Vertical {
				offsetY += (h)
			} else if v.Direction == Horizontal {
				offsetX += (w)
			}
		}
	}
	if v.getRGBA() == nil {

		w, h := v.GetSize()
		img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

		ctx := gg.NewContextForRGBA(img)
		if v.Font != nil {
			ctx.SetFontFace(v.Font)
		}
		if v.Drawer != nil {
			v.Drawer(ctx)
		} else {

			if v.MaxScrollLength < 0 {
				if v.Direction == Vertical {
					rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
					ctx.SetColor(color.RGBA{0x80, 0x80, 0x80, 0x99})
					scrollBarWidth := h * h / (h + math.Abs(v.MaxScrollLength))
					if scrollBarWidth < 10 {
						scrollBarWidth = 10
					}
					x0, y0 := w-10, (h-scrollBarWidth)*rate
					x1, y1 := x0+10, scrollBarWidth*(1-rate)+rate*h
					ctx.SetColor(colornames.Gray)
					ctx.DrawRectangle(x0, y0, x1-x0, y1-y0)
					ctx.Fill()

				} else if v.Direction == Horizontal {
					rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
					ctx.SetColor(color.RGBA{0x80, 0x80, 0x80, 0x99})
					scrollBarWidth := w * w / (w + math.Abs(v.MaxScrollLength))
					if scrollBarWidth < 10 {
						scrollBarWidth = 10
					}
					x0, y0 := (w-scrollBarWidth)*rate, h-10
					x1, y1 := scrollBarWidth*(1-rate)+rate*w, y0+10

					ctx.SetColor(colornames.Gray)
					ctx.DrawRectangle(x0, y0, x1-x0, y1-y0)
					ctx.Fill()
				}
			}

		}

		v.setRGBA(img)
	}

	canvas.DrawImage(v.getRGBA(), 0, 0)

}
func (v *View) Scroll(x, y float64, distance float64) bool {
	if v.Scroller != nil {
		return v.Scroller(x, y, distance)
	}

	if v.Children != nil {
		for _, c := range v.Children {
			if c == nil {
				continue
			}
			l, t := c.getPosition()
			width, height := c.GetSize()

			if x >= l && x <= l+width && y >= t && y <= t+height {

				if c.Scroll(x-l, y-t, distance) {
					return true
				}

			}

		}
	}
	if v.Direction != None {

		v.SetScroll(distance * 5.5)
		return true
	}

	return false
}
func (v *View) Event(x, y float64, action CursorType) bool {
	if v.Eventer != nil {
		return v.Eventer(x, y, action)
	}
	switch action {
	case Down:
		v.isDownScollBar = false
		if v.Direction == Vertical {
			w, h := v.GetSize()
			rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
			scrollBarWidth := h * h / (h + math.Abs(v.MaxScrollLength))
			if scrollBarWidth < ScollBarMinWidth {
				scrollBarWidth = ScollBarMinWidth
			}
			x0, y0 := w-ScollBarMinWidth, (h-scrollBarWidth)*rate
			x1, y1 := x0+ScollBarMinWidth, scrollBarWidth*(1-rate)+rate*h

			if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
				v.downPoint = []float64{x, y}
				v.isDownScollBar = true
				return true
			}
		} else if v.Direction == Horizontal {
			w, h := v.GetSize()
			rate := math.Abs(v.ScrollLength / v.MaxScrollLength)

			scrollBarWidth := w * w / (w + math.Abs(v.MaxScrollLength))
			if scrollBarWidth < ScollBarMinWidth {
				scrollBarWidth = ScollBarMinWidth
			}
			x0, y0 := (w-scrollBarWidth)*rate, h-ScollBarMinWidth
			x1, y1 := scrollBarWidth*(1-rate)+rate*w, y0+ScollBarMinWidth
			if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
				v.downPoint = []float64{x, y}
				v.isDownScollBar = true
				return true
			}
		}
		for _, c := range v.Children {
			if c == nil {
				continue
			}
			l, t := c.GetCurrentPosition()
			width, height := c.GetSize()

			if x >= l && x <= l+width && y >= t && y <= t+height {

				if c.Event(x-l, y-t, action) {
					v.cursorView = c
					return true
				}
			}

		}

		break
	case Up:

		if v.cursorView != nil {
			l, t := v.cursorView.GetCurrentPosition()
			v.cursorView.Event(x-l, y-t, action)
		} else {

			for _, c := range v.Children {
				if c == nil {
					continue
				}
				l, t := c.GetCurrentPosition()
				width, height := c.GetSize()

				if x >= l && x <= l+width && y >= t && y <= t+height {
					if c.Event(x-l, y-t, action){
						return true
					}
				}

			}

		}
		v.cursorView = nil
		v.downPoint = nil
		v.isDownScollBar = false
		if v.Clicker != nil {
			v.Clicker(v, x, y)
		}

		v.Window.SetFocus(false)
		if v.Focuser!=nil{


			if v.SetFocus(x,y,true){
				return true
			}
		}
		break
	case Hover:

		if v.cursorView != nil {

			//l, t := v.cursorView.getPosition()
			//width, height := v.cursorView.GetSize()
			v.cursorView.Event(x, y, action)
			//if x >= l && x <= l+width && y >= t && y <= t+height {
			//
			//} else {
			//	v.cursorView.Event(x-l, y-t, Over)
			//	v.cursorView = nil
			//}
			return true

		}
		for _, c := range v.Children {
			if c == nil {
				continue
			}

			l, t := c.GetCurrentPosition()
			width, height := c.GetSize()

			if x >= l && x <= l+width && y >= t && y <= t+height {

				if c.Event(x-l, y-t, action) {
					v.cursorView = c
					return true
				}
			}

		}

		if v.downPoint != nil && v.isDownScollBar {
			w, h := v.GetSize()
			if v.Direction == Vertical {
				if len(v.downPoint) < 3 {
					v.downPoint = append(v.downPoint, v.downPoint[1])
				}

				scrollBarWidth := h * h / (h + math.Abs(v.MaxScrollLength))
				if scrollBarWidth < ScollBarMinWidth {
					scrollBarWidth = ScollBarMinWidth
				}
				distance := v.MaxScrollLength * (y - v.downPoint[2]) / (h - scrollBarWidth)
				v.SetScroll(distance)
				v.downPoint[2] = y
			} else if v.Direction == Horizontal {
				if len(v.downPoint) < 3 {
					v.downPoint = append(v.downPoint, v.downPoint[0])
				}
				scrollBarWidth := w * w / (h + math.Abs(v.MaxScrollLength))
				if scrollBarWidth < ScollBarMinWidth {
					scrollBarWidth = ScollBarMinWidth
				}
				distance := v.MaxScrollLength * (x - v.downPoint[2]) / (w - scrollBarWidth)
				v.SetScroll(distance)
				v.downPoint[2] = x

			}
		}
		break
	case Out:
		if v.cursorView != nil {
			v.cursorView.Event(x, y, action)
		} else {
			for _, c := range v.Children {
				if c != nil {
					c.Event(x, y, action)
				}

			}
		}
		v.cursorView = nil
		v.downPoint = nil

		break
	}
	return false
}
func NewButtonView(v *View) Viewer {

	view := &TextView{v}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	var e error
	view.Font,e= gg.LoadFontFace(view.FontPath,float64( view.FontSize))
	if e!=nil{
		panic(e)
	}
	view.ShouldMeasure = true
	view.Backgrounder = view.DrawBackground
	if view.BackgroundColor==nil{
		view.BackgroundColor = colornames.Green
	}


	//view.TextColor = colornames.White
	//view.PaddingTop = 4
	//view.PaddingBottom = 4
	//view.PaddingLeft = 6
	//view.PaddingRight = 6
	return v
}
func NewTextView(v *View) Viewer {
	view := &TextView{v}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	var e error
	view.Font,e= gg.LoadFontFace(view.FontPath,float64( view.FontSize))
	if e!=nil{
		panic(e)
	}
	view.ShouldMeasure = true
	return v
}

type TextView struct {
	*View
}

func (v *TextView) MeasureSize(w, h float64) (float64, float64) {
	mw, mh := v.GetMeasureSize()
	if mw == AUTO_SCOLL_PARENT {
		mw = math.MaxFloat64
	}

	if mh == AUTO_SCOLL_PARENT {
		mh = math.MaxFloat64
	}
	str:=v.Title
	_, strw, strh := v.WordWrap(str, mw-v.PaddingLeft-v.PaddingRight, mh-v.PaddingTop-v.PaddingBottom)
	rw, rh := strw+v.PaddingLeft+v.PaddingRight, strh+v.PaddingTop+v.PaddingBottom
	lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())+v.PaddingBottom+v.PaddingTop
	if w == AUTO {
		w = rw
	}
	if w == FULL_PARENT {
		w = mw
	}
	if h == AUTO {
		h = rh
		if v.LineCount>0{
			h = lineHeight*float64(v.LineCount)
		}
	}
	if h == FULL_PARENT {
		h = mh
	}
	return w, h
}
func (v *TextView) DrawBackground(canvas *gg.Context) {
	w, h := v.GetSize()
	canvas.SetColor(v.BackgroundColor)
	canvas.DrawRoundedRectangle(0, 0, w, h, 5)
	canvas.Fill()
}
func (v *TextView) Draw(canvas *gg.Context) {
	w, h := v.GetSize()

	list, _, _ := v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)

	lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
	topOffset := float64(v.Font.Metrics().Ascent.Ceil())
	for i, l := range list {
		if v.TextColor == nil {
			v.TextColor = colornames.Black
		}

		if len(list)==1&&v.LineCount==1{
			canvas.SetColor(v.TextColor)
			canvas.DrawString(l, v.PaddingLeft, (h-lineHeight)/2+topOffset)
			canvas.Fill()
		}else{
			canvas.SetColor(v.TextColor)
			canvas.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
			canvas.Fill()
		}

	}

}
func (v *View) WordWrap(s string, width, height float64) ([]string, float64, float64) {
	if width == AUTO_SCOLL_PARENT {
		width = math.MaxFloat64
	}
	if height == AUTO_SCOLL_PARENT {
		height = math.MaxFloat64
	}
	results := make([]string, 0)
	ctx := gg.NewContext(1, 1)
	list := ctx.WordWrap(s, width*2)
	h := 0.0
	w := 0.0

	for _, l := range list {
		strs := []rune(l)
		start := 0
		for i := 0; i < len(strs); i++ {

			strw := font.MeasureString(v.Font, string(strs[start:i+1]))

			if float64(strw.Ceil()) > width {
				results = append(results, string(strs[start:i]))
				start = i
				h += float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
				w = width

			} else {
				w = math.Max(w, float64(strw.Ceil()))
			}
			if i == len(strs)-1 && string(strs[start:]) != "" {
				results = append(results, string(strs[start:]))
				h += float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())

			}
		}
	}

	if height > 0 && h > height {
		h = height
	}
	if h==0{
		h = float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
	}

	return results, w, h
}


func NewEditView(v *View) Viewer {
	view := &EditView{&TextView{v},0,0,0}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	var e error
	view.Font,e= gg.LoadFontFace(view.FontPath,float64( view.FontSize))
	if e!=nil{
		panic(e)
	}
	view.ShouldMeasure = true
	view.Focuser = view.Focus
	return v
}

type EditView struct {
	*TextView
	times int
	timestamp int64
	offsetX float64
}
func (v *EditView) Focus(x,y float64,b bool) bool{
	v.isFocus=b
	if b{
		v.timestamp = time.Now().Unix()
		v.times = 0
		go func(t int64) {
			for v.isFocus&&t==v.timestamp{
				v.RequestLayout()
				<-time.Tick(time.Millisecond*600)
				v.times++
			}
			v.times=0
		}(v.timestamp)
		w,h:=v.GetSize()
		x=x-v.PaddingLeft
		y=y-v.PaddingTop
		var list []string
		if v.LineCount!=1{
			list,_,_=v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)
			if len(list)==0{
				list=[]string{""}
			}
		}else{
			list=[]string{v.Title}
		}

		lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
		s:=""
		widthindex :=0
		line:=0
		minW:=math.MaxFloat64
		if y>0{
			line=int(y/lineHeight)


		}
		if v.LineCount==1{
			line = 0
		}
		if x>0{
			for i,l:=range []rune(list[line]){
				s+=string(l)
				advance:=float64(font.MeasureString(v.Font,s).Ceil())+v.offsetX

				if float64(advance)-x<minW{
					minW = math.Abs(float64(advance)-x)
					widthindex = i
				}

			}
		}else  {
			widthindex = -1
		}

		lineNum:=0
		if line>0{
			for i,l:=range list{
				if i>=line{
					break
				}
				lineNum+=len([]rune(l))
			}
		}
		v.cursorIndex= widthindex +1+lineNum
		if v.cursorIndex>len([]rune(v.Title)){
			v.cursorIndex =len([]rune(v.Title))
		}




	}
	v.RequestLayout()
	return b
}
func (v *EditView) Draw(canvas *gg.Context) {
	w, h := v.GetSize()
	var list []string
	if v.LineCount!=1{
		list, _, _ = v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)
	}else{
		list = []string{v.Title}
	}


	lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())


	topOffset := float64(v.Font.Metrics().Ascent.Ceil())
	if !v.isFocus{

		for i, l := range list {
			if v.TextColor == nil {
				v.TextColor = colornames.Black
			}
			if len(list)==1&&v.LineCount==1{
				canvas.SetColor(v.TextColor)
				canvas.DrawString(l, v.PaddingLeft, (h-lineHeight)/2+topOffset)
				canvas.Fill()
			}else{
				canvas.SetColor(v.TextColor)
				canvas.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
				canvas.Fill()
			}
			//canvas.SetColor(v.TextColor)
			//canvas.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
			//canvas.Fill()
		}
	}

	if v.isFocus{

		var strlist []string
		if v.LineCount!=1{
			strlist,_,_=v.WordWrap(string([]rune(v.Title)[:v.cursorIndex]),w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)

		}else{
			strlist=[]string{string([]rune(v.Title)[:v.cursorIndex])}
		}
		if len(strlist)==0{
			strlist =append(strlist,"")
		}
		advance:=font.MeasureString(v.Font,strlist[len(strlist)-1]).Ceil()
		if v.LineCount==1{

			strw:=font.MeasureString(v.Font,v.Title).Ceil()

			if float64(strw)> w-v.PaddingLeft-v.PaddingRight&&float64(advance)+v.offsetX>w-v.PaddingLeft-v.PaddingRight{
				if v.cursorIndex==len([]rune(v.Title)){
					v.offsetX =(w-v.PaddingLeft-v.PaddingRight)-float64(strw)
				}else{
					v.offsetX =(w-v.PaddingLeft-v.PaddingRight)-float64(advance)
				}

			}

			if float64(advance)+v.offsetX<0{
				v.offsetX = -float64(advance)
			}

		}



		for i, l := range list {
			if v.TextColor == nil {
				v.TextColor = colornames.Black
			}

			//canvas.SetColor(v.TextColor)
			//canvas.DrawString(l, v.PaddingLeft+v.offsetX, v.PaddingTop+topOffset+float64(i)*lineHeight)
			//canvas.Fill()

			if len(list)==1&&v.LineCount==1{
				canvas.SetColor(v.TextColor)
				canvas.DrawString(l, v.PaddingLeft+v.offsetX, (h-lineHeight)/2+topOffset)
				canvas.Fill()
			}else{
				canvas.SetColor(v.TextColor)
				canvas.DrawString(l, v.PaddingLeft+v.offsetX, v.PaddingTop+topOffset+float64(i)*lineHeight)
				canvas.Fill()
			}
		}


		if v.times%2==0{
			if v.AccentColor==nil{
				v.AccentColor = colornames.Black
			}

			if lineHeight==1{
				canvas.SetColor(v.AccentColor)
				canvas.DrawLine(v.PaddingLeft+float64(advance)+v.offsetX,h/2-lineHeight/2,v.PaddingLeft+float64(advance)+v.offsetX,h/2+lineHeight/2)
				canvas.Stroke()
			}else{
				canvas.SetColor(v.AccentColor)
				canvas.DrawLine(v.PaddingLeft+float64(advance)+v.offsetX,v.PaddingTop+lineHeight*float64(len(strlist)-1),v.PaddingLeft+float64(advance)+v.offsetX,lineHeight*float64(len(strlist))+v.PaddingTop)
				canvas.Stroke()
			}



		}
		//fmt.Println(strlist[len(strlist)-1])

	}

}


type ImageView struct {
	*View
	rgba image.Image
}
type ImageScaleType int

const (
	Fit   ImageScaleType = 0
	Cover ImageScaleType = 1
)

func NewImageView(v *View)Viewer {
	view := &ImageView{v,nil}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	view.ShouldMeasure = true
	return view
}
func (v *ImageView)loadingData()  {
	if v.Src != "" {
		if strings.HasPrefix(v.Src, "http") {
			resp, e := http.Get(v.Src)
			if e != nil {
				return
			}
			defer resp.Body.Close()
			v.rgba, _, e = image.Decode(resp.Body)
			if e != nil {
				return
			}

		}else{
			f, e := os.Open(v.Src)
			if e != nil {
				return
			}
			defer f.Close()
			v.rgba, _, e = image.Decode(f)
			if e != nil {
				return
			}
		}
	}
	if v.rgba!=nil{
		w, h := v.GetSize()
			s:=1.0
			if v.ScaleType==Fit{

				s=w/float64(v.rgba.Bounds().Dx())
				if h/float64(v.rgba.Bounds().Dy())<s{
					s = h/float64(v.rgba.Bounds().Dy())
				}
			}else if v.ScaleType==Cover{
				s=w/float64(v.rgba.Bounds().Dx())
				if h/float64(v.rgba.Bounds().Dy())>s{
					s = h/float64(v.rgba.Bounds().Dy())
				}
			}

			ctx:=gg.NewContext(int(w),int(h))
		    ctx.DrawRoundedRectangle(0,0,w,h,v.BorderRoundWidth)
		    ctx.Clip()
			ctx.Scale(s,s)
			ctx.DrawImageAnchored(v.rgba,int(w/(2*s)), int(h/(2*s)),0.5,0.5)

			//mask:=gg.NewContext(int(w),int(h))
			//mask.DrawRoundedRectangle(0,0,w,h,v.BorderRoundWidth)
			//mask.Clip()

			//ctx.SetMask(mask.AsMask())

			v.rgba = ctx.Image()

		v.RequestLayout()
	}
}
func (v *ImageView) Draw(canvas *gg.Context) {


	if v.rgba == nil {
		go v.loadingData()
	}
	if v.rgba!=nil{
		canvas.DrawImage(v.rgba, 0,0)
	}

}


func (v *ImageView) MeasureSize(w, h float64) (float64, float64) {

	return w, h
}
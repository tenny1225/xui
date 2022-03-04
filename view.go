package xui

import (
	"github.com/fogleman/gg"
	"github.com/google/uuid"
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
	render(canvas XCanvas)

	measure(w, h float64) (float64, float64)
	setSize(w, h float64)
	setWindow(w XWindow)
	getSettingSize() (float64, float64)
	getPosition() (float64, float64)
	setCurrentPosition(x, y float64)
	GetCurrentPosition() (float64, float64)
	getEvent() func(x, y float64, t CursorType) bool
	setPosition(x, y float64)
	getTitle() string
	isShouldMeasure() bool
	RequestLayout()
	setParent(view *View)
	init(w XWindow)
	getRGBA() *image.RGBA
	pushString(s string) bool
	backspace() bool
	addCursorIndex(i int) bool

	setRGBA(rgba *image.RGBA)
	getBackground() *image.RGBA
	setBackground(rgba *image.RGBA)
	getRenderRect() (float64, float64, float64, float64)
	setRenderPosition()
	getId()string

	Event(x, y float64, t CursorType) bool
	Scroll(x, y float64, distance float64) bool
	GetMeasureSize() (float64, float64)
	GetFont() font.Face
	SetScroll(s float64)
	GetSize() (float64, float64)
	GetDirection() Direction
	SetChildCount(c int)
	SetFocus(x, y float64, b bool) bool
	Recycle()
	DrawRect(c XCanvas,x,y ,w,h float64,color *color.RGBA)
	DrawImage(c XCanvas,x,y float64,img image.Image)
	DrawCircle(c XCanvas,x,y ,r float64,color *color.RGBA)
}

type View struct {
	id string
	Parent           *View
	BackgroundColor  *color.RGBA
	PrimaryColor     color.Color
	AccentColor      color.Color
	Children         []Viewer
	Width            float64
	Height           float64
	PaddingLeft      float64
	BorderWidth      float64
	BorderRoundWidth float64
	PaddingTop       float64
	PaddingRight     float64
	PaddingBottom    float64
	TextColor        color.Color
	Direction        Direction
	ScrollLength     float64
	MaxScrollLength  float64
	ScrollingDistance float64
	LineCount        int
	FontSize         int
	FontPath         string
	ScaleType        ImageScaleType
	Src              string
	currentLeft      float64
	currentTop       float64


	Left          float64
	Top           float64
	Right float64
	Bottom float64
	WidthRatio float64
	HeightRatio float64

	Title         string
	ShouldMeasure bool
	Drawer        func(canvas XCanvas)
	Clicker       func(v *View, x, y float64)
	Measurer      func(w, h float64) (float64, float64)
	Eventer       func(x, y float64, t CursorType) bool
	Scroller      func(x, y float64, distance float64) bool
	Focuser       func(x, y float64, b bool) bool
	GetView       func(i int) Viewer
	GetChildCount func() int

	Backgrounder func(canvas XCanvas)

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
func (v*View)setWindow(w XWindow){
	v.Window=w
	v.id = uuid.New().String()
}
func (v*View)getId()string  {
	return v.id
}
func (v *View) getRenderRect() (float64, float64, float64, float64) {
	w, h := v.GetSize()
	mw, mh := v.Window.GetSize()
	l, t := v.GetCurrentPosition()

	if v.Parent != nil {
		px, py, px1, py1 := v.Parent.getRenderRect()

		//pw,ph:=px1-px,py1-py

		return math.Max(px, l), math.Max(py, t), math.Min(float64(px1), l+w), math.Min(float64(py1), t+h)

	} else {
		return math.Max(0, l), math.Max(0, t), math.Min(float64(mw), l+w), math.Min(float64(mh), t+h)
	}

}
func (v *View) setRenderPosition() {
	left,top:=0.0,0.0
	l:=v
	for {

		left+=l.Left
		top+=l.Top


		if l.Parent!=nil{
			if l.Parent.Direction==Vertical{

				top+=l.Parent.ScrollLength
				list:=l.Parent.Children
				if list!=nil{
					for _,c:=range list{

						if c.getId()==l.id{
							break
						}
						_,h:=c.GetSize()
						_,t:=c.getPosition()
						top+=h+t


					}
				}


			}else if l.Parent.Direction==Horizontal{
				left+=l.Parent.ScrollLength
				list:=l.Parent.Children
				if list!=nil{
					for _,c:=range list{

						if c.getId()==l.id{
							break
						}
						w,_:=c.GetSize()
						l,_:=c.getPosition()
						left+=w+l
					}
				}

			}

			l = l.Parent

			left+=l.PaddingLeft
			top+=l.PaddingTop
		}else{
			break
		}
	}
	v.setCurrentPosition(left,top)
}
func (v *View) Recycle() {
	v.Window.UI(func() {
		v.background = nil
		v.rgba = nil
	})
}
func (v *View) setCurrentPosition(x, y float64) {
	v.currentLeft, v.currentTop = x, y
}
func (v *View) GetCurrentPosition() (float64, float64) {
	return v.currentLeft, v.currentTop
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
func (v *View)DrawCircle(c XCanvas,x,y ,r float64,color *color.RGBA){
	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()
	if x2>x1&&y2>y1{
		c.DrawCircle(x+l, y+t, r, color)
	}
}
func (v *View)DrawRect(c XCanvas,x,y ,w,h float64,color *color.RGBA){
	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()
	if x2>x1&&y2>y1{
		c.DrawRect(math.Max(x+l,x1), math.Max(y+t,y1), math.Min(w,x2-x1), math.Min(h,y2-y1), color)
	}

}
func (v *View)DrawImage(c XCanvas,x,y float64,img image.Image){
	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()
	//x2,y2=x2-v.PaddingRight,y2-v.PaddingBottom
	//x1,y1=x1-v.PaddingLeft,y1-v.PaddingTop
	if x2>x1&&y2>y1{
		c.DrawImageInRetangle(l+x, t+y, img, x1, y1, x2-x1, y2-y1)
	}

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

		if v.MaxScrollLength > 0 {
			v.MaxScrollLength = 0
		}
		if v.WidthRatio>0{
			width *=v.WidthRatio
		}
		if v.HeightRatio>0{
			height *=v.HeightRatio
		}
		width,height= width-v.Right, height-v.Bottom
		if  v.Direction == Vertical {
			v.MaxScrollLength = height - max

		} else if  v.Direction == Horizontal {
			v.MaxScrollLength = width - max
		}
		if v.MaxScrollLength>0{
			v.MaxScrollLength = 0
		}
		return  width,height
	}


}
func (v *View) setParent(view *View) {
	v.Parent = view

}
func (v *View) pushString(s string) bool {
	for _, l := range v.Children {
		if l.pushString(s) {
			return true
		}
	}
	if v.isFocus {
		v.Title = string([]rune(v.Title)[:v.cursorIndex]) + s + string([]rune(v.Title)[v.cursorIndex:])
		v.cursorIndex++
		v.Recycle()
		return true
	}
	return false
}
func (v *View) backspace() bool {
	for _, l := range v.Children {
		if l.backspace() {
			return true
		}
	}
	if v.isFocus {
		if v.Title != "" {
			list := make([]rune, 0)
			for i, l := range []rune(v.Title) {
				if i+1 != v.cursorIndex {
					list = append(list, l)
				}
			}
			v.Title = string(list)
			v.cursorIndex--
		}

		if v.cursorIndex < 0 {
			v.cursorIndex = 0
		}
		v.Recycle()
		return true
	}
	return false
}
func (v *View) addCursorIndex(i int) bool {
	for _, l := range v.Children {
		if l.addCursorIndex(i) {
			return true
		}
	}
	if v.isFocus {
		if v.Title != "" {
			v.cursorIndex += i
		}

		if v.cursorIndex < 0 {
			v.cursorIndex = 0
		} else if v.cursorIndex > len([]rune(v.Title)) {
			v.cursorIndex = len([]rune(v.Title))
		}
		v.Recycle()
		return true
	}
	return false
}
func (v *View) SetFocus(x, y float64, b bool) bool {
	for _, l := range v.Children {
		if l.SetFocus(x, y, b) {
			v.cursorView = l

			return true
		}
	}
	if !b {
		v.isFocus = b
	}
	if v.Focuser != nil {
		if  v.Focuser(x, y, b){
			return true
		}
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


}

func (v *View) GetMeasureSize() (float64, float64) {

	w := v.Width
	h := v.Height

	if (w == AUTO ) && v.Parent != nil {
		if v.Parent.Direction == Horizontal {
			w = AUTO_SCOLL_PARENT
		} else {
			w, _ = v.Parent.GetMeasureSize()
		}

	}
	if (h == AUTO) && v.Parent != nil {
		if v.Parent.Direction == Vertical {
			h = AUTO_SCOLL_PARENT
		} else {
			_, h = v.Parent.GetMeasureSize()
		}

	}
	if w==FULL_PARENT{
		if v.Parent!=nil{
			w,_ = v.Parent.GetMeasureSize()
		}
	}
	if h==FULL_PARENT{
		if v.Parent!=nil{
			_,h = v.Parent.GetMeasureSize()
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
	v.MaxScrollLength = 0
	v.ScrollingDistance=0
	v.ScrollLength = 0
	v.id = uuid.New().String()
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
func (v *View) render(canvas XCanvas) {

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
	if v.Direction!=None&&v.ScrollingDistance!=0{
		v.SetScroll(v.ScrollingDistance)
	}
	if v.ShouldMeasure {
		v.setSize(v.measure(v.getSettingSize()))
		v.setRenderPosition()
	}

	if v.Backgrounder == nil {
		if v.BackgroundColor != nil {
			x1, y1, x2, y2 := v.getRenderRect()
			if x2>x1&&y2>y1{
				canvas.DrawRect(x1, y1, x2-x1, y2-y1, v.BackgroundColor)
			}

		}

	} else {
		v.Backgrounder(canvas)

	}

	if len(v.Children) > 0 {

		for i, l := range v.Children {
			if l == nil && v.GetView != nil {
				l = v.GetView(i)
				l.init(v.Window)
				l.setParent(v)
				l.setSize(l.measure(l.getSettingSize()))
				v.Children[i] = l
			}
			if l == nil {
				continue
			}

			l.setRenderPosition()
			x, y := l.GetCurrentPosition()
			w, h := l.GetSize()
			x1,y1,x2,y2:=l.getRenderRect()


			if float64(x)+w < x1 || float64(x) > x2 || float64(y)+h < y1 || float64(y) > y2 {
				l.setCurrentPosition(float64(x), float64(y))
			} else {
				l.setCurrentPosition(float64(x), float64(y))
				l.render(canvas)

			}

		}
	}

	w, h := v.GetSize()

	if v.Drawer != nil {
		v.Drawer(canvas)
	} else {

		l, t := v.GetCurrentPosition()
		if v.MaxScrollLength < 0 {
			if v.Direction == Vertical {
				rate := math.Abs(v.ScrollLength / v.MaxScrollLength)

				scrollBarWidth := h * h / (h + math.Abs(v.MaxScrollLength))
				if scrollBarWidth < 10 {
					scrollBarWidth = 10
				}
				x0, y0 := w-10, (h-scrollBarWidth)*rate
				x1, y1 := x0+10, scrollBarWidth*(1-rate)+rate*h

				canvas.DrawRect(l+x0, t+y0, x1-x0, y1-y0, &colornames.Lightgrey)

			} else if v.Direction == Horizontal {
				rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
				scrollBarWidth := w * w / (w + math.Abs(v.MaxScrollLength))
				if scrollBarWidth < 10 {
					scrollBarWidth = 10
				}
				x0, y0 := (w-scrollBarWidth)*rate, h-10
				x1, y1 := scrollBarWidth*(1-rate)+rate*w, y0+10

				canvas.DrawRect(l+x0, t+y0, x1-x0, y1-y0,  &colornames.Lightgrey)
			}
		}

	}

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
			l, t := c.GetCurrentPosition()
			width, height := c.GetSize()

			if x >= l && x <= l+width && y >= t && y <= t+height {

				if c.Scroll(x-l, y-t, distance) {
					return true
				}

			}

		}
	}
	if v.Direction != None {



		if v.ScrollLength==v.MaxScrollLength&&distance<0{
			return false
		}
		if v.ScrollLength==0&&distance>0{
			return false
		}

		v.ScrollingDistance=distance
		v.id=uuid.New().String()
		go func(v *View,id string) {

			time.Sleep(time.Millisecond*100)
			if id==v.id{
				v.ScrollingDistance=0
			}


		}(v,v.id)
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
			l,t:=v.GetCurrentPosition()
			rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
			scrollBarWidth := h * h / (h + math.Abs(v.MaxScrollLength))
			if scrollBarWidth < ScollBarMinWidth {
				scrollBarWidth = ScollBarMinWidth
			}
			x0, y0 := w-ScollBarMinWidth+l, (h-scrollBarWidth)*rate+t
			x1, y1 := x0+ScollBarMinWidth+l, scrollBarWidth*(1-rate)+rate*h+t

			if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
				v.downPoint = []float64{x, y}
				v.isDownScollBar = true

				return true
			}
		} else if v.Direction == Horizontal {
			w, h := v.GetSize()
			rate := math.Abs(v.ScrollLength / v.MaxScrollLength)
			l,t:=v.GetCurrentPosition()
			scrollBarWidth := w * w / (w + math.Abs(v.MaxScrollLength))
			if scrollBarWidth < ScollBarMinWidth {
				scrollBarWidth = ScollBarMinWidth
			}
			x0, y0 := (w-scrollBarWidth)*rate+l, h-ScollBarMinWidth+t
			x1, y1 := scrollBarWidth*(1-rate)+rate*w+l, y0+ScollBarMinWidth+t
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
			//l, t := c.GetCurrentPosition()
			//width, height := c.GetSize()
			x1,y1,x2,y2:=c.getRenderRect()

			if x >= x1 && x <=x2 && y >= y1 && y <= y2 {

				if c.Event(x, y, action) {
					v.cursorView = c
					return true
				}
			}

		}

		break
	case Up:

		if v.cursorView != nil {

			v.cursorView.Event(x, y, action)
		} else {

			for _, c := range v.Children {
				if c == nil {
					continue
				}
				x1,y1,x2,y2:=c.getRenderRect()

				if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
					if c.Event(x, y, action) {
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
		if v.Focuser != nil {
			l,t:=v.GetCurrentPosition()

			if v.SetFocus(x-l, y-t, true) {
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

			x1,y1,x2,y2:=c.getRenderRect()

			if x >= x1 && x <= x2 && y >= y1 && y <= y2 {

				if c.Event(x, y, action) {
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
	view.Font, e = gg.LoadFontFace(view.FontPath, float64(view.FontSize))
	if e != nil {
		panic(e)
	}
	view.ShouldMeasure = true
	view.Backgrounder = view.DrawBackground
	if view.BackgroundColor == nil {
		view.BackgroundColor = &colornames.Green
	}
	return v
}
func NewTextView(v *View) Viewer {
	view := &TextView{v}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	var e error
	view.Font, e = gg.LoadFontFace(view.FontPath, float64(view.FontSize))

	if e != nil {
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
	str := v.Title
	_, strw, strh := v.WordWrap(str, mw-v.PaddingLeft-v.PaddingRight, mh-v.PaddingTop-v.PaddingBottom)
	rw, rh := strw+v.PaddingLeft+v.PaddingRight, strh+v.PaddingTop+v.PaddingBottom
	lineHeight := float64(v.Font.Metrics().Height.Ceil()+v.Font.Metrics().Descent.Ceil()) + v.PaddingBottom + v.PaddingTop
	if w == AUTO {
		w = rw
	}
	if w == FULL_PARENT {
		w = mw
	}
	if h == AUTO {
		h = rh
		if v.LineCount > 0 {
			h = lineHeight * float64(v.LineCount)
		}
	}
	if h == FULL_PARENT {
		h = mh
	}
	return w, h
}
func (v *TextView) DrawBackground(canvas XCanvas) {
	if v.background == nil {
		w, h := v.GetSize()
		v.background = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		ctx := gg.NewContextForRGBA(v.background)
		ctx.SetColor(v.BackgroundColor)
		ctx.DrawRoundedRectangle(0, 0, w, h, 5)
		ctx.Fill()
	}
	//if v.background!=nil{
	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()
	if y2>y1&&x2>x1{
		canvas.DrawImageInRetangle(l, t, v.background, x1, y1, x2-x1, y2-y1)
	}

	//}

}
func (v *TextView) Draw(canvas XCanvas) {
	if v.rgba == nil {
		w, h := v.GetSize()
		v.rgba = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		ctx := gg.NewContextForRGBA(v.rgba)
		ctx.SetFontFace(v.GetFont())

		list, _, _ := v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)

		lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
		topOffset := float64(v.Font.Metrics().Ascent.Ceil())
		for i, l := range list {
			if v.TextColor == nil {
				v.TextColor = colornames.Black
			}

			if len(list) == 1 && v.LineCount == 1 {
				ctx.SetColor(v.TextColor)
				ctx.DrawString(l, v.PaddingLeft, (h-lineHeight)/2+topOffset)
				ctx.Fill()
			} else {
				ctx.SetColor(v.TextColor)
				ctx.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
				ctx.Fill()
			}

		}

	}
	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()

	if y2>y1&&x2>x1{
		canvas.DrawImageInRetangle(l, t, v.rgba, x1, y1, x2-x1, y2-y1)

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
	if h == 0 {
		h = float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
	}

	return results, w, h
}

func NewEditView(v *View) Viewer {
	view := &EditView{&TextView{v}, 0, 0, 0}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	var e error
	view.Font, e = gg.LoadFontFace(view.FontPath, float64(view.FontSize))
	if e != nil {
		panic(e)
	}
	view.ShouldMeasure = true
	view.Focuser = view.Focus
	return v
}

type EditView struct {
	*TextView
	times     int
	timestamp int64
	offsetX   float64
}

func (v *EditView) Focus(x, y float64, b bool) bool {
	v.isFocus = b
	if b {
		v.timestamp = time.Now().Unix()
		v.times = 0
		go func(t int64) {
			for v.isFocus && t == v.timestamp {
				v.Recycle()
				<-time.Tick(time.Millisecond * 600)
				v.times++
			}
			v.times = 0
		}(v.timestamp)
		w, h := v.GetSize()
		x = x - v.PaddingLeft
		y = y - v.PaddingTop
		var list []string
		if v.LineCount != 1 {
			list, _, _ = v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)
			if len(list) == 0 {
				list = []string{""}
			}
		} else {
			list = []string{v.Title}
		}

		lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())
		s := ""
		widthindex := 0
		line := 0
		minW := math.MaxFloat64
		if y > 0 {
			line = int(y / lineHeight)

		}
		if v.LineCount == 1 {
			line = 0
		}
		if x > 0 {
			for i, l := range []rune(list[line]) {
				s += string(l)
				advance := float64(font.MeasureString(v.Font, s).Ceil()) + v.offsetX

				if float64(advance)-x < minW {
					minW = math.Abs(float64(advance) - x)
					widthindex = i
				}

			}
		} else {
			widthindex = -1
		}

		lineNum := 0
		if line > 0 {
			for i, l := range list {
				if i >= line {
					break
				}
				lineNum += len([]rune(l))
			}
		}
		v.cursorIndex = widthindex + 1 + lineNum
		if v.cursorIndex > len([]rune(v.Title)) {
			v.cursorIndex = len([]rune(v.Title))
		}

	}
	v.Window.UI(func() {
		v.Recycle()
	})
	return b
}
func (v *EditView) Draw(canvas XCanvas) {
	if v.rgba == nil {
		w, h := v.GetSize()
		v.rgba = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		ctx := gg.NewContextForRGBA(v.rgba)
		ctx.SetFontFace(v.GetFont())
		var list []string
		if v.LineCount != 1 {
			list, _, _ = v.WordWrap(v.Title, w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)
		} else {
			list = []string{v.Title}
		}

		lineHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())

		topOffset := float64(v.Font.Metrics().Ascent.Ceil())
		if !v.isFocus {

			for i, l := range list {
				if v.TextColor == nil {
					v.TextColor = colornames.Black
				}
				if len(list) == 1 && v.LineCount == 1 {
					ctx.SetColor(v.TextColor)
					ctx.DrawString(l, v.PaddingLeft, (h-lineHeight)/2+topOffset)
					ctx.Fill()
				} else {
					ctx.SetColor(v.TextColor)
					ctx.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
					ctx.Fill()
				}
				//canvas.SetColor(v.TextColor)
				//canvas.DrawString(l, v.PaddingLeft, v.PaddingTop+topOffset+float64(i)*lineHeight)
				//canvas.Fill()
			}
		}

		if v.isFocus {

			var strlist []string
			if v.LineCount != 1 {
				strlist, _, _ = v.WordWrap(string([]rune(v.Title)[:v.cursorIndex]), w-v.PaddingLeft-v.PaddingRight, h-v.PaddingTop-v.PaddingBottom)

			} else {
				strlist = []string{string([]rune(v.Title)[:v.cursorIndex])}
			}
			if len(strlist) == 0 {
				strlist = append(strlist, "")
			}
			advance := font.MeasureString(v.Font, strlist[len(strlist)-1]).Ceil()
			if v.LineCount == 1 {

				strw := font.MeasureString(v.Font, v.Title).Ceil()

				if float64(strw) > w-v.PaddingLeft-v.PaddingRight && float64(advance)+v.offsetX > w-v.PaddingLeft-v.PaddingRight {
					if v.cursorIndex == len([]rune(v.Title)) {
						v.offsetX = (w - v.PaddingLeft - v.PaddingRight) - float64(strw)
					} else {
						v.offsetX = (w - v.PaddingLeft - v.PaddingRight) - float64(advance)
					}

				}

				if float64(advance)+v.offsetX < 0 {
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

				if len(list) == 1 && v.LineCount == 1 {
					ctx.SetColor(v.TextColor)
					ctx.DrawString(l, v.PaddingLeft+v.offsetX, (h-lineHeight)/2+topOffset)
					ctx.Fill()
				} else {
					ctx.SetColor(v.TextColor)
					ctx.DrawString(l, v.PaddingLeft+v.offsetX, v.PaddingTop+topOffset+float64(i)*lineHeight)
					ctx.Fill()
				}
			}

			if v.times%2 == 0 {
				if v.AccentColor == nil {
					v.AccentColor = colornames.Black
				}

				if lineHeight == 1 {
					ctx.SetColor(v.AccentColor)
					ctx.DrawLine(v.PaddingLeft+float64(advance)+v.offsetX, h/2-lineHeight/2, v.PaddingLeft+float64(advance)+v.offsetX, h/2+lineHeight/2)
					ctx.Stroke()
				} else {
					ctx.SetColor(v.AccentColor)
					ctx.DrawLine(v.PaddingLeft+float64(advance)+v.offsetX, v.PaddingTop+lineHeight*float64(len(strlist)-1), v.PaddingLeft+float64(advance)+v.offsetX, lineHeight*float64(len(strlist))+v.PaddingTop)
					ctx.Stroke()
				}

			}
			//fmt.Println(strlist[len(strlist)-1])

		}
	}

	l, t := v.GetCurrentPosition()
	x1, y1, x2, y2 := v.getRenderRect()
	if y2>y1&&x2>x1{
		canvas.DrawImageInRetangle(l, t, v.rgba, x1, y1, x2-x1, y2-y1)
	}



}

type ImageView struct {
	*View
	isLoading bool
}
type ImageScaleType int

const (
	Fit   ImageScaleType = 0
	Cover ImageScaleType = 1
)

func NewImageView(v *View) Viewer {
	view := &ImageView{v, false}
	view.Drawer = view.Draw
	view.Measurer = view.MeasureSize
	view.ShouldMeasure = true
	return view
}
func (v *ImageView) loadingData() {
	defer func() {
		v.isLoading=false
	}()
	if v.Src != "" {
		var img image.Image
		if strings.HasPrefix(v.Src, "http") {
			resp, e := http.Get(v.Src)
			if e != nil {
				return
			}
			defer resp.Body.Close()
			img, _, e = image.Decode(resp.Body)
			if e != nil {
				return
			}

		} else {
			f, e := os.Open(v.Src)
			if e != nil {
				return
			}
			defer f.Close()
			img, _, e = image.Decode(f)
			if e != nil {
				return
			}
		}
		w, h := v.GetSize()
		s := 1.0
		if v.ScaleType == Fit {

			s = w / float64(img.Bounds().Dx())
			if h/float64(img.Bounds().Dy()) < s {
				s = h / float64(img.Bounds().Dy())
			}
		} else if v.ScaleType == Cover {
			s = w / float64(img.Bounds().Dx())
			if h/float64(img.Bounds().Dy()) > s {
				s = h / float64(img.Bounds().Dy())
			}
		}

		ctx := gg.NewContext(int(w), int(h))
		ctx.DrawRoundedRectangle(0, 0, w, h, v.BorderRoundWidth)
		ctx.Clip()
		ctx.Scale(s, s)
		ctx.DrawImageAnchored(img, int(w/(2*s)), int(h/(2*s)), 0.5, 0.5)

		//mask:=gg.NewContext(int(w),int(h))
		//mask.DrawRoundedRectangle(0,0,w,h,v.BorderRoundWidth)
		//mask.Clip()

		//ctx.SetMask(mask.AsMask())

		v.rgba = ctx.Image().(*image.RGBA)
	}


}
func (v *ImageView) Draw(canvas XCanvas) {

	if v.rgba == nil &&!v.isLoading{
		v.isLoading=true
		go v.loadingData()
	}
	if v.rgba != nil {
		l, t := v.GetCurrentPosition()
		x1, y1, x2, y2 := v.getRenderRect()
		canvas.DrawImageInRetangle(l, t, v.rgba, x1, y1, x2-x1, y2-y1)

	}

}

func (v *ImageView) MeasureSize(w, h float64) (float64, float64) {

	return w, h
}


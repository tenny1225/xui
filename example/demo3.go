package main

import (
	"github.com/fogleman/gg"
	"github.com/tenny1225/xui"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
)

type ClipPage struct {
	xui.BasePage
}

func (p *ClipPage) GetContentView() xui.Viewer {
	return &xui.View{
		Width: xui.FULL_PARENT,
		Height:xui.FULL_PARENT,
		Children: []xui.Viewer{

			&xui.View{
				Height:300,
				Top:200,
				Width:xui.FULL_PARENT,
				Direction:xui.Vertical,
				GetChildCount: func() int {
					return 20
				},
				GetView: func(i int) xui.Viewer {
					return &xui.View{
						Width:xui.FULL_PARENT,
						Height:60,
						PaddingLeft:  30,
						PaddingRight: 30,
						Top:10,
						Direction:xui.Horizontal,
						Children: []xui.Viewer{
							NewSlideView(&xui.View{
								Top:20,
								Width:        xui.FULL_PARENT,
								Height:       30,
								Right:100,
								FontPath:     "OPPOSans-L.ttf",
								Title:"123",
								FontSize:     12,
								PrimaryColor: &colornames.Gray,
								AccentColor:  &colornames.Orange,
							}),
							xui.NewButtonView(&xui.View{
								Top:10,
								LineCount:1,
								Height:30,
								Left:10,
								PaddingLeft:8,
								PaddingRight:8,
								PaddingTop:4,
								PaddingBottom:4,
								FontPath:     "OPPOSans-L.ttf",
								Title:"123",
								FontSize:12,
							}),
						},
					}
				},
			},
			&xui.View{
				Width:        xui.FULL_PARENT,
				Height:       200,
				Top:0,
				PaddingLeft:  30,
				PaddingRight: 30,
				BackgroundColor: &colornames.Black,
			},
		},
	}
}
func main() {
	ctx := xui.NewXContext()
	ctx.Run(func() {
		w := xui.NewWindow("测试", 500, 500, 100, 100, false, ctx)
		w.AddRoute("clip", &ClipPage{})
		w.StartPage("clip", nil, false)
	})
}

type SlideView struct {
	*xui.View
	lineHeight  float64
	total       float64
	currents    []float64
	titles      []string
	focusIndex  int
	isDown bool
	rgba        *image.RGBA
	spaceHeight float64
}

func NewSlideView(v *xui.View) xui.Viewer {
	view := &SlideView{View: v}
	view.Drawer = view.Draw
	var e error
	view.Font, e = gg.LoadFontFace(view.FontPath, float64(view.FontSize))
	if e != nil {
		panic(e)
	}
	//view.Measurer = view.MeasureSize
	view.ShouldMeasure = true
	view.lineHeight = 8
	view.Eventer = view.Event
	view .total  =100
	view.focusIndex =-1
	view.currents = []float64{30,80}
	view.titles = []string{"30","80"}
	view.spaceHeight  =4
	return view
}
func (v *SlideView) Draw(canvas xui.XCanvas) {
	w, h := v.GetSize()
	w,h = w-v.PaddingLeft-v.PaddingRight,h-v.PaddingTop-v.PaddingBottom

	fontHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())


	top := (h - v.lineHeight - fontHeight - v.spaceHeight*2) / 2+v.PaddingTop

	//l, t := v.GetCurrentPosition()
	v.DrawRect(canvas, v.PaddingLeft, top+v.PaddingTop, w, v.lineHeight, v.PrimaryColor.(*color.RGBA))


	for i,l:=range v.currents{
		x:=l*w/v.total
		r:=8.0
		if i==v.focusIndex {
			r = 10.0
		}
		v.DrawCircle(canvas, x+v.PaddingLeft, top+v.lineHeight/2, r, v.AccentColor.(*color.RGBA))
	}
	if len(v.currents)==2{
		x1:=v.currents[0]*w/v.total
		x2:=v.currents[1]*w/v.total
		x:=x1
		width:=x2-x1
		if x1>x2{
			x =x2
			width = x1-x2
		}
		v.DrawRect(canvas, x+v.PaddingLeft, top, width, v.lineHeight, v.AccentColor.(*color.RGBA))
	}
	if v.rgba==nil{

		v.rgba = image.NewRGBA(image.Rect(0,0,int(w+v.PaddingLeft+v.PaddingRight),int(h+v.PaddingTop+v.PaddingBottom)))
		ctx:=gg.NewContextForRGBA(v.rgba)

		ctx.SetFontFace(v.Font)
		ctx.SetColor(colornames.Gray)
		top:=(h-v.PaddingTop-v.PaddingBottom - v.lineHeight - fontHeight - v.spaceHeight*2) / 2
		for i,l:=range v.titles{
			x:=v.currents[i]*w/v.total
			s:=l
			strw,_:=ctx.MeasureString(s)
			ctx.DrawString(s,x-float64(strw)/2+v.PaddingLeft,top+v.lineHeight+v.spaceHeight*2+fontHeight/2+v.PaddingTop)
		}
	}

	v.DrawImage(canvas,0,0,v.rgba)




}
func (v *SlideView) Event(x, y float64, action xui.CursorType) bool {
	w, h := v.GetSize()
	w,h = w-v.PaddingLeft-v.PaddingRight,h-v.PaddingTop-v.PaddingBottom
	l,t:=v.GetCurrentPosition()
	fontHeight := float64(v.Font.Metrics().Height.Ceil() + v.Font.Metrics().Descent.Ceil())


	top := (h - v.lineHeight - fontHeight - v.spaceHeight*2) / 2+v.PaddingTop
	switch action {
	case xui.Down:

		for i,n:=range v.currents{
			rect:=image.Rect(int(n*w/v.total+v.PaddingLeft+l)-5,int(top+t)-5,int(n*w/v.total+v.PaddingLeft+8+l)+5,int(top+8+t)+5)
			if int(x)>=rect.Min.X&&int(x)<=rect.Max.X&&int(y)>=rect.Min.Y&&int(y)<=rect.Max.Y{
				v.focusIndex =i
				v.isDown=true

				return true
			}
		}
	case xui.Hover:
		if v.focusIndex >=0&&v.isDown{
			v.currents[v.focusIndex] = (x-l-v.PaddingLeft)*v.total/w
			if v.currents[v.focusIndex]<0{
				v.currents[v.focusIndex] = 0

			}
			if v.currents[v.focusIndex]>v.total{
				v.currents[v.focusIndex] = v.total

			}
			v.Recycle()
			return true
		}
		for i,n:=range v.currents{
			rect:=image.Rect(int(n*w/v.total+v.PaddingLeft+l)-5,int(top+t)-5,int(n*w/v.total+v.PaddingLeft+8+l)+5,int(top+8+t)+5)
			if int(x)>=rect.Min.X&&int(x)<=rect.Max.X&&int(y)>=rect.Min.Y&&int(y)<=rect.Max.Y{
				v.focusIndex =i
				return true
			}
		}
		v.focusIndex =-1

	case xui.Up:
		fallthrough
	case xui.Out:
		v.focusIndex =-1
		v.isDown=false
	}
	return false
}
func (v *SlideView) Recycle() {
	v.Window.UI(func() {
		v.rgba = nil
	})
}
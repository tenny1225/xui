package main

import (
	"github.com/tenny1225/xui"
	"golang.org/x/image/colornames"
)


type DemoPage struct {
	xui.BasePage
}

func (p *DemoPage) GetContentView() xui.Viewer {
	return &xui.View{
		Width:  xui.FULL_PARENT,
		Height: xui.FULL_PARENT,
		Direction:xui.Vertical,

		Children: []xui.Viewer{
			&xui.View{
				Width:xui.FULL_PARENT,
				Height:xui.FULL_PARENT,
				BackgroundColor:&colornames.Blue,
				Bottom:100,
				Direction:xui.Vertical,
				Children: []xui.Viewer{
					&xui.View{
						Top:10,
						Width:xui.FULL_PARENT,
						Height:100,
						BackgroundColor:&colornames.Red,
					},
					&xui.View{
						Top:10,
						Width:xui.FULL_PARENT,
						Height:100,
						BackgroundColor:&colornames.Yellow,
					},
					&xui.View{
						Top:10,
						Width:xui.FULL_PARENT,
						Height:100,
						BackgroundColor:&colornames.Orange,
					},
					&xui.View{
						Top:10,
						Width:xui.FULL_PARENT,
						Height:100,
						BackgroundColor:&colornames.Palegoldenrod,
					},
					&xui.View{
						Top:10,
						Width:xui.FULL_PARENT,
						Height:100,
						BackgroundColor:&colornames.Gray,
					},
				},
			},
			&xui.View{
				Width:xui.FULL_PARENT,
				Height:100,
				BackgroundColor:&colornames.Black,
			},

		},
	}
}

func main() {
	ctx := xui.NewXContext()
	ctx.Run(func() {
		w := xui.NewWindow("测试", 500, 500, ctx)
		w.AddRoute("start", &DemoPage{})
		w.StartPage("start", nil, false)
	})
}

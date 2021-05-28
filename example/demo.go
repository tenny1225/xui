package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"xui"
)

type TestPage struct {
	xui.BasePage
}
func (p *TestPage) GetContentView() xui.Viewer {
	return &xui.View{
		Width:500,
		Height:100,
		BackgroundColor:colornames.Gold,
		Direction:xui.Vertical,
		Children: []xui.Viewer{
			&xui.View{
				Title:"test",
				Width:  xui.FULL_PARENT,
				Height:80,
				Top:8,
				Children: []xui.Viewer{


					xui.NewTextView(&xui.View{
						Height:30,
						Left:100,
						Top:8,
						FontPath:        "OPPOSans-M.ttf",
						FontSize:        16,
						Title:"美洲豹的一天",
					}),


				},
			},
			&xui.View{
				Width:  xui.FULL_PARENT,
				Height:80,
				Top:8,
				Children: []xui.Viewer{

					xui.NewImageView(&xui.View{
						Width:60,
						Height:60,
						Left:20,
						Top:8,
						ScaleType:xui.Cover,
						BorderRoundWidth:30,

						Src:"https://fuss10.elemecdn.com/0/6f/e35ff375812e6b0020b6b4e8f9583jpeg.jpeg",
					}),
					xui.NewTextView(&xui.View{
						Height:30,
						Left:100,
						Top:8,
						FontPath:        "OPPOSans-M.ttf",
						FontSize:        16,
						Title:"豹子狩猎",
					}),
					xui.NewTextView(&xui.View{
						Height:30,
						Left:100,
						Width:340,
						Top:38,
						FontPath:        "OPPOSans-L.ttf",
						FontSize:        12,
						Title:"又叫美洲虎，是现存第三大的猫科动物。体重35—150千克，最大亚种雄性亚马孙美洲豹平均体重为98千克，咬力可达1250磅.",
					}),
					xui.NewButtonView(&xui.View{
						Height:40,
						LineCount:1,
						Left:430,
						Top:8,
						FontPath:        "OPPOSans-M.ttf",
						FontSize:        12,
						TextColor:colornames.Green,
						BackgroundColor:colornames.White,

						Title:"查看详情",
						Clicker: func(v *xui.View, x, y float64) {
							fmt.Println("click")
						},
					}),
					&xui.View{
						Width:400,
						Left:50,
						Height:1,
						Top:76,
						BackgroundColor:colornames.Lightgrey,
					},
				},
			},


		},
	}
}
func main()  {
	ctx := xui.NewXContext()
	ctx.Run(func() {
		w := xui.NewWindow("测试", 500, 500, ctx)
		w.AddRoute("test", &TestPage{})
		w.StartPage("test", nil, false)
	})
}

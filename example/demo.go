package main

import "xui"

type TestPage struct {
	xui.BasePage
}
func (p *TestPage) GetContentView() xui.Viewer {
	return xui.NewButtonView(&xui.View{
		Top:10,
		Left:10,
		FontPath:        "OPPOSans-M.ttf",
		FontSize:        15,
		Title:"hello world",
		PaddingLeft:8,
		PaddingTop:8,
		PaddingRight:8,
		PaddingBottom:8,
	})
}
func main()  {
	ctx := xui.NewXContext()
	ctx.Run(func() {
		w := xui.NewWindow("测试", 500, 500, ctx)
		w.AddRoute("test", &TestPage{})
		w.StartPage("test", nil, false)
	})
}

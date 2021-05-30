package main

import (
	"golang.org/x/image/colornames"
	"github.com/tenny1225/xui"
)

type StartPage1 struct {
	xui.BasePage
}
func (p *StartPage1)Create(data map[string]interface{}){

}
func (p *StartPage1) GetContentView() xui.Viewer {
	return &xui.View{
		Children: []xui.Viewer{
			&xui.View{
				BackgroundColor:&colornames.Red,
				Width:  xui.FULL_PARENT,
				Height: 80,
				Children: []xui.Viewer{
					xui.NewTextView(&xui.View{
						Height:40,
						BackgroundColor:&colornames.Black,
						Left:            20,
						Top:             20,
						LineCount:1,
						FontPath:        "OPPOSans-M.ttf",
						FontSize:        20,
						PaddingLeft:     8,
						PaddingTop:      8,
						PaddingRight:    8,
						PaddingBottom:   8,
						TextColor:       colornames.White,
						Title:           "今日新闻",
					}),
					xui.NewEditView(&xui.View{
						Width:           180,
						Left:            250,
						Top:             20,
						Height:40,
						LineCount:1,
						FontPath:        "OPPOSans-L.ttf",
						FontSize:        15,
						PaddingLeft:     8,
						PaddingTop:      8,
						PaddingRight:    8,
						PaddingBottom:   8,
						BackgroundColor: &colornames.White,
						BorderRoundWidth:10,
						TextColor:       colornames.Black,
						Title:           "请输入",
					}),
					xui.NewButtonView(&xui.View{
						Left:            440,
						Top:             20,
						Height:40,
						LineCount:1,
						FontPath:        "OPPOSans-M.ttf",
						FontSize:        15,
						PaddingLeft:     8,
						PaddingTop:      8,
						PaddingRight:    8,
						PaddingBottom:   8,
						BorderRoundWidth:10,
						BackgroundColor:&colornames.White,
						TextColor:       colornames.Black,
						Title:           "搜索",
						Clicker: func(v *xui.View, x, y float64) {
											v.Window.GetTexture2Image()
						},
					}),
				},
			},
			&xui.View{
				Width:  xui.FULL_PARENT,
				Height: 420,
				Top:80,
				Direction:xui.Vertical,
				Children: []xui.Viewer{

					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
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
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									//fmt.Println("click")
									v.Window.StartPage("end",nil,false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
						Children: []xui.Viewer{

							xui.NewImageView(&xui.View{
								Width:60,
								Height:60,
								Left:20,
								Top:8,
								ScaleType:xui.Cover,
								BorderRoundWidth:30,
								Src:"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
							}),
							xui.NewTextView(&xui.View{
								Height:30,
								Left:100,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        16,
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									v.Window.StartPage("end", map[string]interface{}{
										"url":"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
									},false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
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
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									//fmt.Println("click")
									v.Window.StartPage("end",nil,false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
						Children: []xui.Viewer{

							xui.NewImageView(&xui.View{
								Width:60,
								Height:60,
								Left:20,
								Top:8,
								ScaleType:xui.Cover,
								BorderRoundWidth:30,
								Src:"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
							}),
							xui.NewTextView(&xui.View{
								Height:30,
								Left:100,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        16,
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									v.Window.StartPage("end", map[string]interface{}{
										"url":"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
									},false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
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
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									//fmt.Println("click")
									v.Window.StartPage("end",nil,false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
						Children: []xui.Viewer{

							xui.NewImageView(&xui.View{
								Width:60,
								Height:60,
								Left:20,
								Top:8,
								ScaleType:xui.Cover,
								BorderRoundWidth:30,
								Src:"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
							}),
							xui.NewTextView(&xui.View{
								Height:30,
								Left:100,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        16,
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									v.Window.StartPage("end", map[string]interface{}{
										"url":"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
									},false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
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
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									//fmt.Println("click")
									v.Window.StartPage("end",nil,false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
						Children: []xui.Viewer{

							xui.NewImageView(&xui.View{
								Width:60,
								Height:60,
								Left:20,
								Top:8,
								ScaleType:xui.Cover,
								BorderRoundWidth:30,
								Src:"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
							}),
							xui.NewTextView(&xui.View{
								Height:30,
								Left:100,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        16,
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									v.Window.StartPage("end", map[string]interface{}{
										"url":"https://fuss10.elemecdn.com/2/11/6535bcfb26e4c79b48ddde44f4b6fjpeg.jpeg",
									},false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
					&xui.View{
						Width:  xui.FULL_PARENT,
						Height:80,
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
								Title:"美洲豹的一天",
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
								Height:30,
								Left:430,
								Top:8,
								FontPath:        "OPPOSans-M.ttf",
								FontSize:        12,
								TextColor:colornames.Green,
								BackgroundColor:&colornames.White,
								Title:"查看详情",
								Clicker: func(v *xui.View, x, y float64) {
									//fmt.Println("click")
									v.Window.StartPage("end",nil,false)
								},
							}),
							&xui.View{
								Width:400,
								Left:50,
								Height:1,
								Top:76,
								BackgroundColor:&colornames.Lightgrey,
							},
						},
					},
                },
			},
		},
	}
}

type EndPage1 struct {
	xui.BasePage
	url string
}

func (p*EndPage1)Create(data map[string]interface{})  {
	if data!=nil{
		p.url = data["url"].(string)
	}else{
		p.url = "https://fuss10.elemecdn.com/0/6f/e35ff375812e6b0020b6b4e8f9583jpeg.jpeg"
	}
}
func (p *EndPage1) GetContentView() xui.Viewer {
	return &xui.View{
		Width:  500,
		Height: 500,
		Children: []xui.Viewer{
			xui.NewImageView(&xui.View{
				Top:20,
				Left:             50,
				Width:            400,
				Height:           200,
				Src:              p.url,
				ScaleType:        xui.Cover,
			}),
			xui.NewTextView(&xui.View{
				Left:             50,
				Top:240,
				LineCount:1,
				FontPath:        "OPPOSans-M.ttf",
				FontSize:        20,
				Title:"美洲豹的一天",
			}),
			xui.NewTextView(&xui.View{
				Left:             50,
				Top:280,
			    Width:400,
				FontPath:        "OPPOSans-L.ttf",
				FontSize:        14,
				Title:"美洲豹，学名：Panthera onca (Linnaeus, 1758)，又叫美洲虎，是现存第三大的猫科动物。体重35—150千克，最大亚种雄性亚马孙美洲豹平均体重为98千克，咬力可达1250磅。是生活在中南美洲的一种大型猫科动物。它身上的花纹比较像豹，但整个身体的形状又更接近于虎。在猫科动物中，美洲豹的体型仅次于狮、虎。野外寿命约18年。人工饲养的历史达20多年",
			}),
			xui.NewTextView(&xui.View{
				Left:             50,
				Top:430,
				LineCount:1,
				FontPath:        "OPPOSans-M.ttf",
				FontSize:        16,
				TextColor:colornames.Green,
				Title:"返回<",
				Clicker: func(v *xui.View, x, y float64) {
					p.GetWindow().PopPage()
				},
			}),
		},
	}
}

func main() {
	ctx := xui.NewXContext()
	ctx.Run(func() {
		w := xui.NewWindow("测试", 500, 500, ctx)
		w.AddRoute("start", &StartPage1{})
		w.AddRoute("end", &EndPage1{})
		w.StartPage("start", nil, false)
	})
}

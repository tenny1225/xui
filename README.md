### 给golang编写一个UI框架
废话不多说，先看效果
![点击查看详情会跳转第二个页面](https://upload-images.jianshu.io/upload_images/874510-ca79fa951beeb229.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![点击返回会回到上一个页面](https://upload-images.jianshu.io/upload_images/874510-55040a65ff2ff818.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
代码在github上https://github.com/tenny1225/xui，以上例子放在了example目录下
go语言目前没有官方版本的UI库,如果想要看到某种结果必须使用打印或者文件化，目前也没有官方支持的UI库，所以自己在闲暇时刻实现了了一个简陋的框架，底层用的是glfw和opengl2，这里有官方实现的go绑定库https://github.com/go-gl
本人之前看过谷歌gallery2的源码，也尝试过fulltter的开发，所以整体设计是模仿它们的实现方式；为了快速实现，所有的绘制全部基于image，所以性能上有些不尽人意；目前可以使用控件有ImageView，TextView，ButtonView，EditView，以及横向和竖向的滚动条，现在还存在很多bug（本人实力有限，代码也写的乱）

下面是一个简单的测试例子，显示一个hello world的按钮
```
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

```

![2021-05-27 12-50-33屏幕截图.png](https://upload-images.jianshu.io/upload_images/874510-b78eaf5598f5ddc0.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


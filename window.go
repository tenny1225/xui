package xui

import (
	"context"

	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"runtime"
	"sync"
)

type XWindow interface {
	Close()
	GetContext() context.Context
	GetSize() (int, int)
	AddRoute(string, Page)
	StartPage(r string,data map[string]interface{},clearall bool)
	PopPage()
	SetFocus(b bool)
	getMQ()func(w*xwindow)
	pushMQ(f func(w*xwindow))
	RequestLayout()
	Reload()
}

func NewWindow(name string,w, h int, c XContext) XWindow {
	ctx, cancel := context.WithCancel(c.getCtx())
	c.addWait()

	win := &xwindow{
		name:      name,
		width:     w,
		height:    h,
		router:    make(map[string]Page),
		context:   c,
		ctx:       ctx,
		cancel:    cancel,
		pages:     make([]Page,0),
		mqQuenues: make([]func(w*xwindow),0),
		locker:    sync.RWMutex{},
	}
	go func(w XWindow, ctx context.Context) {
		<-ctx.Done()
		w.Close()
	}(win, ctx)
	return win
}

type xwindow struct {
	name         string
	context      XContext
	ctx          context.Context
	cancel       context.CancelFunc
	window       *glfw.Window
	width        int
	height       int
	router       map[string]Page
	pages        []Page
	mqQuenues    []func(w*xwindow)
	canvas       XCanvas
	isClose      bool
	leaveAnimate Differentiator
	enterAnimate Differentiator
	rgba         image.Image
	rootView     *View
	locker       sync.RWMutex

}
func (w *xwindow) getMQ()func(w*xwindow){
	w.locker.Lock()
	defer w.locker.Unlock()
	if len(w.mqQuenues)>0{
		f:=w.mqQuenues[0]
		if(len(w.mqQuenues)==1){
			w.mqQuenues = make([]func(w*xwindow),0)
		}else{
			w.mqQuenues = w.mqQuenues[1:len(w.mqQuenues)-1]
		}

		return f
	}
	return nil
}
func (w *xwindow) pushMQ(f func(w*xwindow)){
	w.locker.Lock()
	defer w.locker.Unlock()
	w.mqQuenues = append(w.mqQuenues,f)
}
func (w *xwindow) Close() {
	if !w.isClose {
		w.cancel()
		w.isClose = true
		w.context.windowDone()
	}


}
func (w *xwindow)RequestLayout(){
	w.pushMQ(func(w *xwindow) {
		//w.rootView.Children[0] = w.pages[len(w.pages)-1].GetContentView()
		w.rootView.init(w)
	})
}
func (w *xwindow)Reload() {
	w.pushMQ(func(w *xwindow) {
		w.rootView.Children[0] = w.pages[len(w.pages)-1].GetContentView()
		w.rootView.init(w)
	})
}
func (w *xwindow) GetSize() (int, int) {
	return w.window.GetSize()
}
func (w *xwindow) GetContext() context.Context {
	return w.ctx
}
func (w *xwindow) AddRoute(str string, p Page) {
	p.setWindow(w)
	w.router[str] = p
}
func (w *xwindow) StartPage(r string,data map[string]interface{},isclear bool) {
	go w.pageTo(r,data,isclear)
}
func (w *xwindow)PopPage(){

	w.pushMQ(func(w*xwindow) {
		if len(w.pages)==0{
			return
		}
		old:=w.pages[len(w.pages)-1]
		if len(w.pages)==1{
			old.Stop()
			old.Destroy()
			w.Close()
			w.pages =make([]Page,0)
			return
		}

		var rgba image.Image = nil
		if w.rootView!=nil&&w.rootView.rgba!=nil{
			ctx:=gg.NewContextForRGBA(w.rootView.rgba)
			rgba=ctx.Image()
		}
		w.pages =w.pages[:len(w.pages)-1]
		curretPage:=w.pages[len(w.pages)-1]
		curretPage.Active()

		width, height := w.window.GetSize()
		w.rootView = &View{
			Width:  float64(width),
			Height: float64(height),
			Title:  "parent",
			Children: []Viewer{
				curretPage.GetContentView(),
			},
		}
		w.rootView.init(w)
		w.rgba=rgba
		if old!=nil{
			w.leaveAnimate = old.GetPopDifferentiator()
			w.enterAnimate= curretPage.GetRecoverDifferentiator()
			old.Stop()
			old.Destroy()
		}
	})

}
func (w *xwindow)SetFocus(b bool){
	w.rootView.SetFocus(-1,-1,b)
	w.pushMQ(func(w *xwindow) {
		w.rootView.RequestLayout()
	})
}
func (w *xwindow)initCurentPage(r string,data map[string]interface{})  {
	w.pushMQ(func(w*xwindow) {
		var old Page =nil
		if len(w.pages)>0{
			old = w.pages[len(w.pages)-1]
		}

		if w.rootView!=nil&&w.rootView.rgba!=nil{
			ctx:=gg.NewContextForRGBA(w.rootView.rgba)
			w.rgba=ctx.Image()
		}
		currentPage:=w.router[r]
		//currentPagePtrType := reflect.TypeOf(w.router[r]) //获取call的指针的reflect.Type
		//
		//currentPageTrueType := currentPagePtrType.Elem() //获取type的真实类型
		//
		//currentPagePtrValue := reflect.New(currentPageTrueType) //返回对象的指针对应的reflect.Value
		//
		//currentPage := currentPagePtrValue.Interface().(Page)





		w.pages = append(w.pages,currentPage)
		currentPage.Create(data)
		currentPage.Active()
		width, height := w.window.GetSize()
		w.rootView = &View{
			Width:  float64(width),
			Height: float64(height),
			Title:  "parent",
			Children: []Viewer{
				currentPage.GetContentView(),
			},
		}
		w.rootView.init(w)
		if old!=nil{
			w.leaveAnimate = old.GetQueneDifferentiator()
			w.enterAnimate=currentPage.GetEnterDifferentiator()
			old.Stop()
		}
	})
}
func (w *xwindow) pageTo(r string,data map[string]interface{},b bool) error {
	if w.window == nil {
		runtime.LockOSThread()
		if err := glfw.Init(); err != nil {
			panic(err)
		}
		defer glfw.Terminate()
		//glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)

		window, err := glfw.CreateWindow(w.width, w.height, w.name, nil, nil)
		if err != nil {
			panic(err)
		}
		w.window = window

		window.Restore()
		window.MakeContextCurrent()
		if err := gl.Init(); err != nil {
			panic(err)
		}

		gl.ClearColor(0.5, 0.5, 0.5, 0.0)
		gl.ClearDepth(1)
		gl.DepthFunc(gl.LEQUAL)
		ambient := []float32{0.5, 0.5, 0.5, 1}
		diffuse := []float32{1, 1, 1, 1}
		lightPosition := []float32{-5, 5, 10, 0}
		gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
		gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
		gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
		gl.Enable(gl.LIGHT0)

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		//gl.Frustum(-1, 1, -1, 1, 1.0, 5.0)
		gl.Ortho(-3, 3, -3, 3, -3.0, 100.0)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		canvas := NewCanvas(w)
		w.initCurentPage(r,nil)
		window.SetFocusCallback(func(win *glfw.Window, focused bool) {
			if !focused{
				w.rootView.SetFocus(-1,-1,false)
			}

		})
		window.SetCloseCallback(func(win *glfw.Window) {
			w.Close()
		})
		window.SetCharCallback(func(win *glfw.Window, char rune) {
			w.rootView.pushString(string(char))
			w.rootView.RequestLayout()
		})
		window.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key==glfw.KeyBackspace&&action== glfw.Release {
				w.rootView.backspace()
				w.rootView.RequestLayout()
			}else if key==glfw.KeyLeft&&action== glfw.Release {
				w.rootView.addCursorIndex(-1)
				w.rootView.RequestLayout()
			}else if key==glfw.KeyRight&&action== glfw.Release {
				w.rootView.addCursorIndex(1)
				w.rootView.RequestLayout()
			}
		})

		window.SetSizeCallback(func(win *glfw.Window, width int, height int) {

			gl.Viewport(0, 0, int32(width), int32(height))
			w.rootView.Width = float64(width)
			w.rootView.Height = float64(height)
			w.rootView.init(w)
		})
		window.SetScrollCallback(func(win *glfw.Window, x float64, y float64) {
			if w.rootView.currentPoint != nil&&w.enterAnimate==nil&&w.leaveAnimate==nil {
				w.rootView.Scroll(w.rootView.currentPoint[0], w.rootView.currentPoint[1], y)
			}

		})

		window.SetCursorPosCallback(func(win *glfw.Window, xpos float64, ypos float64) {
			if w.enterAnimate==nil&&w.leaveAnimate==nil{
				width, height := win.GetSize()

				w.rootView.Width, w.rootView.Height = float64(width), float64(height)
				x, y := win.GetCursorPos()
				l, t := w.rootView.getPosition()

				if x >= l && x <= l+float64(width) && y >= t && y <= t+float64(height) {
					w.rootView.Event(x, y, Hover)
					w.rootView.currentPoint = []float64{x, y}
				} else {
					w.rootView.currentPoint = nil
				}
			}


		})
		window.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			if button == glfw.MouseButtonLeft &&w.enterAnimate==nil&&w.leaveAnimate==nil{
				x, y := win.GetCursorPos()
				width, height := win.GetSize()
				l, t := w.rootView.getPosition()
				ct := CursorType(0)
				if action == glfw.Press {
					ct = Down
				} else if action == glfw.Release {
					ct = Up
				}
				if x >= l && x <= l+float64(width) && y >= t && y <= t+float64(height) {

					w.rootView.Event(x, y, ct)
				}

			}
		})

		for !w.isClose&&!window.ShouldClose() {
			for{
				if fun:=w.getMQ();fun!=nil{
					fun(w)
					continue
				}
				break
			}
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			gl.ClearColor(1, 1, 1, 0)
			width, height := window.GetSize()
			px, py := window.GetCursorPos()
			if px < 0 || px > float64(width) || py < 0 || py > float64(height) {
				w.rootView.Event(px, py, Out)
			}

			if w.rootView.rgba == nil {
				img := image.NewRGBA(image.Rect(0, 0, int(w.rootView.Width), int(w.rootView.Height)))
				ctx := gg.NewContextForRGBA(img)
				w.rootView.render(ctx)
				w.rootView.setRGBA(img)
			}

			if w.rgba!=nil&&w.leaveAnimate!=nil{

				offsetX,offsetY:=w.leaveAnimate.GetCurrent()
				canvas.DrawImage(offsetX, offsetY,w.rgba)
			}

			offsetX,offsetY:=0.0,0.0
			if w.enterAnimate!=nil{
				offsetX,offsetY = w.enterAnimate.GetCurrent()
			}
			canvas.DrawImage(offsetX, offsetY, w.rootView.getRGBA())

			window.SwapBuffers()
			glfw.PollEvents()

			if  w.enterAnimate!=nil&&w.enterAnimate.IsDone(){
				w.enterAnimate = nil

			}
			if  w.leaveAnimate!=nil&&w.leaveAnimate.IsDone(){
				w.leaveAnimate = nil
				w.rgba=nil
			}
			if w.enterAnimate!=nil{
				w.enterAnimate.Calc()

			}
			if w.leaveAnimate!=nil{
				w.leaveAnimate.Calc()
			}
		}

	} else {

		if b{
			for _,l:=range w.pages {
				l.Stop()
				l.Destroy()
			}
			w.pages =make([]Page,0)
		}

		w.initCurentPage(r,data)
	}

	return nil
}

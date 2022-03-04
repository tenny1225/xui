package xui

import (
	"context"
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
	UI(f func())
	RequestLayout()
	Reload()
	GetTexture2Image()*image.RGBA
}

func NewWindow(name string,w, h,x,y int,resizable bool, c XContext) XWindow {
	ctx, cancel := context.WithCancel(c.getCtx())
	c.addWait()

	win := &xwindow{
		name:      name,
		x:x,
		y:y,
		width:     w,
		height:    h,
		router:    make(map[string]Page),
		context:   c,
		ctx:       ctx,
		cancel:    cancel,
		pages:     make([]Page,0),
		mqQuenues: make([]func(w*xwindow),0),
		locker:    sync.RWMutex{},
		resizable:resizable,
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
	x int
	y int
	router       map[string]Page
	pages        []Page
	mqQuenues    []func(w*xwindow)
	canvas       XCanvas
	isClose      bool
	leaveAnimate Differentiator
	enterAnimate Differentiator
	rgba         image.Image
	locker       sync.RWMutex
	resizable bool

}
func (w *xwindow)UI(f func()){
	w.locker.Lock()
	defer w.locker.Unlock()
	w.mqQuenues = append(w.mqQuenues, func(w *xwindow) {
		f()
	})
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
		w.window.Destroy()
		w.cancel()
		w.isClose = true
		w.context.windowDone()
	}


}
func (w *xwindow)RequestLayout(){
	w.pushMQ(func(w *xwindow) {
		//w.rootView.Children[0] = w.pages[len(w.pages)-1].GetContentView()
		w.pages[len(w.pages)-1].getRoot().init(w)
	})
}
func (w *xwindow)Reload() {
	w.pushMQ(func(w *xwindow) {
		//w.pages[len(w.pages)-1].getRoot().Children[0] = w.pages[len(w.pages)-1].GetContentView()
		//w.pages[len(w.pages)-1].getRoot().init(w)
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
	w.rgba = w.GetTexture2Image()

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

		w.pages =w.pages[:len(w.pages)-1]
		curretPage:=w.pages[len(w.pages)-1]
		curretPage.Active()


		if old!=nil{
			w.leaveAnimate = old.GetPopDifferentiator()
			w.enterAnimate= curretPage.GetRecoverDifferentiator()
			old.Stop()
			old.Destroy()
		}
	})

}
func (w *xwindow)SetFocus(b bool){
	w.pages[len(w.pages)-1].getRoot().SetFocus(-1,-1,b)
}
func (w *xwindow)initCurentPage(r string,data map[string]interface{})  {
	w.pushMQ(func(w*xwindow) {
		var old Page =nil
		if len(w.pages)>0{
			old = w.pages[len(w.pages)-1]
		}


		currentPage:=w.router[r]
		w.pages = append(w.pages,currentPage)
		currentPage.Create(data)
		currentPage.Active()
		width, height := w.window.GetSize()
		root:= &View{
			Width:  float64(width),
			Height: float64(height),
			Title:  "parent",
			Children: []Viewer{
				currentPage.GetContentView(),
			},
		}
		root.init(w)
		currentPage.setRoot(root)

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
		if !w.resizable{
			glfw.WindowHint(glfw.Resizable, glfw.False)
		}
		glfw.WindowHint(glfw.ScaleToMonitor, glfw.False)
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)


		window, err := glfw.CreateWindow(w.width, w.height, w.name, nil, nil)
		if err != nil {
			panic(err)
		}
		window.SetPos(w.x,w.y)
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
		gl.LightModeli(gl.FRONT, gl.AMBIENT_AND_DIFFUSE)
		gl.Enable( gl.COLOR_MATERIAL)
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
				w.pages[len(w.pages)-1].getRoot().SetFocus(-1,-1,false)
			}

		})
		window.SetCloseCallback(func(win *glfw.Window) {
			w.Close()
		})
		window.SetCharCallback(func(win *glfw.Window, char rune) {
			w.pages[len(w.pages)-1].getRoot().pushString(string(char))
			//w.rootView.RequestLayout()
		})
		window.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			rootView:=w.pages[len(w.pages)-1].getRoot()
			if key==glfw.KeyBackspace&&action== glfw.Release {
				rootView.backspace()
				//w.rootView.RequestLayout()
			}else if key==glfw.KeyLeft&&action== glfw.Release {
				rootView.addCursorIndex(-1)
				//w.rootView.RequestLayout()
			}else if key==glfw.KeyRight&&action== glfw.Release {
				rootView.addCursorIndex(1)
				//w.rootView.RequestLayout()
			}
		})

		window.SetSizeCallback(func(win *glfw.Window, width int, height int) {
			rootView:=w.pages[len(w.pages)-1].getRoot()
			gl.Viewport(0, 0, int32(width), int32(height))
			rootView.Width = float64(width)
			rootView.Height = float64(height)
			rootView.init(w)
		})
		window.SetScrollCallback(func(win *glfw.Window, x float64, y float64) {
			rootView:=w.pages[len(w.pages)-1].getRoot()
			if rootView.currentPoint != nil&&w.enterAnimate==nil&&w.leaveAnimate==nil {
				rootView.Scroll(rootView.currentPoint[0], rootView.currentPoint[1], y*15)
				//w.offsetX+=y*22
			}

		})

		window.SetCursorPosCallback(func(win *glfw.Window, xpos float64, ypos float64) {
			if w.enterAnimate==nil&&w.leaveAnimate==nil{
				//rootView:=w.pages[len(w.pages)-1].getRoot()
				x, y := win.GetCursorPos()
				//l, t := w.pages[len(w.pages)-1].getRoot().getPosition()
				//width,height:=rootView.GetSize()

				//if x >= l && x <= l+float64(width) && y >= t && y <= t+float64(height) {
				//	w.pages[len(w.pages)-1].getRoot().Event(x, y, Hover)
				//	w.pages[len(w.pages)-1].getRoot().currentPoint = []float64{x, y}
				//} else {
				//	w.pages[len(w.pages)-1].getRoot().currentPoint = nil
				//}
				w.pages[len(w.pages)-1].getRoot().Event(x, y, Hover)
				w.pages[len(w.pages)-1].getRoot().currentPoint = []float64{x, y}
			}


		})

		window.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
			if button == glfw.MouseButtonLeft &&w.enterAnimate==nil&&w.leaveAnimate==nil{
				x, y := win.GetCursorPos()
				width, height := win.GetSize()
				l, t := w.pages[len(w.pages)-1].getRoot().getPosition()
				ct := CursorType(0)
				if action == glfw.Press {
					ct = Down
				} else if action == glfw.Release {
					ct = Up
				}
				if ct==Down{
					if x >= l && x <= l+float64(width) && y >= t && y <= t+float64(height) {

						w.pages[len(w.pages)-1].getRoot().Event(x, y, ct)
					}
				}else{
					w.pages[len(w.pages)-1].getRoot().Event(x, y, ct)
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


			if  w.leaveAnimate!=nil&&w.rgba!=nil{
				dx,dy:=w.leaveAnimate.GetCurrent()

				width,height:=w.GetSize()
				rx,ry,rx1,ry1:=dx,dy,dx+float64(width),dy+float64(height)
				if rx<0{
					rx = 0
				}
				if ry<0{
					ry=0
				}
				if rx>float64(width){
					rx = float64(width)
				}
				if ry>float64(height){
					ry = float64(height)
				}

				if rx1<0{
					rx1 = 0
				}
				if ry1<0{
					ry1=0
				}
				if rx1>float64(width){
					rx1 = float64(width)
				}
				if ry1>float64(height){
					ry1 = float64(height)
				}
				if rx1>rx&&ry1>ry{
					canvas.Save()
					canvas.DrawImageInRetangle(dx,dy,w.rgba,rx,ry,rx1-rx,ry1-ry)
					canvas.Restore()
				}

			}

			dx,dy:=0.0,0.0
			if  w.enterAnimate!=nil{
				dx,dy = w.enterAnimate.GetCurrent()
			}
			canvas.SetTranslate(dx,dy)
			canvas.Save()

			w.pages[len(w.pages)-1].getRoot().render(canvas)



			canvas.Restore()
			canvas.SetTranslate(0,0)





			window.SwapBuffers()
			glfw.PollEvents()

			if  w.enterAnimate!=nil&&w.enterAnimate.IsDone(){
				w.enterAnimate = nil

			}
			if  w.leaveAnimate!=nil&&w.leaveAnimate.IsDone(){
				w.leaveAnimate = nil
				//w.rgba = nil

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
		w.rgba = w.GetTexture2Image()

		w.initCurentPage(r,data)
	}

	return nil
}
func (w*xwindow)GetTexture2Image()*image.RGBA  {
	width,height:=w.GetSize()
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 4);
	rgba:=image.NewRGBA(image.Rect(0,0,width,height))
	gl.ReadPixels(0,0,int32(width),int32(height),gl.RGBA,gl.UNSIGNED_BYTE,gl.Ptr(rgba.Pix))
	dist:=image.NewRGBA(image.Rect(0,0,width,height))
	for i:=0;i< rgba.Bounds().Max.X;i++{
		for j:=0;j< rgba.Bounds().Max.Y;j++{
			c:=rgba.At(i,j)
			dist.Set(i,rgba.Bounds().Max.Y-j,c)
		}
	}
	return dist
}
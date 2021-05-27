package xui

type Page interface {
	Create(data map[string]interface{})
	Active()
	Stop()
	Destroy()
	GetContentView() Viewer
	GetWindow( )XWindow
	GetEnterDifferentiator( )Differentiator
	GetQueneDifferentiator( )Differentiator
	GetPopDifferentiator( )Differentiator
	GetRecoverDifferentiator( )Differentiator
	setWindow(window XWindow)

}
type BasePage struct {
	window XWindow

}

func (p *BasePage) Create(data map[string]interface{}) {

}

func (*BasePage) Active() {

}

func (*BasePage) Stop() {

}

func (*BasePage) Destroy() {

}

func (*BasePage) GetContentView() Viewer {
	return nil
}
func (p *BasePage) setWindow(window XWindow) {
	p.window = window
}
func (p *BasePage) GetWindow( )XWindow{
	return p.window
}
func (p *BasePage) GetEnterDifferentiator( )(Differentiator){
	w,_:=p.GetWindow().GetSize()
	return NewlineDifferentiator(float64(w),0,0,0)
}
func (p *BasePage) GetQueneDifferentiator( )(Differentiator){
	w,_:=p.GetWindow().GetSize()
	return NewlineDifferentiator(0,-float64(w),0,0)
}
func (p *BasePage) GetPopDifferentiator( )(Differentiator){
	w,_:=p.GetWindow().GetSize()
	return NewlineDifferentiator(0,float64(w),0,0)
}
func (p *BasePage) GetRecoverDifferentiator( )Differentiator{
	w,_:=p.GetWindow().GetSize()
	return NewlineDifferentiator(-float64(w),0,0,0)
}
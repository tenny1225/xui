package xui

type Differentiator interface {
	Calc()
	IsDone()bool
	GetCurrent()(float64,float64)
	Clone()Differentiator
}

type lineDifferentiator struct {
	startX float64
	endX float64
	currentX float64

	startY float64
	endY float64
	currentY float64
}

func NewlineDifferentiator(startX,endX,startY,endY float64) Differentiator {
	return &lineDifferentiator{
		startX:startX,
		endX:endX,
		currentX:startX,
		startY:startY,
		endY:endY,
		currentY:startY,

	}
}
func (l*lineDifferentiator)Clone()Differentiator{
	return &lineDifferentiator{
		startX:l.startX,
		endX:l.endX,
		currentX:l.startX,
		startY:l.startY,
		endY:l.endY,
		currentY:l.startY,
	}
}
func (l*lineDifferentiator) GetCurrent() (float64,float64){
	return l.currentX,l.currentY
}
func (l*lineDifferentiator) Calc() {

	if l.currentX !=l.endX{
		l.currentX += (l.endX-l.startX)*0.1
		if l.endX>l.startX&&l.currentX>l.endX{
			l.currentX = l.endX
		}
		if l.endX<l.startX&&l.currentX<l.endX{
			l.currentX = l.endX
		}
	}

	if l.currentY !=l.endY{
		l.currentY += (l.endY-l.startY)*0.1
		if l.endY>l.startY&&l.currentY>l.endY{
			l.currentY = l.endY
		}
		if l.endY<l.startY&&l.currentY<l.endY{
			l.currentY = l.endY
		}
	}
}

func (l*lineDifferentiator) IsDone() bool {
	return l.currentX==l.endX&&l.currentY==l.endY
}

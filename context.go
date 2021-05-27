package xui

import (
	"context"
	"sync"
	"time"
)
type XContext interface {
	Close()
	Run(f func())
	getCtx()context.Context
	addWait()
	windowDone()

}
type xcontext struct {
	isClose chan bool
	cancel  context.CancelFunc
	ctx context.Context
	wait *sync.WaitGroup
	count int
}

func (c *xcontext)addWait(){
	c.wait.Add(1)
	c.count++
}
func (c *xcontext)windowDone(){
	c.wait.Done()
}
func (c *xcontext)getCtx()context.Context{
	return c.ctx
}
func (c *xcontext) Close() {
	c.cancel()
	for i:=0;i<c.count;i++{
		c.wait.Done()
	}
}
func (c *xcontext) Run(f func()) {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		f()
	}()
	<-time.Tick(time.Second*2)
	c.wait.Wait()
}

func NewXContext() *xcontext {
	return &xcontext{isClose: make(chan bool),wait:&sync.WaitGroup{}}
}

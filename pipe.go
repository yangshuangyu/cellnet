package cellnet

import (
	"reflect"
)

type EventPipe interface {
	AddQueue() EventQueue

	Start()

	Stop(int)

	Wait() int
}

type evPipe struct {
	qarray []*evQueue

	arrayLock  bool
	exitSignal chan int
}

func (self *evPipe) AddQueue() EventQueue {

	if self.arrayLock {
		panic("Pipe already start, can not addqueue any more")
	}

	q := newEventQueue()

	self.qarray = append(self.qarray, q)

	return q
}

func (self *evPipe) Start() {

	// 开始后, 不能修改数组
	self.arrayLock = true

	go func() {

		cases := make([]reflect.SelectCase, len(self.qarray))

		// 按照队列(peer)数量开始做case
		for i, q := range self.qarray {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(q.queue)}
		}

		for {

			if index, value, ok := reflect.Select(cases); ok {

				self.qarray[index].CallData(value.Interface())
			}

		}

	}()

}

func (self *evPipe) Stop(result int) {
	self.exitSignal <- result
}

func (self *evPipe) Wait() int {
	return <-self.exitSignal
}

func NewEventPipe() EventPipe {
	return &evPipe{
		qarray:     make([]*evQueue, 0),
		exitSignal: make(chan int),
	}
}

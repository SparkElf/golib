package promise

import "sync"

//Scheduler The scheduler is an experimental attempt to give you control over any phase of an asynchronous workflow.
type Scheduler struct {
	Head        *Promise
	Cur         *Promise
	Done        chan bool //notice that it is a bufferd channel.
	Interrupted bool
	RWMutex     *sync.RWMutex
}

//Interrupt Interrupt the whole workflow.
func (this *Scheduler) Interrupt() {
	this.RWMutex.Lock()
	this.Interrupted = true
	this.RWMutex.Unlock()
}

//IsInterrupted Return the value of the variable interrupt for Scheduler.
//
//It can only ensure that what is read is not a corrupt variable.If you need more code synchronization,
//you can get the lock and write your own synchronization code.
func (this *Scheduler) IsInterrupted() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.Interrupted
}

//NewScheduler return a new scheduler.
func NewScheduler() *Scheduler {
	s := new(Scheduler)
	return s
}

//InitScheduler Init the scheduler by the first function to excute.
func (this *Scheduler) InitScheduler(run func() (interface{}, error)) *Promise {
	this.Done = make(chan bool, 1)
	p := &Promise{
		scheduler: this,
		okHandler: func(interface{}) (interface{}, error) {
			return run()
		},
	}
	this.Head = p
	return p
}

//Await Block the main thread until the asynchronous workflow returns results and
//stops working.
func (this *Scheduler) Await() (interface{}, error) {
	<-this.Done
	return this.Cur.next.V, this.Cur.next.E
}

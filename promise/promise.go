package promise

import (
	"sync"
	"time"
)

type OkHandler = func(interface{}) (interface{}, error)
type ErrHandler = func(error) (interface{}, error)
type ThenHandler = func(interface{}, error) (interface{}, error)
type InterruptHandler = func(*Promise) (interface{}, error)

//Promise The task entity.
//E and V represent the errors and values of the last task
//After the task is done,the error and value generated by this task will be assigned to the next promise.
type Promise struct {
	E error
	V interface{}

	errHandler       ErrHandler
	okHandler        OkHandler
	thenHandler      ThenHandler
	interruptHandler InterruptHandler

	next      *Promise
	timeout   time.Duration
	scheduler *Scheduler

	Interrupted bool
	RWMutex     *sync.RWMutex
}

func (this *Promise) GetScheduler() *Scheduler {
	return this.scheduler
}

func (this *Promise) GetNext() *Promise {
	return this.next
}

//NewPromise Create a new scheduler with run function and then return the first promise.
/**
 *!Notice:
 *!This function will create a new scheduler,not just a promise,do not use it casually.
 */
func NewPromise(run func() (interface{}, error)) *Promise {
	return NewScheduler().InitScheduler(run)
}

//SetTimeout Set timeout of a promise.
func (this *Promise) SetTimeout(t time.Duration) *Promise {
	this.timeout = t
	return this
}

//Interrupt Interrupt a single task without affecting the entire workflow.
func (this *Promise) Interrupt() {
	this.RWMutex.Lock()
	this.Interrupted = true
	this.RWMutex.Unlock()
}

//IsInterrupted Return the value of the variable interrupt for promise.
//The RLock ensure what is read is not a corrupt variable.If you need
//synchronize code blocks,you can get the lock and write your own synchronization code.
func (this *Promise) IsInterrupted() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.Interrupted
}

//Then Add callback to handle values and errors
func (this *Promise) Then(then ThenHandler) *Promise {
	p := &Promise{
		thenHandler: then,
		scheduler:   this.scheduler,
	}
	this.next = p
	return p
}

//WithInterrupt Add callback to handle values and errors,and it provides finer grained interrupt control.
/**
 * !Notice:
 */
//In go, the main thread cannot actively interrupt the sub thread,so if you don't write your own code to poll the
//interrupt flag,the sub thread will not be interrupted before the execution is completed.
func (this *Promise) WithInterrupt(i InterruptHandler) *Promise {
	p := &Promise{
		interruptHandler: i,
		scheduler:        this.scheduler,
	}
	this.next = p
	return p
}

//OnSuccess Add callback to handle values if task success.
//It introduced to make programming more semantic and it only works with OnError.
func (this *Promise) OnSuccess(okHandler OkHandler) *Promise {
	if this.thenHandler != nil || this.interruptHandler != nil {
		p := &Promise{
			okHandler: okHandler,
			scheduler: this.scheduler,
		}
		this.next = p
		return p
	} else if this != this.scheduler.Head {
		this.okHandler = okHandler
	}
	return this
}

//OnError Add callback to handle errors if task fail.
//It introducd to make programming more semantic and it only works with OnSuccess.
func (this *Promise) OnError(errHandler ErrHandler) *Promise {
	if this.thenHandler != nil || this.interruptHandler != nil {
		p := &Promise{
			errHandler: errHandler,
			scheduler:  this.scheduler,
		}
		this.next = p
		return p
	} else if this != this.scheduler.Head {
		this.errHandler = errHandler
	}
	return this
}

//StartTasks Start to excute callbacks.
func (this *Promise) StartTasks() *Scheduler {
	this.next = &Promise{} //用作保存结果
	go func() {
		s := this.scheduler
		for i := s.Head; i.next != nil; i = i.next {
			s.Cur = i
			done := make(chan bool)

			if s.IsInterrupted() {
				i.Interrupt()
				i.next.E = &Err{code: INTERRUPT}
				s.Done <- true
				return
			}

			if i.timeout != 0 {
				go i.run(&done)
				select {
				case <-done:
				case <-time.After(i.timeout):
					i.next.E = &Err{code: TIMEOUT}
					i.Interrupt()
				}
			} else {
				go i.run(&done)
				<-done
			}
		}
		s.Done <- true
	}()
	return this.scheduler
}

//run Excute the current callback function.
/**
 * !Notice:
 */
//then fucntion and withInterrupt function work alone,onSuccess and onError work togother.
func (this *Promise) run(done *chan bool) {
	var v interface{}
	var err error
	if this.thenHandler != nil {
		v, err = this.thenHandler(this.V, this.E)
	} else if this.interruptHandler != nil {
		v, err = this.interruptHandler(this)
	} else if this.E != nil {
		v, err = this.errHandler(this.E)
	} else {
		v, err = this.okHandler(this.V)
	}
	this.next.V = v
	this.next.E = err
	*done <- true
}

package main

import (
	"github.com/SparkElf/golib/promise"
	"time"
)

func main() {
	s := promise.
		NewPromise(func() (interface{}, error) {
			return 1, nil
		}).
		Then(func(i interface{}, e error) (interface{}, error) {
			return i.(int) * 2, nil
		}).
		Then(func(i interface{}, e error) (interface{}, error) {
			println(i.(int))
			return i.(int) * 2, nil
		}).
		Then(func(i interface{}, e error) (interface{}, error) {
			time.Sleep(10 * time.Second)
			return 10, nil
		}).
		OnError(func(e error) (interface{}, error) {
			println(e.Error())
			return nil, e
		}).
		OnSuccess(func(i interface{}) (interface{}, error) {
			println(i.(int))
			return i, nil
		}).
		StartTasks()
	println(time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(2 * time.Second)
	println(time.Now().Format("2006-01-02 15:04:05"))
	s.Interrupt()
	v, err := s.Await()
	println(time.Now().Format("2006-01-02 15:04:05"))
	println("final result:", v, err.Error())
}

package timed_task_tool

import (
	"time"
)

type TimedTask struct {
	ch chan bool
}

/*
*
上一次任务和下一次任务之间间隔delay秒
时间间隔 秒 ; 执行方法
*/
func DelayAtFixedRate(delay uint, fc func()) (tt *TimedTask) {
	tt = &TimedTask{make(chan bool)}
	go func() {
		for {
			select {
			case <-tt.ch:
				return
			default:
				fc()
			}
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}()
	return
}

/*
*
每隔delay秒发起一次任务
时间间隔 秒 ; 执行方法
*/
func DelayWithFixedDelay(delay uint, fc func()) (tt *TimedTask) {
	tt = &TimedTask{make(chan bool)}
	go func() {
		ticker := time.NewTicker(time.Duration(delay) * time.Second)
		for range ticker.C {
			select {
			case <-tt.ch:
				ticker.Stop()
				return
			default:
				fc()
			}
		}
	}()
	return
}

/*
*
任务终止，等待上一次任务执行完成后
*/
func (tt *TimedTask) Close() {
	tt.ch <- true
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func StartWorker(wg *sync.WaitGroup,id int, terminate <-chan int,jobs <-chan int){
	wg.Add(1)
	for {
		select {
		case <-terminate:
			fmt.Printf("Остановлен воркер №%v\n",id)
			wg.Done()
			return
		case <-jobs:
			fmt.Printf("Воркер №%v:%v\n",id,<-jobs)
			time.Sleep(time.Second/20)
		}
	}
}

func WriteChannelINFINITE(jobs chan<-int){
	for {
		jobs<-rand.Intn(100)
		time.Sleep(time.Duration(rand.Float64()/10)*time.Second)
	}
}

func main(){
	//Опять же, для серьезного проекта использовал бы контекст с таймаутом
	var wg sync.WaitGroup
	n:=10
	timer:=time.NewTimer(time.Duration(n)*time.Second)
	workersAmount:=2
	jobs:=make(chan int)
	terminateChanel:=make(chan int)
	go WriteChannelINFINITE(jobs)
	for i:=1;i<workersAmount+1;i++{
		go StartWorker(&wg,i,terminateChanel,jobs)
	}
	<-timer.C
	close(terminateChanel)
	wg.Wait()
	os.Exit(0)
}

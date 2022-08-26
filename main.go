package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func test() {
	fmt.Println("Test start!")
	defer fmt.Println("Test END!")
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Recovered in test:", e)
		}
	}()
	fmt.Println("Test before panic!")
	panic("__panic err: goliam problem__")
	fmt.Println("After panic!")
}

func main() {
	defer func() {
		fmt.Println("Main END!")
	}()
	go signalCheck()
	// time.Sleep(1 * time.Second)
	domemalloc()
	// test()
}

func signalCheck() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("******* awaiting signal")
	<-done
	fmt.Println("exiting")
}

func domemalloc() {
	defer func() {
		if err := recover(); err != nil {
			println("panic for: ", err)
		}
	}()
	slimslice := make([]int64, 1<<20)
	fatslice := make([]int64, 0, (1<<20)*10) // 50MB// /*1024*1024*1024/2*/
	start := time.Now()
	for {
		fatslice = append(fatslice, slimslice...)
		if len(fatslice)+len(slimslice) > cap(fatslice) {
			size := len(fatslice) * 8
			t := time.Since(start).Seconds()
			fmt.Printf("%d MB at %.2f = %.2f MB/s\n", size>>20, t, float64(size)/(1<<20)/t)
		}
	}
}

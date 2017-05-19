package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

type cond struct {
	condition *sync.Cond
}

type Set map[string]struct{}

func main() {

	s := make(Set)
	s["item1"] = struct{}{}
	s["item2"] = struct{}{}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": s,
		})
	})
	r.Run(":8083") // listen and serve on 0.0.0.0:8080

	//var dataStream chan interface{}
	//dataStream := make(chan interface{})
	//To declare a unidirectional channel, you’ll simply include the <- operator. \
	// To both declare and instantiate a channel which can only read,
	// place the <- operator on the left-hand side like so:

	//var dataStream <-chan interface{}
	//dataStream := make(<-chan interface{})

	// And to declare and create a channel which can only send you place the
	// <- operator on the right-hand side like so:

	//var dataStream chan<- interface{}
	//dataStream := make(chan<- interface{})

}

func runOne() {

	fmt.Println("runOne running")

	fmt.Println("runOne running triggered")
	//wg.Done()
}

func runTwo() {

	fmt.Println("runTwo running")

	fmt.Println("runTwo running triggered")

}

func runner() {
	fmt.Println("Heros of the telemark")

	cond1 := cond{condition: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var gor sync.WaitGroup
		gor.Add(1)
		go func() {
			gor.Done()
			c.L.Lock()
			defer c.L.Unlock()

			c.Wait()
			fn()
		}()
		gor.Wait()
	}

	fmt.Println("Adding to queue")
	grw := &sync.WaitGroup{}
	grw.Add(2)
	subscribe(cond1.condition, func() {

		runOne()
		grw.Done()
	})
	subscribe(cond1.condition, func() {

		runTwo()
		grw.Done()
	})
	cond1.condition.Broadcast()
	grw.Wait()
	//c.L.Unlock()

	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()
	fmt.Println("plp")
	instance := myPool.Get()
	myPool.Put(instance)
	fmt.Println("lll")

	myPool.Get()

	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		}}

	// Seed the pool with 4KB
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

			// Assume something interesting, but quick is being done with
			// this memory.
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

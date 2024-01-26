package main

import (
	"sync"

	"github.com/rezaAmiri123/sternx/philosopher"
)

const (
	philoCount = 9
)

func main() {
	agent := philosopher.NewAgent(philoCount)
	var wg sync.WaitGroup
	wg.Add(philoCount)

	for i := 0; i < philoCount; i++ {
		go func(philoNum int) {
			philo := philosopher.NewPhilosopher(
				philoNum,
				agent.GetLeftChopstick(philoNum),
				agent.GetRightChopstick(philoNum+1),
				agent.GetRequestChan(),
				agent.GetResponseChan(philoNum),
				agent.GetDoneChan(),
			)
			philo.Run()
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		agent.Terminate()
	}()
	agent.Run()
	agent.PrintStatic()
}

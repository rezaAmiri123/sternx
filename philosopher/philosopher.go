package philosopher

import (
	"fmt"
	"sync"
	"time"
)

const (
	defaultEatTimes       = 30
	defaultNumEatingPhilo = 2
)

type (
	Chopstick struct {
		idx int // Including the index can make debugging simpler
		mu  sync.Mutex
	}
	EatRequest struct {
		who int // Who is making the request
	}
	EatResponse struct {
		who      int
		duration time.Duration
	}
	EatDone struct {
		who int
	}
	Philosopher struct {
		philoNum        int
		eatTimes        int
		leftCS, rightCS *Chopstick
		requestChan     chan EatRequest
		responseChan    chan EatResponse
		doneChan        chan EatDone
	}
)

func NewPhilosopher(
	philoNum int,
	leftCS, rightCS *Chopstick,
	requestChan chan EatRequest,
	responseChan chan EatResponse,
	doneChan chan EatDone,
	options ...PhilosopherOption,
) *Philosopher {
	p := &Philosopher{
		philoNum:     philoNum,
		eatTimes:     defaultEatTimes,
		leftCS:       leftCS,
		rightCS:      rightCS,
		requestChan:  requestChan,
		responseChan: responseChan,
		doneChan:     doneChan,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *Philosopher) Run() {
	for numEat := 0; numEat < p.eatTimes; numEat++ {
		// once the philosopher intends to eat, lock the corresponding chopsticks
		for {
			p.leftCS.mu.Lock()
			// Attempt to get the right Chopstick -
			// if someone else has it we replace the left chopstick and try again
			// (in order to avoid deadlocks)
			if p.rightCS.mu.TryLock() {
				break
			}
			p.leftCS.mu.Unlock()
		}

		// We have the chopsticks but need the hosts permission
		p.requestChan <- EatRequest{
			who: p.philoNum,
		}
		response := <-p.responseChan

		duration := response.duration / time.Millisecond
		fmt.Printf("philosopher %d starting to eat (duration %d) (%d feed)\n", p.philoNum, duration, numEat)
		time.Sleep(response.duration)
		fmt.Printf("philosopher %d finished eating (duration %d) (%d feed)\n", p.philoNum, duration, numEat)

		p.rightCS.mu.Unlock()
		p.leftCS.mu.Unlock()
		// Tell host that we have finished eating
		p.doneChan <- EatDone{who: p.philoNum}
	}
	fmt.Printf("philosopher %d is full\n", p.philoNum)

}

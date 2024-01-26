package philosopher

import (
	"fmt"
	"time"
)

type (
	philoData struct {
		idx          int
		feedCount    int
		feedDuration time.Duration
		responseChan chan EatResponse
	}

	Agent struct {
		philoCount     int
		numEatingPhilo int
		chopsticks     []*Chopstick
		requestChan    chan EatRequest
		doneChan       chan EatDone
		philos         []*philoData
	}
)

func NewAgent(philoCount int, options ...AgentOptions) *Agent {
	chopsticks := make([]*Chopstick, philoCount)
	philoDatas := make([]*philoData, philoCount)
	for i := 0; i < philoCount; i++ {
		philoDatas[i] = &philoData{
			idx:          i,
			responseChan: make(chan EatResponse),
		}
		chopsticks[i] = &Chopstick{idx: i}
	}
	a := &Agent{
		philoCount:     philoCount,
		numEatingPhilo: defaultNumEatingPhilo,
		chopsticks:     chopsticks,
		requestChan:    make(chan EatRequest),
		doneChan:       make(chan EatDone),
		philos:         philoDatas,
	}

	for _, option := range options {
		option(a)
	}

	return a
}

func (a *Agent) Run() {

	awaitRequest := a.requestChan

	var whoEating []int // tracks who is currently eating

	for {
		select {
		case request, ok := <-awaitRequest:
			if !ok {
				// Closed channel means that we are done (finishedChan is guaranteed to be empty)
				return
			}

			// Sanity check - confirm that philosopher is not being greedy! (should never happen)
			if index(whoEating, request.who) != -1 {
				panic("Multiple requests from same philosopher")
			}
			whoEating = append(whoEating, request.who)
			fmt.Printf("%d started eating (currently eating %v)\n", request.who, whoEating)
			var duration time.Duration
			average := a.getAverageTime()

			switch {
			case average > a.philos[request.who].feedDuration:
				duration = time.Millisecond
			default:
				duration = 2 * time.Millisecond
			}

			a.philos[request.who].responseChan <- EatResponse{
				who:      request.who,
				duration: duration,
			}

			a.philos[request.who].feedDuration += duration
			a.philos[request.who].feedCount++

		case finished := <-a.doneChan:
			idx := index(whoEating, finished.who)
			if idx != -1{
				whoEating = append(whoEating[:idx], whoEating[idx+1:]...)
			}
			
			fmt.Printf("%d completed eating (total duration %d) (curretly eating %v)\n", finished.who, a.philos[finished.who].feedDuration/time.Millisecond, whoEating)
		}
		// There has been a change in the number of philosopher's eating
		if len(whoEating) < a.numEatingPhilo {
			awaitRequest = a.requestChan
		} else {
			// Ignore new eat requests until a philosopher finishes
			// (nil channel will never be selected)
			awaitRequest = nil
		}
	}

}

func (a Agent) GetResponseChan(numPhilo int) chan EatResponse {
	return a.philos[numPhilo].responseChan
}

func (a Agent) GetRequestChan() chan EatRequest {
	return a.requestChan
}

func (a Agent) GetDoneChan() chan EatDone {
	return a.doneChan
}

func (a Agent) GetLeftChopstick(philoNum int) *Chopstick {
	return a.chopsticks[philoNum]
}

func (a Agent) GetRightChopstick(philoNum int) *Chopstick {
	return a.chopsticks[philoNum%a.philoCount]
}

func (a *Agent) Terminate() {
	close(a.requestChan)
	close(a.doneChan)
	for _, philo := range a.philos{
		close(philo.responseChan)
	}
}

func (a Agent) PrintStatic() {
	for _, philo := range a.philos {
		fmt.Printf("philo number %d, total feed %d, total duration %d\n", philo.idx, philo.feedCount, philo.feedDuration)
	}
}
func (a Agent) getAverageTime() time.Duration {
	var sum time.Duration
	for _, philo := range a.philos {
		sum += philo.feedDuration
	}
	if len(a.philos) == 0 {
		return 1
	}
	return sum / time.Duration(len(a.philos))
}

func index(s []int, v int) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

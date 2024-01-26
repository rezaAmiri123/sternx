package philosopher

import (
	"testing"
	"time"
)

func TestPhilosopher_Run(t *testing.T) {
	type fields struct {
		philoNum        int
		eatTimes        int
		leftCS, rightCS *Chopstick
		requestChan     chan EatRequest
		responseChan    chan EatResponse
		doneChan        chan EatDone
	}
	type wants struct {
		eatRequestOK  bool
		doneRequestOK bool
	}
	tests := map[string]struct {
		fields fields
		wants  wants
		on     func(p *Philosopher)
	}{
		"OK": {
			fields: fields{
				philoNum:     1,
				eatTimes:     1,
				leftCS:       &Chopstick{},
				rightCS:      &Chopstick{},
				requestChan:  make(chan EatRequest),
				responseChan: make(chan EatResponse),
				doneChan:     make(chan EatDone),
			},
			wants: wants{
				eatRequestOK:  true,
				doneRequestOK: true,
			},
		},
		"Close Done Request Channel": {
			fields: fields{
				philoNum:     1,
				eatTimes:     1,
				leftCS:       &Chopstick{},
				rightCS:      &Chopstick{},
				requestChan:  make(chan EatRequest),
				responseChan: make(chan EatResponse),
				doneChan:     make(chan EatDone),
			},
			on: func(p *Philosopher) {
				close(p.doneChan)
			},
			wants: wants{
				eatRequestOK:  true,
				doneRequestOK: false,
			},
		},
		"Close Eat Request Channel": {
			fields: fields{
				philoNum:     1,
				eatTimes:     1,
				leftCS:       &Chopstick{},
				rightCS:      &Chopstick{},
				requestChan:  make(chan EatRequest),
				responseChan: make(chan EatResponse),
				doneChan:     make(chan EatDone),
			},
			on: func(p *Philosopher) {
				close(p.requestChan)
			},
			wants: wants{
				eatRequestOK:  false,
				doneRequestOK: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			philo := &Philosopher{
				philoNum:     tt.fields.philoNum,
				eatTimes:     tt.fields.eatTimes,
				leftCS:       tt.fields.leftCS,
				rightCS:      tt.fields.rightCS,
				requestChan:  tt.fields.requestChan,
				responseChan: tt.fields.responseChan,
				doneChan:     tt.fields.doneChan,
			}
			go func() {
				philo.Run()
			}()
			
			if tt.on != nil {
				tt.on(philo)
			}
			eatRequest, ok := <-philo.requestChan
			if ok != tt.wants.eatRequestOK {
				t.Errorf("eat request channel ok is %v, wanted ok is %v", ok, tt.wants.eatRequestOK)
				return
			}
			if ok {
				assertWrongPhilosopherIndex(t, eatRequest.who, philo.philoNum)
			} else {
				return
			}
			assertChopstickIsNotLock(t, philo.leftCS)
			assertChopstickIsNotLock(t, philo.rightCS)

			philo.responseChan <- EatResponse{
				who:      eatRequest.who,
				duration: time.Microsecond,
			}

			doneRequest, ok := <-philo.doneChan
			if ok != tt.wants.doneRequestOK {
				t.Errorf("done request channel ok is = %v, wanted ok is %v", ok, tt.wants.doneRequestOK)
				return
			}
			if ok {
				assertWrongPhilosopherIndex(t, doneRequest.who, philo.philoNum)
			}

		})
	}
}

func assertWrongPhilosopherIndex(t *testing.T, gotIndex, wantIndex int) {
	t.Helper()
	if gotIndex != wantIndex {
		t.Errorf("Philosopher is not the same got = %v, wanted = %v", gotIndex, wantIndex)
	}
}

		
func assertChopstickIsNotLock(t *testing.T,cs *Chopstick){
	t.Helper()
	if cs.mu.TryLock(){
		t.Errorf("Chopstick is not Lock")
		cs.mu.Unlock()
	}
}

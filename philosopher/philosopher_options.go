package philosopher

type PhilosopherOption func(*Philosopher)

func WithEatTimes(eatTimes int)PhilosopherOption{
	return func(p *Philosopher) {
		p.eatTimes = eatTimes
	}
}

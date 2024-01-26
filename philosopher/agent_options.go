package philosopher

type AgentOptions func(*Agent)

func WithNumEatingPhilo(numEatingPhilo int) AgentOptions{
	return func(a *Agent) {
		a.numEatingPhilo = numEatingPhilo
	}
}
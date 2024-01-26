# sternx
There are two major struct. 
agent and philosopher. 
philosophers are connected to the agent through three channels
philosophers are not connected to each other.

I want to summarize the whole cycle. 

First, we define an agent that is going to handle our philosophers
this agent defines Chopsticks and channels. 

Second, we define some philosophers in different goroutines,
and we pass to each of them two Chopsticks and three channels.

Third, we run the agent and philosophers.

Fourth, every philosopher which is being run in a different goroutine tries to lock one Chopstick. If it was successful, It would try to lock the other Chopstick. If it was successful, It would send an eating request to the agent through the eating request channel. after that, it waits for the response.  meanwhile, the agent which is listening to the eating request channel receives the request. the agent makes some calculations and sends a time duration back to the philosopher through the eating response channel. next,  the philosopher sleeps for the time duration that was received. (It simulates the eating). after that, the philosopher unlocks the Chopsticks and sends a done signal through the done eating channel to the agent. next, it checks a condition which if is true, the philosopher tries to lock Chopsticks and does the whole cycle again. meanwhile, the agent receives the done eating signal. next, it updates its statistical data, and it listens to the eating request channel and does the whole cycle again.

Note:
whenever a philosopher tries to lock both Chopsticks If it is not successful, It will unlock both Chopsticks. this mechanism helps to prevent deadlock.

Note:
The agent has statistical data that helps to manage philosophers
the agent uses this data to send a better time duration to philosophers.
wherever a philosopher uses fewer resources than the average, the agent sends more time duration.
this mechanism helps to prevent starvation.

Note:
there is a unit test for the run method at Philosopher struct

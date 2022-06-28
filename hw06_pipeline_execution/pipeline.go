package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for i := 0; i < len(stages); i++ {
		in = stages[i](proxyInput(in, done))
	}

	return in
}

func proxyInput(in, done In) Bi {
	out := make(Bi)

	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
	}()

	return out
}

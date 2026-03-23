package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = wrap(stage, done)(in)
	}
	return in
}

func wrap(stage Stage, done In) Stage {
	return func(in In) Out {
		proxyIn := stageWithDone(in, done)
		stageOut := stage(proxyIn)
		proxyOut := stageWithDone(stageOut, done)

		return proxyOut
	}
}

func stageWithDone(in In, done In) Out {
	proxy := make(Bi)
	go func() {
		defer func() {
			close(proxy)
			drain(in)
		}()
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case proxy <- val:
				}
			}
		}
	}()

	return proxy
}

func drain(in In) {
	for range in { //nolint:revive
	}
}

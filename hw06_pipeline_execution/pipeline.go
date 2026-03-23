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

//nolint:gocognit
func wrap(stage Stage, done In) Stage {
	return func(in In) Out {
		proxyIn := make(Bi)
		proxyOut := make(Bi)
		stageOut := stage(proxyIn)

		// inner chan consumer
		go func() {
			defer func() {
				close(proxyIn)
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
					case proxyIn <- val:
					}
				}
			}
		}()

		// outer chan producer
		go func() {
			defer func() {
				close(proxyOut)
				drain(stageOut)
			}()
			for {
				select {
				case <-done:
					return
				case val, ok := <-stageOut:
					if !ok {
						return
					}
					select {
					case <-done:
						return
					case proxyOut <- val:
					}
				}
			}
		}()

		return proxyOut
	}
}

func drain(in In) {
	for range in { //nolint:revive
	}
}

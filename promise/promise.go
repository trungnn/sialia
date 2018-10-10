package promise

type Promise struct {
	doneC chan struct{}
	startC chan struct{}

	AutoStart bool
	IsSettled bool
	Res       interface{}
	Err       error
}

func newPromise() *Promise {
	return &Promise{
		doneC: make(chan struct{}),
		startC: make(chan struct{}),
	}
}

func New(fn func()(interface{}, error)) *Promise {
	p := newPromise()

	go func() {
		if !p.AutoStart {
			<-p.startC
		}

		p.settle(fn())
	}()

	return p
}

func (p *Promise) settle(res interface{}, err error) {
	if p.IsSettled {
		return // redundant
	}

	p.Res, p.Err = res, err
	p.IsSettled = true
	p.doneC <- struct{}{}

	close(p.doneC)
	close(p.startC)
}

func (p *Promise) Then(fn func(interface{})(interface{}, error)) *Promise {
	return New(func()(interface{}, error) {
		res, err := p.Await()

		if err != nil {
			return nil, err
		}

		return fn(res)
	})
}

func (p *Promise) Await() (interface{}, error) {
	if !p.AutoStart {
		p.startC <- struct{}{}
	}

	<- p.doneC

	return p.Res, p.Err
}

package promise

type Promise struct {
	doneC chan struct{}

	Res interface{}
	Err error
}

func New(fn func()(interface{}, error)) *Promise {
	p := &Promise{
		doneC: make(chan struct{}),
	}

	go func() { p.fulfill(fn()) }()

	return p
}

func (p *Promise) fulfill(res interface{}, err error) {
	p.Res, p.Err = res, err
	p.doneC <- struct{}{}
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
	<- p.doneC
	return p.Res, p.Err
}

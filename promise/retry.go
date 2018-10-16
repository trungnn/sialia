package promise

type RetryCheckFn func(error) bool

type AbortErr struct {
	err error
}

func (aerr *AbortErr) Error() string {
	return aerr.err.Error()
}

type RetryOpts struct {
	RetryCheck RetryCheckFn
	MaxTries   int
}

func NewWithRetry(fn PromisifiableFn, opts *RetryOpts) *Promise {
	p := New(fn)

	for tries := 1; tries < opts.MaxTries; tries++ {
		p = p.Catch(func(err error) (interface{}, error) {
			if _, ok := err.(*AbortErr); ok {
				return nil, err
			}

			if opts.RetryCheck != nil && !opts.RetryCheck(err) {
				return nil, &AbortErr{err: err}
			}

			return fn()
		})
	}

	return p.Catch(func(err error) (interface{}, error) {
		if aerr, ok := err.(*AbortErr); ok {
			return nil, aerr.err
		}

		return nil, err
	})
}

package promise

type ResList []interface{}
type PromiseList []*Promise

type PromiseAllOpts struct {
	MaxConcurrency int
	WaitAllSettled bool
	Promises       []*Promise
}

func All(opts PromiseAllOpts) *Promise {
	if len(opts.Promises) > 0 {
		allP := newPromise()
		allP.AutoStart = true
		allP.all(opts.Promises, opts.WaitAllSettled, opts.MaxConcurrency)

		return allP
	} else {
		return New(func() (interface{}, error) {
			var res interface{}
			if opts.WaitAllSettled {
				res = make(PromiseList, 0)
			} else {
				res = make(ResList, 0)
			}

			return res, nil
		})
	}
}

func (p *Promise) all(childPromises []*Promise, waitAllSettled bool, maxConcurrency int) {
	var workC chan struct{}
	if maxConcurrency > 0 {
		workC = make(chan struct{}, maxConcurrency)
	}

	var res interface{}
	if waitAllSettled {
		res = make(PromiseList, len(childPromises))
	} else {
		res = make(ResList, len(childPromises))
	}

	remaining := len(childPromises)

	updateProgress := func(i int, v interface{}) {
		if workC != nil {
			workC <- struct{}{}
		}

		if p.IsSettled {
			return
		}

		switch val := v.(type) {
		case *Promise:
			pList := res.(PromiseList)
			pList[i] = val
		default:
			rList := res.(ResList)
			rList[i] = val
		}

		if remaining--; remaining == 0 {
			if workC != nil {
				close(workC)
			}

			p.settle(res, nil)
		}
	}

	for index, cp := range childPromises {
		go func(i int, promise *Promise) {
			if workC != nil {
				<-workC
			}

			res, err := promise.Await()

			if waitAllSettled {
				updateProgress(i, promise)
			} else {
				if err != nil {
					p.settle(nil, err)
				} else {
					updateProgress(i, res)
				}
			}
		}(index, cp)

		if workC != nil && index < maxConcurrency {
			workC <- struct{}{}
		}
	}
}

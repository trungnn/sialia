package promise

type ResList []interface{}
type PromiseList []*Promise

type PromiseAllOpts struct {
	MaxConcurrency int
	WaitAllSettled bool
	Promises []*Promise
}

func All(opts PromiseAllOpts) *Promise {
	allP := newPromise()
	allP.AutoStart = true

	allP.all(opts.Promises, opts.WaitAllSettled, opts.MaxConcurrency)

	return allP
}

func (p *Promise) all(childPromises []*Promise, waitAllSettled bool, maxConcurrency int) {
	workC := make(chan struct{}, maxConcurrency)

	var res interface{}
	if waitAllSettled {
		res = make(PromiseList, len(childPromises))
	} else {
		res = make(ResList, len(childPromises))
	}

	remaining := len(childPromises)

	updateProgress := func(i int, v interface{}) {
		workC <- struct{}{}

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
			close(workC)
			p.settle(res, nil)
		}
	}

	for index, cp := range childPromises {
		go func(i int, promise *Promise) {
			<-workC

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

		if index < maxConcurrency {
			workC <- struct{}{}
		}
	}
}


package promise

func All(promises ...*Promise) *Promise {
	allP := &Promise{
		doneC: make(chan struct{}),
	}

	res := make([]interface{}, len(promises))
	remaining := len(promises)

	updateProgress := func(i int, v interface{}) {
		if allP.IsSettled {
			return
		}

		res[i] = v

		if remaining--; remaining == 0 {
			allP.settle(res, nil)
		}
	}

	for index, promise := range promises {
		go func(i int, p *Promise) {
			res, err := p.Await()

			if err != nil {
				allP.settle(nil, err)
			} else {
				updateProgress(i, res)
			}
		}(index, promise)
	}

	return allP
}

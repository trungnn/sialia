package promise

type PromiseMapFn func(interface{}) *Promise

type PromiseMapOpts struct {
	Items []interface{}
	MapFn PromiseMapFn

	MaxConcurrency int
	WaitAllSettled bool
}

func Map(opts PromiseMapOpts) *Promise {
	var ps []*Promise

	for _, item := range opts.Items {
		ps = append(ps, opts.MapFn(item))
	}

	return All(PromiseAllOpts{
		Promises: ps,
		WaitAllSettled: opts.WaitAllSettled,
		MaxConcurrency: opts.MaxConcurrency,
	})
}

Sialia
===

Sialia is another way to refer to bluebirds. The name is inspired by the [Bluebird project](https://github.com/petkaantonov/bluebird) for Javascript.

The eventual goal is to achieve a somewhat similar functionality in Go.

Usage
---

Given the following function:

```golang
func Add(x, y int) (int, error) {
    time.Sleep(time.Duration(x + y) * time.Second)
    if (x + y) % 2 == 0 {
        return x + y, nil
    } else {
        return nil, errors.New("error")
    }
}
```

Create a new Promise

```golang
p := promise.New(func()(interface{}, error) {   
    return Add(1, 1)
})

// the following line will block until the promise is fulfilled
res, err := p.Await() // => 2, nil
```

Promise chaining

```golang
p := promise.New(func()(interface{}, error) {   
    return Add(1, 1)
}).Then(func(i interface{})(interface{}, error) {
    return Add(i.(int), 2)
}).Then(func(i interface{})(interface{}, error) {
    return Add(i.(int), 4)
})

// the following line will block until the promise is fulfilled
res, err := p.Await() // => 8, nil

p2 := promise.New(func()(interface{}, error) {   
    return Add(1, 2)
}).Then(func(i interface{})(interface{}, error) {
    return Add(i.(int), 2)
}).Then(func(i interface{})(interface{}, error) {
    return Add(i.(int), 4)
})

// the following line will block until the promise is fulfilled
res, err := p2.Await() // => nil, error
```

Waiting for more than one promise to complete

```golang
p1 := promise.New(func()(interface{}, error) {   
    return Add(1, 1)
})
p2 := promise.New(func()(interface{}, error) {   
    return Add(2, 2)
})
p3 := promise.New(func()(interface{}, error) {   
    return Add(3, 3)
})

p := promise.All(p1, p2, p3)
// the following line will block until p1, p2, p3 are fulfilled
res, err := p.Await() // => [2, 4, 6], nil


p1 := promise.New(func()(interface{}, error) {   
    return Add(1000, 1000)
})
p2 := promise.New(func()(interface{}, error) {   
    return Add(1, 2)
})

p := promise.All(p1, p2)
// the following line will block until p2 failed
res, err := p.Await() // => nil, error
```
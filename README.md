# syncslice

Goroutine-safe wrapper interface over a generic slice using [sync.RWMutex](https://pkg.go.dev/sync#RWMutex).

You don't need this library. Slices are already safe. See https://github.com/nbd-wtf/syncslice/issues/1 and https://go.dev/blog/slices-intro.

```go
s := syncslice.Make[string](2, 3)

s.Set(0, "hello")
s.Set(1, "world")
s.Append("!")

for item := range s.Iter() {
    fmt.Println(item.Value)
}
```

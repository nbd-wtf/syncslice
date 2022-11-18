# syncslice

Goroutine-safe wrapper interface over a generic slice using [sync.RWMutex](https://pkg.go.dev/sync#RWMutex).

```go
s := syncslice.Make[string](2, 3)

s.Set(0, "hello")
s.Set(1, "world")
s.Append("!")

for item := range s.Iter() {
    fmt.Println(item.Value)
}
```

# Library of Slice-like Collections for Go

A library of data structures that all implement a Go slice-like interface, so 
that they act like an `[]interface{}` type. For example, the following code:

```go
s := slice.EmptySlice(0, 0)
s = s.Append(10, 13, 62, 124)
s = s.Slice(1, 2)
fmt.Printf("%d", s.Index(0).(int) + s.Index(1).(int))
```

Acts in the same way as:
```go
s := make([]interface{}, 0, 0)
s = append(s, 10, 13, 62, 124)
s = s[1:2]
fmt.Printf("%d", s[0].(int) + s[1].(int))
```

All functions in the `Slice` interface have a Go equivalent, except the 
`Prepend` function, which exists for simplicity and optimisation reasons

Technical documentation is [here](https://pkg.go.dev/github.com/bhollier/slice)

### Currently, the following data structures are supported:

- `Wrapper` (Simple wrapper around `[]interface{}`)
- Linked List (`Singly` and `Doubly`)
- `Distributed` Slice

The `Distributed` type is a custom data structure, that stores its elements in 
"buckets". This means that elements can be added onto the start or end of the
slice without reallocating the array. This also increases the speed of prepend 
operations, as prepending works as 'backwards appending', rather than
reallocating and copying

## Benchmarks:

| Data Structure    | Append (ns/op) | Prepend (ns/op) | Erase (ns/op) | Index (ns/op) | Iter (ns/op) |
| ----------------- | -------------- | --------------- | ------------- | ------------- | ------------ |
| **[]int**         | 1035           | 460465          | 505           | 316           | 188          |
| **[]interface{}** | 3225           | 2534283         | 815           | 340           | 290          |
| **Wrapper**       | 5483           | 3737938         | 937           | 1692          | 2552         |
| **Distributed**   | 3079           | 480965          | 45005         | 2534          | 3422         |
| **Doubly**        | 9249           | 9819            | 7220          | 2762764       | 14149        |
| **Singly**        | 7966           | 8056            | 16760         | 5137144       | 24290        |

For more info, see `baseline_test.go` and `slice_test.go`

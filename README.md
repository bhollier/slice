# Library of Slice-like Collections for Go

A library of data structures that all implement a Go slice-like interface, so 
that they act like an `[]interface{}` type. For example, the following code:

```go
s := slice.EmptySlice(0, 0)
s = s.Append(10, 13, 62, 124)
s = s.Slice(1, 2)
fmt.Printf("%d", s.Index(0) + s.Index(1))
```

Acts in the same way as:
```go
s := make([]int, 0, 0)
s = append(s, 10, 13, 62, 124)
s = s[1:2]
fmt.Printf("%d", s[0] + s[1])
```

All functions in the `Slice` interface have a Go equivalent, except the 
`Prepend` function, which exists for simplicity and optimisation reasons

Technical documentation is [here](https://pkg.go.dev/github.com/bhollier/slice)

### Currently, the following data structures are supported:

- `Wrapper` (Simple wrapper around `[]T`)
- Linked List (`Singly` and `Doubly`)
- `Distributed` Slice

The `Distributed` type is a custom data structure, that stores its elements in 
"buckets". This means that elements can be added onto the start or end of the
slice without reallocating the array

## Benchmarks:

| Data Structure    | Append (ns/op) | Prepend (ns/op) | Erase (ns/op)* | Index (ns/op) | Iter (ns/op) |
|-------------------|----------------|-----------------|----------------|---------------|--------------|
| **[]int**         | 588            | 644303          | 445            | 179           | 175          |
| **[]interface{}** | 4588           | 2746719         | 837            | 330           | 294          |
| **Wrapper**       | 674            | 295611          | 15601          | 1134          | 1732         |
| **Distributed**   | 1135           | 537222          | 62131          | 2504          | 2569         |
| **Doubly**        | 3957           | 4137            | 187978         | 2440531       | 12598        |
| **Singly**        | 3036           | 3194            | 139828         | 2638009       | 13186        |

&ast; Currently, the `Erase` method clones the `Slice` unnecessarily, so does not perform optimally

For more info, see `baseline_test.go` and `slice_test.go`

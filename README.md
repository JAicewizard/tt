# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 

# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}. everything marked with a S at the end benchmarks data containing a slice instead of a nested Data value

```
BenchmarkGobData-8               1000000              2846 ns/op             368 B/op         11 allocs/op
BenchmarkGobDataEncode-8         2000000              1527 ns/op             120 B/op          3 allocs/op
BenchmarkGobMap-8                 500000              4456 ns/op             494 B/op         13 allocs/op
BenchmarkGobMapEncode-8          1000000              3850 ns/op             462 B/op         11 allocs/op
BenchmarkGobDataS-8              1000000              2806 ns/op             368 B/op         11 allocs/op
BenchmarkGobDataEncodeS-8        2000000              1399 ns/op             120 B/op          3 allocs/op
BenchmarkGobMapS-8               1000000              4299 ns/op             510 B/op         13 allocs/op
BenchmarkGobMapEncodeS-8         1000000              3832 ns/op             478 B/op         11 allocs/op
```
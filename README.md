# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 

# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}. everything marked with a S at the end benchmarks data containing a slice instead of a nested Data value

i7 975k@3.6 ddr3@1596
```
BenchmarkGobData-8               5000000              2703 ns/op             368 B/op         11 allocs/op
BenchmarkGobDataEncode-8        10000000              1403 ns/op             120 B/op          3 allocs/op
BenchmarkGobMap-8                2000000              9011 ns/op             844 B/op         26 allocs/op
BenchmarkGobMapEncode-8          3000000              5008 ns/op             368 B/op         10 allocs/op
BenchmarkGobDataS-8              5000000              2699 ns/op             264 B/op         12 allocs/op
BenchmarkGobDataEncodeS-8       10000000              1388 ns/op             120 B/op          3 allocs/op
BenchmarkGobMapS-8               2000000              7297 ns/op             450 B/op         15 allocs/op
BenchmarkGobMapEncodeS-8         3000000              4645 ns/op             208 B/op          6 allocs/op
```
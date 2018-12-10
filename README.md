# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 

# limitations
tt only supports up to 255 child nodes per node
tt only supports up to 2^32 number of nodes(each submap, array, and value is a node)
text is up to 2^32 bytes long, same for the key

# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}. everything marked with a S at the end benchmarks data containing a slice instead of a nested Data value


i7 975k@3.6 ddr3@1596 make bench
```
BenchmarkGobData-8               5000000              3044 ns/op             536 B/op         10 allocs/op
BenchmarkGobDataEncode-8        10000000              1733 ns/op             288 B/op          4 allocs/op
BenchmarkGobMap-8                2000000              9232 ns/op             844 B/op         26 allocs/op
BenchmarkGobMapEncode-8          3000000              5119 ns/op             368 B/op         10 allocs/op
BenchmarkGobDataS-8              5000000              2896 ns/op             440 B/op         11 allocs/op
BenchmarkGobDataEncodeS-8       10000000              1589 ns/op             288 B/op          4 allocs/op
BenchmarkGobMapS-8               2000000              7286 ns/op             450 B/op         15 allocs/op
BenchmarkGobMapEncodeS-8         3000000              4609 ns/op             208 B/op          6 allocs/op
```

the increase in allocations in the latest update is caused by bytes.buffer `
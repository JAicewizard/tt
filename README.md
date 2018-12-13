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
BenchmarkGobData-8               3000000              5095 ns/op             986 B/op         23 allocs/op
BenchmarkGobDataEncode-8         3000000              4425 ns/op             300 B/op          5 allocs/op
BenchmarkGobMap-8                1000000             20113 ns/op            1656 B/op         56 allocs/op
BenchmarkGobMapEncode-8          1000000             12045 ns/op             432 B/op         12 allocs/op
```

this update consists of a change in benchmarks, dont compare these to the previous benchmarks.
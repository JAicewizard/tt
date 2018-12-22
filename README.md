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
BenchmarkGobData-8               2000000              7674 ns/op            1420 B/op         29 allocs/op
BenchmarkGobDataDecode-8         5000000              3551 ns/op             783 B/op         23 allocs/op
BenchmarkGobDataEncode-8         2000000              7031 ns/op             637 B/op          6 allocs/op
BenchmarkGobMap-8                 500000             34040 ns/op            2140 B/op         82 allocs/op
BenchmarkGobMapDecode-8          1000000             18148 ns/op            1484 B/op         62 allocs/op
BenchmarkGobMapEncode-8          1000000             21973 ns/op             656 B/op         20 allocs/op
```
this update consists of a change in benchmarks, dont compare these to the previous benchmarks.
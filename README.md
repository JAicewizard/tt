# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 

# Benchmarks

BenchmarkGobData-8               5000000              2822 ns/op             368 B/op         11 allocs/op
BenchmarkGobDataEncode-8        10000000              1460 ns/op             120 B/op          3 allocs/op
BenchmarkGobMap-8                3000000              4931 ns/op             586 B/op         15 allocs/op
BenchmarkGobMapEncode-8          3000000              3844 ns/op             462 B/op         11 allocs/op
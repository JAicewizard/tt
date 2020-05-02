# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 


# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}.

r5 3600 make bench
```
BenchmarkV3-12                   3577809              3260 ns/op             945 B/op         17 allocs/op
BenchmarkV3Decode-12             4935484              2343 ns/op             624 B/op         17 allocs/op
BenchmarkV3Encode-12            13271485               835 ns/op             343 B/op          0 allocs/op
BenchmarkGobData-12              4368571              2513 ns/op            1491 B/op         27 allocs/op
BenchmarkGobDataDecode-12        5306721              2271 ns/op             784 B/op         23 allocs/op
BenchmarkGobDataEncode-12        8861355              1235 ns/op             700 B/op          4 allocs/op
BenchmarkGobMap-12                869761             13557 ns/op            1916 B/op         68 allocs/op
BenchmarkGobMapDecode-12         1592424              7375 ns/op            1260 B/op         48 allocs/op
BenchmarkGobMapEncode-12         1000000             10094 ns/op             656 B/op         20 allocs/op
```

# v3 vs v2

As seen in the benchmarks v3 is significantly slower than v2, this has a couple of reasons:
- v3 supports all kinds of maps and structs meaning we have to use reflect for a lot of things. We have some fast-paths for commonly used maps but this doesnt negate all performance loss.
- v3 supports encoding/decoding from all kinds of readerswriters, the added indirection adds a significant amount of overhead.
- v3 uses the google varints, this means that we have removed any size limitations but it does mean slower encoding/decoding. see #limitations in ttv2

other diferences are:
- map[interaface{}]interface{} and map[string]interface{} get converted to Data in v2 but not in v3. In v3 all maps are converted into map[interface{}]interface{} if no other type is provided, much like json and gob.

# limitations in ttv2

ttv2 only supports up to 255 child nodes per node
ttv2 only supports up to 2^32 number of nodes(each submap, array, and value is a node)
text is up to 2^32 bytes long, same for the key

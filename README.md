# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 


# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}.

r5 3600 make bench
```
BenchmarkV3              3015724              3935 ns/op            1157 B/op         33 allocs/op
BenchmarkV3Decode        4057526              2961 ns/op             739 B/op         28 allocs/op
BenchmarkV3Encode       13256769               896 ns/op             383 B/op          5 allocs/op
BenchmarkV2              4915644              2423 ns/op            1389 B/op         27 allocs/op
BenchmarkV2Decode        9181489              1292 ns/op             784 B/op         23 allocs/op
BenchmarkV2Encode       10495838              1113 ns/op             592 B/op          4 allocs/op
BenchmarkGobMap           979080             11936 ns/op            1916 B/op         68 allocs/op
BenchmarkGobMapDecode    1848817              6436 ns/op            1260 B/op         48 allocs/op
BenchmarkGobMapEncode    2300364              5208 ns/op             656 B/op         20 allocs/op
```

## notes
- all tests are run using `GOMAXPROCS=1`, this is because on zen running on multiple threads will cause horrible cache-invalidation. A single alloc/op would cause the GC to run at some point, this would kick the benching to a diferent core. The reason I decided to run using `GOMAXPROCS=1` is because this doesnt have a big impact on Intel cpus, and any real world application would be generating garbage anyways, so eleminitin the GC from running should be part of the benchmark. Another reason coul be: real world applications would so something else in between runs causing cache-invalidation anyways.

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

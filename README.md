# tt

TinyTranscoder is a library that can encode/decode map[interface{}]interface{} where the interface is any of the supported types or a Transmitter.
tt is faster than GOB but does support less types and does use a significant more bytes to transmis the data 


# Benchmarks

Data is are benchmarks that use Data, map benchmarks just convert the data to a map[interface{}]interface{}.

r5 3600 make bench
```
BenchmarkV3               303894              3921 ns/op            1337 B/op         33 allocs/op
BenchmarkV3Decode         401946              2827 ns/op             778 B/op         28 allocs/op
BenchmarkV3Encode        1208140              1049 ns/op             558 B/op          5 allocs/op
BenchmarkV3int64         1657701               730 ns/op             572 B/op         11 allocs/op
BenchmarkV2               450891              2442 ns/op            1615 B/op         26 allocs/op
BenchmarkV2Decode         914095              1216 ns/op             776 B/op         22 allocs/op
BenchmarkV2Encode         962227              1174 ns/op             791 B/op          4 allocs/op
BenchmarkGobMap            99562             11806 ns/op            1916 B/op         68 allocs/op
BenchmarkGobMapDecode     183130              6511 ns/op            1260 B/op         48 allocs/op
BenchmarkGobMapEncode     235747              4983 ns/op             656 B/op         20 allocs/op
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

alloc:
	go test -bench=BenchmarkV3 -c
	env GODEBUG=allocfreetrace=1 ./tt.test -test.run=none -test.bench=BenchmarkV3$$ -test.benchtime=10ms 2>trace.log
	
bench:
	go test -bench=. -benchtime=10s -benchmem

shortbench:
	go test -bench=. -benchtime=1s -benchmem
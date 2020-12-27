alloc:
	go test -bench=BenchmarkV3 -c
	env GODEBUG=allocfreetrace=1 ./tt.test -test.run=none -test.bench=BenchmarkV3$$ -test.benchtime=10ms 2>trace.log
	
bench:
	env GOMAXPROCS=1 go test -bench=. -run=^\$ -benchtime=10s -benchmem

shortbench:
	env GOMAXPROCS=1 go test -bench=. -run=^\$ -benchtime=1s -benchmem

test:
	go test -short
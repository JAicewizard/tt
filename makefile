alloc:
	go test -bench=BenchmarkGobData -c
	env GODEBUG=allocfreetrace=1 ./tt.test -test.run=none -test.bench=BenchmarkGobDataDecode$$ -test.benchtime=10ms 2>trace.log
	
bench:
	go test -bench=BenchmarkGob -benchtime=10s -benchmem

shortbench:
	go test -bench=BenchmarkGob -benchtime=1s -benchmem
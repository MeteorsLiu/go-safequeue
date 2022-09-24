# go-safequeue
Based on CAS Lock-Free Algorithm Queue for Go


# Benchmark

```
BenchmarkGoMutexInsert-8    	 2280337	       582.4 ns/op	     101 B/op	       1 allocs/op

BenchmarkGoMutexRead-8      	 3821852	       314.8 ns/op	      32 B/op	       1 allocs/op

BenchmarkLockFreeInsert-8   	 3396618	       359.5 ns/op	      48 B/op	       2 allocs/op

BenchmarkLockFreeRead-8     	 3603032	       330.8 ns/op	      24 B/op	       1 allocs/op
```
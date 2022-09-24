# go-safequeue
Based on CAS Lock-Free Algorithm Queue for Go

# Feature

**Generic Support**

**Lightweight And Fast**



# Benchmark

```
BenchmarkGoMutexInsert-8    	 2606114	       576.9 ns/op	     103 B/op	       1 allocs/op
BenchmarkGoMutexRead-8      	 3668961	       319.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkLockFreeInsert-8   	 3465559	       379.5 ns/op	      48 B/op	       2 allocs/op
BenchmarkLockFreeRead-8     	 3687990	       336.2 ns/op	      24 B/op	       1 allocs/op
BenchmarkChannelInsert-8    	 3518539	       331.0 ns/op	      24 B/op	       1 allocs/op
BenchmarkChannelRead-8      	 3545941	       341.1 ns/op	      24 B/op	       1 allocs/op
```

# Conclusion

However, according to the benchmark result, Golang's buffered channel is faster than CAS Lock-free Queue.

I recommend you to use golang's buffered channel as the queue.
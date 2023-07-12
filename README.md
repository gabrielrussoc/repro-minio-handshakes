# Repro

AWS_ACCESS_KEY_ID="..." AWS_SECRET_ACCESS_KEY="..." go test -trace=trace.out

And then can look at trace with:
```
go tool trace -http localhost:9091 trace.out
```

and cpu profile with:
```
go tool pprof -http localhost:9092 /tmp/cpu-2532064764.pprof
```

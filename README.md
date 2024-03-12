# Steps

1. start the first async service:
```shell
go run delegate.go async_1.go
```
2. start other async services:
```shell
go run delegate.go async_2.go --join 192.168.2.100:7946
go run delegate.go async_3.go --join 192.168.2.100:7946
```

# How to apply to xCherry

Assume the application is going to start X async services with Y shards.

1. Start the first async service
2. Start the left (X - 1) async services on different ports and join the address of the first async service
3. In each service, pull data from shards assigned to it:
```shell
for key := 0; key < Y; key++ {
    node, ok := devt.consistent.GetNode(strconv.Itoa(key))
    if !ok {
        continue
    } 
    
    pull data with shard == key from the db and process.
}
```

# Questions

There might be times when re-assignments are underway, and multiple nodes might end up with the same shard. 
To prevent it, we might incorporate additional methods for double-checking.
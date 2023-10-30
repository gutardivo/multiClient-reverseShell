## Let's go:
```
$ git clone https://github.com/gutardivo/multiClient-reverseShell.git
$ cd multiClient-reverseShell
```
#### For linux target:
```
$ cd client/linux
$ GOOS=linux go build client.go
```

#### For windows target:
```
$ cd client/windows
$ GOOS=windows go build client.go
```

#### Run server:
```
$ cd server
$ go run server.go
```

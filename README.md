# Scoreboard for roller Derby

## Running

1. Open `dist` folder 
2. Launch  `server` application
3. Open your browser [http://localhost:3000/scoreboard](http://localhost:3000/scoreboard)

## Development

### Pre-requisite

* have a sane installation of `go`
* have a sane installation of `nodejs`

### Install tools

Install javascript tools:
```
npm install -g bower gulp
```

### Compile all

just type 
```
make
```

### Compile server

From the folder `src/server` launch
```
go build main.go
```

### Compile client

From the folder `src/server` launch
```
gulp
```

### Develop

Launch server
```
cd src/server
go run main.go
```

Launch client
```
cd src/client
gulp serve
```
and open your browser [http://localhost:3000/scoreboard/dev.html]
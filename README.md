[![Docker automated build](https://img.shields.io/badge/docker-automated--build-blue.svg?style=flat-square)](https://hub.docker.com/r/xjewer/go-test-teamcity/)

# Golang test TeamCity converter

Convert go test output to TeamCity format

Support Run, Skip, Pass, Fail

### Installation

```bash
go get github.com/2tvenom/go-test-teamcity
```

### How use
```bash
go test -v ./... | go-test-teamcity
```

### Docker
```bash
go test -v ./... | docker run -i xjewer/go-test-teamcity
```

### Docker multi-stage build
Extending Golang Dockerhub instructions to `Start a Go instance in your app`:
https://hub.docker.com/_/golang

> The most straightforward way to use this image is to use a Go container as both the build and runtime environment. In your Dockerfile, writing something along the lines of the following will compile and run your project:

```Dockerfle

...

COPY --from=xjewer/go-test-teamcity /converter /usr/local/bin/go-test-teamcity
RUN  go test -v ./... | go-test-teamcity
```

### Links
- https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity

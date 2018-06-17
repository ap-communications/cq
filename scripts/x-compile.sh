export GOOS=windows
export GOARCH=amd64
go build -o cq-$GOOS-$GOARCH.exe src/cq.go

export GOOS=windows
export GOARCH=386
go build -o cq-$GOOS-$GOARCH.exe src/cq.go

export GOOS=linux
export GOARCH=amd64
go build -o cq-$GOOS-$GOARCH src/cq.go

export GOOS=linux
export GOARCH=386
go build -o cq-$GOOS-$GOARCH src/cq.go

export GOOS=linux
export GOARCH=arm
go build -o cq-$GOOS-$GOARCH src/cq.go

export GOOS=darwin
export GOARCH=amd64
go build -o cq-$GOOS-$GOARCH src/cq.go

export GOOS=darwin
export GOARCH=386
go build -o cq-$GOOS-$GOARCH src/cq.go

unset GOOS
unset GOARCH

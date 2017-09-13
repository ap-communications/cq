export GOOS=windows
export GOARCH=amd64
go build -o cq.exe ../
zip cq-1.1-$GOOS-$GOARCH.zip cq.exe
rm cq.exe

export GOOS=windows
export GOARCH=386
go build -o cq.exe ../
zip cq-1.1-$GOOS-$GOARCH.zip cq.exe
rm cq.exe

export GOOS=linux
export GOARCH=amd64
go build ../
tar zcvf cq-1.1-$GOOS-$GOARCH.tar.gz cq
rm cq

export GOOS=linux
export GOARCH=386
go build ../
tar zcvf cq-1.1-$GOOS-$GOARCH.tar.gz cq
rm cq

export GOOS=linux
export GOARCH=arm
go build ../
tar zcvf cq-1.1-$GOOS-$GOARCH.tar.gz cq
rm cq

export GOOS=darwin
export GOARCH=amd64
go build ../
tar zcvf cq-1.1-$GOOS-$GOARCH.tar.gz cq
rm cq

export GOOS=darwin
export GOARCH=386
go build ../
tar zcvf cq-1.1-$GOOS-$GOARCH.tar.gz cq
rm cq

unset GOOS
unset GOARCH
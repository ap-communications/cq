set GOOS=windows
set GOARCH=amd64
go build -o %GOOS%_%GOARCH%\%~n1.exe %1

set GOOS=windows
set GOARCH=386
go build -o %GOOS%_%GOARCH%\%~n1.exe %1

set GOOS=linux
set GOARCH=amd64
go build -o %GOOS%_%GOARCH%\%~n1 %1

set GOOS=linux
set GOARCH=386
go build -o %GOOS%_%GOARCH%\%~n1 %1

set GOOS=linux
set GOARCH=arm
go build -o %GOOS%_%GOARCH%\%~n1 %1

set GOOS=darwin
set GOARCH=amd64
go build -o %GOOS%_%GOARCH%\%~n1 %1

set GOOS=darwin
set GOARCH=386
go build -o %GOOS%_%GOARCH%\%~n1 %1
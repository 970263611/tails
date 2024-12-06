@echo off
if exist tails (
  del tails
)
if exist tails.exe (
  del tails.exe
)
echo cpu架构列表:
echo  1、arm64
echo  2、amd64
set /p a=请选择架构:
if %a% == 1  (
   set a="arm64"
) else (
   set a="amd64"
)
echo 操作系统列表:
echo  1、linux
echo  2、windows
echo  3、darwin
set /p b=请选择系统:
if %b% == 1  (
    set b="linux"
) else (
    if %b% == 2 (
       set b="windows"
    ) else (
       set b="darwin"
    )
)
go env -w GOARCH=%a%
go env -w GOOS=%b%
if %b% == "windows" (
  go build  -o tails.exe ./boot/main.go
) else (
  go build  -o tails ./boot/main.go
)
echo 打包完成，cpu架构'%a%' 操作系统'%b%'
@echo off
echo cpu�ܹ��б�:
echo  1��arm64
echo  2��amd64
set /p a=��ѡ��ܹ�:
if %a% == 1  (
   set a="arm64"
) else (
   set a="amd64"
)
echo ����ϵͳ�б�:
echo  1��linux
echo  2��windows
echo  3��darwin
set /p b=��ѡ��ϵͳ:
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
go build ./boot/main.go
echo �����ɣ�cpu�ܹ�'%a%' ����ϵͳ'%b%'
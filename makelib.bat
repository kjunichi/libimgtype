go build -buildmode=c-archive -o libimgtype.a main.go
gcc -m64 -shared -o imgtype.dll imgtype.def libimgtype.a -Wl,--allow-multiple-definition -static -lwinmm -lntdll -lWs2_32 
lib /machine:x64 /def:imgtype.def

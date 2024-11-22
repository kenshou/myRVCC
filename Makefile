#编译本项目的可执行文件为rvcc

rvcc:
	go build  -o build/rvcc main.go

test: rvcc
	sh ./test.sh

clean:
	rm -rf build/rvcc build/tmp*
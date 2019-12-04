.PHONY=dist clean darwin linux windows
BUILDARGS=-ldflags '-w -s' -trimpath

all: dist darwin linux windows

clean:
	$(RM) -fvr dist

dist:
	mkdir dist

darwin:
	GOOS=darwin go build ${BUILDARGS} -o dist/ezpwd_darwin ./ezpwd/

linux:
	GOOS=linux go build ${BUILDARGS} -o dist/ezpwd_linux ./ezpwd/

windows:
	GOOS=windows go build ${BUILDARGS} -o dist/ezpwd_windows ./ezpwd/

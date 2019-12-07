.PHONY=dist clean darwin linux windows
BUILDARGS=-ldflags='-w -s' -trimpath

all: dist darwin linux windows

clean:
	$(RM) -fvr dist

dist:
	mkdir dist

darwin:
	CGO_ENABLED=0 GOOS=darwin go build ${BUILDARGS} -o dist/ezpwd_darwin ./ezpwd/

linux:
	CGO_ENABLED=0 GOOS=linux go build ${BUILDARGS} -o dist/ezpwd_linux ./ezpwd/

windows:
	CGO_ENABLED=0 GOOS=windows go build ${BUILDARGS} -o dist/ezpwd_windows ./ezpwd/

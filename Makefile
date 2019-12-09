.PHONY=dist clean darwin linux windows
BUILDARGS=-ldflags='-w -s' -trimpath

all: dist darwin linux windows

clean:
	$(RM) -fvr dist

dist:
	mkdir dist

darwin:
	CGO_ENABLED=0 GOOS=darwin go build ${BUILDARGS} -o dist/ezpwd_darwin ./ezpwd/
	CGO_ENABLED=0 GOOS=darwin go build ${BUILDARGS} -o dist/ezpwd_gui_darwin ./ezpwd_gui/

linux:
	CGO_ENABLED=0 GOOS=linux go build ${BUILDARGS} -o dist/ezpwd_linux ./ezpwd/
	CGO_ENABLED=0 GOOS=linux go build ${BUILDARGS} -o dist/ezpwd_gui_linux ./ezpwd_gui/

windows:
	CGO_ENABLED=0 GOOS=windows go build ${BUILDARGS} -o dist/ezpwd_windows ./ezpwd/
	CGO_ENABLED=0 GOOS=windows go build ${BUILDARGS} -o dist/ezpwd_gui_windows ./ezpwd_gui/

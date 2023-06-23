NAME=hozonsite
VERSION := $(shell cat main.go | grep "var version" | awk '{print $$4}' | sed "s/\"//g")
# Linux、Illumos
PREFIX=/usr
# FreeBSDとOpenBSD
#PREFIX=/usr/local
# NetBSD
#PREFIX=/usr/pkg
MANPREFIX=${PREFIX}/share/man
# Linux、Illumos、とOpenBSD
CNFPREFIX=/etc
# FreeBSD
#CNFPREFIX=/usr/local/etc
# NetBSD
#CNFPREFIX=/usr/pkg/etc
CC=CGO_ENABLED=0 go build
# リリース。なし＝デバッグ。
RELEASE=-ldflags="-s -w" -buildvcs=false

all:
	${CC} ${RELEASE} -o ${NAME}

release:
	env GOOS=linux GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-amd64
	env GOOS=linux GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-arm64
	env GOOS=linux GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-arm
	env GOOS=linux GOARCH=riscv64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-riscv64
	env GOOS=linux GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-386
	env GOOS=linux GOARCH=ppc64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-ppc64
	env GOOS=linux GOARCH=mips64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-mips64
	env GOOS=openbsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-amd64
	env GOOS=openbsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-386
	env GOOS=openbsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-arm64
	env GOOS=openbsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-arm
	env GOOS=openbsd GOARCH=mips64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-mips64
	env GOOS=freebsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-amd64
	env GOOS=freebsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-386
	env GOOS=freebsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-arm4
	env GOOS=freebsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-arm64
	env GOOS=freebsd GOARCH=riscv64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-riscv64
	env GOOS=netbsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-amd64
	env GOOS=netbsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-386
	env GOOS=netbsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-arm4
	env GOOS=netbsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-arm64
	env GOOS=illumos GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-illumos-amd64

clean:
	rm -f ${NAME} ${NAME}-${VERSION}.tar.gz

dist: clean
	mkdir -p ${NAME}-${VERSION}
	cp -R LICENSE.txt Makefile README.md CHANGELOG.md\
		view static logo.jpg\
		${NAME}.1 *.go *.json go.mod ${NAME}-${VERSION}
	tar zcfv ${NAME}-${VERSION}.tar.gz ${NAME}-${VERSION}
	rm -rf ${NAME}-${VERSION}

install: all
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f ${NAME} ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/${NAME}
	mkdir -p ${DESTDIR}${MANPREFIX}/man1
	sed "s/VERSION/${VERSION}/g" < ${NAME}.1 > ${DESTDIR}${MANPREFIX}/man1/${NAME}.1
	chmod 644 ${DESTDIR}${MANPREFIX}/man1/${NAME}.1
	mkdir -p ${DESTDIR}${PREFIX}/share/${NAME}/archive
	chmod 755 ${DESTDIR}${PREFIX}/share/${NAME}/archive
	mkdir -p ${DESTDIR}${CNFPREFIX}/${NAME}
	chmod 755 ${DESTDIR}${CNFPREFIX}/${NAME}

uninstall:
	rm -f ${DESTDIOR}${PREFIX}/bin/${NAME}\
		${DESTDIR}${MANPREFIX}/man1/${NAME}.1\
		${DESTDIR}${CNFPREFIX}/${NAME}\
		${DESTDIR}${PREFIX}/share/${NAME}

.PHONY: all release clean dist install uninstall

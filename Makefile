GOCMD=go
GOBUILD=$(GOCMD) build
GOENV=$(GOCMD) env
OUTDIR=out
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.version=${VERSION}"
NAME=server
MAIN=main.go

build:
	@rm -fr ${OUTDIR}
	@mkdir -p ${OUTDIR}
	GOOS=linux GOARCH=arm GOARM=6 ${GOBUILD} ${LDFLAGS} -o ${OUTDIR}/${NAME}.rpi ${MAIN}
	GOOS=linux ${GOBUILD}  ${LDFLAGS} -o ${OUTDIR}/${NAME}.lin ${MAIN}
	GOOS=darwin ${GOBUILD} ${LDFLAGS} -o ${OUTDIR}/${NAME}.mac ${MAIN}
	zip ${OUTDIR}/${NAME}.rpi.zip ${OUTDIR}/${NAME}.rpi 
	zip ${OUTDIR}/${NAME}.lin.zip ${OUTDIR}/${NAME}.lin
	zip ${OUTDIR}/${NAME}.mac.zip ${OUTDIR}/${NAME}.mac

alpine:
	GOOS=linux ${GOBUILD}  ${LDFLAGS} -o ${OUTDIR}/${NAME}.lin ${MAIN}

test:
	go test ./... -v

coverage:
	@scripts/coverage.sh

image:
	sudo docker build -t thiagozs/paccount .

run.docker:
	sudo docker run --rm --name=paccount --publish=8080:8080 thiagozs/paccount:latest

swagger:
	@swag init -g main.go
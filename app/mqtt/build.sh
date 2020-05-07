#!/bin/bash
# This is how we want to name the binary output
OUTPUT=main

Version="v0.1"
Author="MACPower"
# These are the values we want to pass for Version and BuildTime
#GitTag=`git describe --tags`
GitTag="MAC`git describe --tags`"

BuildTime=`date +%FT%T%z`
GitCommit=`git rev-parse HEAD`
GoVersion=`go version | sed 's/[ ][ ]*/_/g'`

# Setup the -ldflags option for go build here, interpolate the variable values
# shellcheck disable=SC2027
LDFLAGS="-X ${OUTPUT}.Version=${Version} -X ${OUTPUT}.GitTag=${GitTag} -X ${OUTPUT}.BuildTime=${BuildTime} -X ${OUTPUT}.GitCommit=${GitCommit} -X ${OUTPUT}.GoVersion=${GoVersion} -X ${OUTPUT}.Author=${Author}"

MACOS="CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 "
LNXOS="CGO_ENABLED=0 GOOS=linux GOARCH=amd64 "
WINOS="CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "

rm -f -r ./release $OUTPUT.* $OUTPUT ./run.sh *.pem Sysinfo.json *.tgz
mkdir -p ./release

CMD="go build"

echo $MACOS $CMD -ldflags \"$LDFLAGS\" -o ./release/$OUTPUT.mac  > ./run.sh
echo $LNXOS $CMD -ldflags \"$LDFLAGS\" -o ./release/$OUTPUT.lnx >> ./run.sh
echo $WINOS $CMD -ldflags \"$LDFLAGS\" -o ./release/$OUTPUT.exe >> ./run.sh

echo go build -ldflags \"$LDFLAGS\" -o ./$OUTPUT.exe >> ./run.sh

chmod +x ./run.sh
./run.sh
timestamp=`date +%Y%m%d_%H%S_%s`
tar zcvf app${timestamp}.tgz release config *.bat main.exe

#${MACOS} go build -ldflags "$LDFLAGS" -o ./MAC/${OUTPUT}
#${LNXOS} go build -ldflags "$LDFLAGS" -o ./LNX/${OUTPUT}
#${WINOS} go build -ldflags "$LDFLAGS" -o ./WIN/${OUTPUT}

#go build -ldflags "$LDFLAGS" -o ${OUTPUT}
#./${OUTPUT} --version



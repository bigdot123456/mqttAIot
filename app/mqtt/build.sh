#!/bin/bash
# This is how we want to name the binary output
OUTPUT=main

Version="v0.1"
Author="bigdot"
# These are the values we want to pass for Version and BuildTime
#GitTag=`git describe --tags`
GitTag="Tag_`git describe --tags`"
BuildTime=`date +%FT%T%z`
GitCommit=`git rev-parse HEAD`
GoVersion=`go version | sed 's/[ ][ ]*/_/g'`

# Setup the -ldflags option for go build here, interpolate the variable values
# shellcheck disable=SC2027
LDFLAGS="-X ${OUTPUT}.Version=${Version} -X ${OUTPUT}.GitTag=${GitTag} -X ${OUTPUT}.BuildTime=${BuildTime} -X ${OUTPUT}.GitCommit=${GitCommit} -X ${OUTPUT}.GoVersion=${GoVersion} -X ${OUTPUT}.Author=${Author}"

MACOS="CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 "
LNXOS="CGO_ENABLED=0 GOOS=linux GOARCH=amd64 "
WINOS="CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "

rm -f -r ./release $OUTPUT
mkdir -p ./release

CMD="go build"

echo $MACOS $CMD -ldflags \"$LDFLAGS\" -o ./release/MAC.$OUTPUT  > ./run.sh
echo $LNXOS $CMD -ldflags \"$LDFLAGS\" -o ./release/LNX.$OUTPUT >> ./run.sh
echo $WINOS $CMD -ldflags \"$LDFLAGS\" -o ./release/$OUTPUT.exe >> ./run.sh

echo go build -ldflags \"$LDFLAGS\" -o $OUTPUT >> ./run.sh

chmod +x ./run.sh
./run.sh

#${MACOS} go build -ldflags "$LDFLAGS" -o ./MAC/${OUTPUT}
#${LNXOS} go build -ldflags "$LDFLAGS" -o ./LNX/${OUTPUT}
#${WINOS} go build -ldflags "$LDFLAGS" -o ./WIN/${OUTPUT}

#go build -ldflags "$LDFLAGS" -o ${OUTPUT}
#./${OUTPUT} --version



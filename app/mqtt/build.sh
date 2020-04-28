# This is how we want to name the binary output
# This is how we want to name the binary output
OUTPUT=main

Version="v0.1"
Author="bigdot"
# These are the values we want to pass for Version and BuildTime
#GitTag=`git describe --tags`
GitTag="Tag_`git describe --tags`"
BuildTime=`date +%FT%T%z`
GitCommit=`git rev-parse HEAD`
#GoVersion=`go version | sed 's/[ ][ ]*/_/g'`

# Setup the -ldflags option for go build here, interpolate the variable values
# shellcheck disable=SC2027
LDFLAGS="-X ${OUTPUT}.Version=${Version} -X ${OUTPUT}.GitTag=${GitTag} -X ${OUTPUT}.BuildTime=${BuildTime} -X ${OUTPUT}.GitCommit=${GitCommit} -X ${OUTPUT}.GoVersion=${GoVersion} -X ${OUTPUT}.Author=${Author}"

go build -ldflags "$LDFLAGS" -o ${OUTPUT}
#echo go build -ldflags \"$LDFLAGS\"  > run.sh
#chmod +x run.sh
#./run.sh


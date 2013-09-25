#!/usr/bin/env zsh

make

./git-version-proxy &
go get -x 127.0.0.1:8080/gh/msiebuhr/master/dummyGraphiteData.git
# Run the printed git commands with GIT_TRACE_PACKET=1 set in order to see what
# they send over the wire.

kill %1


#!/bin/bash

curpath=$(cd "$(dirname "$0")"; pwd)
project=$(cd "$curpath/.."; pwd)

binary='gogofast'

include_path="-I=."
include_path="$include_path -I=$GOPATH/src" 
include_path="$include_path -I=$GOPATH/src/github.com/gogo/protobuf/protobuf"
output_flag='Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types'
output_flag=$output_flag:.

cd $project

protoc $include_path --${binary}_out=$output_flag web/msg.proto

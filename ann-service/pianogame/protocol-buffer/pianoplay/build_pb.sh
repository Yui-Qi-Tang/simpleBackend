# !/bin/sh
protoc -I . ./pianoplay.proto --go_out=plugins=grpc:.
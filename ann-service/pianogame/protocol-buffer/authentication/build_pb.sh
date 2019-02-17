# !/bin/sh
protoc -I . ./authentication.proto --go_out=plugins=grpc:.
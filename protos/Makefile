FLAG_GEN_gRPC_OUT=--go-grpc_out=./gen/go/
FLAG_GEN_gRPC_OPT=--go-grpc_opt=paths=source_relative

FLAG_GEN_GO_OUT=--go_out=./gen/go/
FLAG_GEN_GO_OPT=--go_opt=paths=source_relative


gen_protoc:
	protoc -I proto proto/sso/sso.proto $(FLAG_GEN_GO_OUT) $(FLAG_GEN_GO_OPT) $(FLAG_GEN_gRPC_OUT) $(FLAG_GEN_gRPC_OPT)

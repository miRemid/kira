.PHONY: proto
proto:
	protoc --proto_path=. auth.proto --micro_out=pb/auth --go_out=pb/auth
	protoc --proto_path=. file.proto --micro_out=pb/file --go_out=pb/file
	protoc --proto_path=. user.proto --micro_out=pb/user --go_out=pb/user
	protoc --proto_path=. upload.proto --micro_out=pb/upload --go_out=pb/upload
	protoc-go-inject-tag -input=pb/auth/auth.pb.go
	protoc-go-inject-tag -input=pb/file/file.pb.go
	protoc-go-inject-tag -input=pb/user/user.pb.go
	protoc-go-inject-tag -input=pb/upload/upload.pb.go

pro:
	protoc --proto_path=. *.proto --micro_out=pb --go_out=pb
	protoc-go-inject-tag -input=pb/auth.pb.go
	protoc-go-inject-tag -input=pb/file.pb.go
	protoc-go-inject-tag -input=pb/user.pb.go
	protoc-go-inject-tag -input=pb/upload.pb.go
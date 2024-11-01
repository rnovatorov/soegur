package api

//go:generate protoc -I ../../proto --go_opt=paths=source_relative --go_out=../.. internal/api/sagaspecpb/spec.proto
//go:generate protoc -I ../../proto --go_opt=paths=source_relative --go_out=../.. internal/api/sagaeventspb/events.proto

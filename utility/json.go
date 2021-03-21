package utility

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func MarshalToJson(message proto.Message) ([]byte, error) {
	options := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	return options.Marshal(message)
}

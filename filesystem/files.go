package filesystem

import "github.com/golang/protobuf/proto"

//go:generate protoc --go_out=. files.proto

func MarshalMetaData(md *MetaData) ([]byte, error) {
	return proto.Marshal(&MetaData{
		Name: proto.String(*md.Name),
		Size: proto.Int64(int64(*md.Size)),
		Sha1: proto.String(*md.Sha1),
	})
}

/*func UnmarshalMetaData(data []byte, md *MetaData) error {
	var pmd MetaData
	if err := proto.Unmarshal(data, &pmd); err != nil {
		return err
	}
	md.Name = pmd.GetName()
}*/

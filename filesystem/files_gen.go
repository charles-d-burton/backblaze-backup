package filesystem

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *MsgpMetaData) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "size":
			z.Size, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "sha1":
			z.Sha1, err = dc.ReadString()
			if err != nil {
				return
			}
		case "backed":
			z.BackedUp, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "versions":
			var zbai uint32
			zbai, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Versions) >= int(zbai) {
				z.Versions = (z.Versions)[:zbai]
			} else {
				z.Versions = make([]string, zbai)
			}
			for zxvk := range z.Versions {
				z.Versions[zxvk], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "delete":
			z.Delete, err = dc.ReadBool()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *MsgpMetaData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "name"
	err = en.Append(0x86, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "size"
	err = en.Append(0xa4, 0x73, 0x69, 0x7a, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.Size)
	if err != nil {
		return
	}
	// write "sha1"
	err = en.Append(0xa4, 0x73, 0x68, 0x61, 0x31)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Sha1)
	if err != nil {
		return
	}
	// write "backed"
	err = en.Append(0xa6, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.BackedUp)
	if err != nil {
		return
	}
	// write "versions"
	err = en.Append(0xa8, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Versions)))
	if err != nil {
		return
	}
	for zxvk := range z.Versions {
		err = en.WriteString(z.Versions[zxvk])
		if err != nil {
			return
		}
	}
	// write "delete"
	err = en.Append(0xa6, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.Delete)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *MsgpMetaData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "name"
	o = append(o, 0x86, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "size"
	o = append(o, 0xa4, 0x73, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt64(o, z.Size)
	// string "sha1"
	o = append(o, 0xa4, 0x73, 0x68, 0x61, 0x31)
	o = msgp.AppendString(o, z.Sha1)
	// string "backed"
	o = append(o, 0xa6, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x64)
	o = msgp.AppendBool(o, z.BackedUp)
	// string "versions"
	o = append(o, 0xa8, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Versions)))
	for zxvk := range z.Versions {
		o = msgp.AppendString(o, z.Versions[zxvk])
	}
	// string "delete"
	o = append(o, 0xa6, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65)
	o = msgp.AppendBool(o, z.Delete)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MsgpMetaData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "size":
			z.Size, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "sha1":
			z.Sha1, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "backed":
			z.BackedUp, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "versions":
			var zajw uint32
			zajw, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Versions) >= int(zajw) {
				z.Versions = (z.Versions)[:zajw]
			} else {
				z.Versions = make([]string, zajw)
			}
			for zxvk := range z.Versions {
				z.Versions[zxvk], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "delete":
			z.Delete, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *MsgpMetaData) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.Int64Size + 5 + msgp.StringPrefixSize + len(z.Sha1) + 7 + msgp.BoolSize + 9 + msgp.ArrayHeaderSize
	for zxvk := range z.Versions {
		s += msgp.StringPrefixSize + len(z.Versions[zxvk])
	}
	s += 7 + msgp.BoolSize
	return
}

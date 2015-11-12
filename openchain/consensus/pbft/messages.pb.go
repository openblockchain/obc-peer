// Code generated by protoc-gen-go.
// source: messages.proto
// DO NOT EDIT!

/*
Package pbft is a generated protocol buffer package.

It is generated from these files:
	messages.proto

It has these top-level messages:
	Message
	Request
	PrePrepare
	Prepare
	Commit
*/
package pbft

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "google/protobuf"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Message struct {
	// Types that are valid to be assigned to Payload:
	//	*Message_Request
	//	*Message_PrePrepare
	//	*Message_Prepare
	//	*Message_Commit
	Payload isMessage_Payload `protobuf_oneof:"payload"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

type isMessage_Payload interface {
	isMessage_Payload()
}

type Message_Request struct {
	Request *Request `protobuf:"bytes,1,opt,name=request,oneof"`
}
type Message_PrePrepare struct {
	PrePrepare *PrePrepare `protobuf:"bytes,2,opt,name=prePrepare,oneof"`
}
type Message_Prepare struct {
	Prepare *Prepare `protobuf:"bytes,3,opt,name=prepare,oneof"`
}
type Message_Commit struct {
	Commit *Commit `protobuf:"bytes,4,opt,name=commit,oneof"`
}

func (*Message_Request) isMessage_Payload()    {}
func (*Message_PrePrepare) isMessage_Payload() {}
func (*Message_Prepare) isMessage_Payload()    {}
func (*Message_Commit) isMessage_Payload()     {}

func (m *Message) GetPayload() isMessage_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetRequest() *Request {
	if x, ok := m.GetPayload().(*Message_Request); ok {
		return x.Request
	}
	return nil
}

func (m *Message) GetPrePrepare() *PrePrepare {
	if x, ok := m.GetPayload().(*Message_PrePrepare); ok {
		return x.PrePrepare
	}
	return nil
}

func (m *Message) GetPrepare() *Prepare {
	if x, ok := m.GetPayload().(*Message_Prepare); ok {
		return x.Prepare
	}
	return nil
}

func (m *Message) GetCommit() *Commit {
	if x, ok := m.GetPayload().(*Message_Commit); ok {
		return x.Commit
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Message) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), []interface{}) {
	return _Message_OneofMarshaler, _Message_OneofUnmarshaler, []interface{}{
		(*Message_Request)(nil),
		(*Message_PrePrepare)(nil),
		(*Message_Prepare)(nil),
		(*Message_Commit)(nil),
	}
}

func _Message_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Message)
	// payload
	switch x := m.Payload.(type) {
	case *Message_Request:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Request); err != nil {
			return err
		}
	case *Message_PrePrepare:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PrePrepare); err != nil {
			return err
		}
	case *Message_Prepare:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Prepare); err != nil {
			return err
		}
	case *Message_Commit:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Commit); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Message.Payload has unexpected type %T", x)
	}
	return nil
}

func _Message_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Message)
	switch tag {
	case 1: // payload.request
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Request)
		err := b.DecodeMessage(msg)
		m.Payload = &Message_Request{msg}
		return true, err
	case 2: // payload.prePrepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PrePrepare)
		err := b.DecodeMessage(msg)
		m.Payload = &Message_PrePrepare{msg}
		return true, err
	case 3: // payload.prepare
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Prepare)
		err := b.DecodeMessage(msg)
		m.Payload = &Message_Prepare{msg}
		return true, err
	case 4: // payload.commit
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Commit)
		err := b.DecodeMessage(msg)
		m.Payload = &Message_Commit{msg}
		return true, err
	default:
		return false, nil
	}
}

type Request struct {
	Timestamp   *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Transaction []byte                     `protobuf:"bytes,2,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}

func (m *Request) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type PrePrepare struct {
	View           uint64   `protobuf:"varint,1,opt,name=view" json:"view,omitempty"`
	SequenceNumber uint64   `protobuf:"varint,2,opt,name=sequenceNumber" json:"sequenceNumber,omitempty"`
	RequestDigest  [][]byte `protobuf:"bytes,3,rep,name=requestDigest,proto3" json:"requestDigest,omitempty"`
}

func (m *PrePrepare) Reset()         { *m = PrePrepare{} }
func (m *PrePrepare) String() string { return proto.CompactTextString(m) }
func (*PrePrepare) ProtoMessage()    {}

type Prepare struct {
	View           uint64 `protobuf:"varint,1,opt,name=view" json:"view,omitempty"`
	SequenceNumber uint64 `protobuf:"varint,2,opt,name=sequenceNumber" json:"sequenceNumber,omitempty"`
	RequestDigest  string `protobuf:"bytes,3,opt,name=requestDigest" json:"requestDigest,omitempty"`
	ReplicaId      uint64 `protobuf:"varint,4,opt,name=replicaId" json:"replicaId,omitempty"`
}

func (m *Prepare) Reset()         { *m = Prepare{} }
func (m *Prepare) String() string { return proto.CompactTextString(m) }
func (*Prepare) ProtoMessage()    {}

type Commit struct {
	View           uint64 `protobuf:"varint,1,opt,name=view" json:"view,omitempty"`
	SequenceNumber uint64 `protobuf:"varint,2,opt,name=sequenceNumber" json:"sequenceNumber,omitempty"`
	ReplicaId      uint64 `protobuf:"varint,3,opt,name=replicaId" json:"replicaId,omitempty"`
}

func (m *Commit) Reset()         { *m = Commit{} }
func (m *Commit) String() string { return proto.CompactTextString(m) }
func (*Commit) ProtoMessage()    {}

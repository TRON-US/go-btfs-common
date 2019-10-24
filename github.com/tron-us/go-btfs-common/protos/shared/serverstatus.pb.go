// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/shared/serverstatus.proto

package shared

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type StatusReport_HealthStatus int32

const (
	StatusReport_WRONG        StatusReport_HealthStatus = 0
	StatusReport_RUNNING      StatusReport_HealthStatus = 1
	StatusReport_BOOTSTRAP    StatusReport_HealthStatus = 2
	StatusReport_PARTIAL_STOP StatusReport_HealthStatus = 3
)

var StatusReport_HealthStatus_name = map[int32]string{
	0: "WRONG",
	1: "RUNNING",
	2: "BOOTSTRAP",
	3: "PARTIAL_STOP",
}

var StatusReport_HealthStatus_value = map[string]int32{
	"WRONG":        0,
	"RUNNING":      1,
	"BOOTSTRAP":    2,
	"PARTIAL_STOP": 3,
}

func (x StatusReport_HealthStatus) String() string {
	return proto.EnumName(StatusReport_HealthStatus_name, int32(x))
}

func (StatusReport_HealthStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ca601dd99a24d24a, []int{1, 0}
}

type StatusRequest struct {
	RequestAddress       []byte               `protobuf:"bytes,1,opt,name=request_address,json=requestAddress,proto3" json:"request_address,omitempty"`
	CurentTime           *timestamp.Timestamp `protobuf:"bytes,2,opt,name=curent_time,json=curentTime,proto3" json:"curent_time,omitempty"`
	Signature            []byte               `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca601dd99a24d24a, []int{0}
}

func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}
func (m *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(m, src)
}
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

func (m *StatusRequest) GetRequestAddress() []byte {
	if m != nil {
		return m.RequestAddress
	}
	return nil
}

func (m *StatusRequest) GetCurentTime() *timestamp.Timestamp {
	if m != nil {
		return m.CurentTime
	}
	return nil
}

func (m *StatusRequest) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type StatusReport struct {
	PeerId               []byte                    `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Address              []byte                    `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	ServiceName          []byte                    `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Status               StatusReport_HealthStatus `protobuf:"varint,4,opt,name=status,proto3,enum=shared.StatusReport_HealthStatus" json:"status,omitempty"`
	CurentTime           *timestamp.Timestamp      `protobuf:"bytes,5,opt,name=curent_time,json=curentTime,proto3" json:"curent_time,omitempty"`
	StartTime            *timestamp.Timestamp      `protobuf:"bytes,6,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	GitHash              []byte                    `protobuf:"bytes,7,opt,name=git_hash,json=gitHash,proto3" json:"git_hash,omitempty"`
	Version              []byte                    `protobuf:"bytes,8,opt,name=version,proto3" json:"version,omitempty"`
	DbStatusExtra        []byte                    `protobuf:"bytes,9,opt,name=db_status_extra,json=dbStatusExtra,proto3" json:"db_status_extra,omitempty"`
	QueueStatusExtra     []byte                    `protobuf:"bytes,10,opt,name=queue_status_extra,json=queueStatusExtra,proto3" json:"queue_status_extra,omitempty"`
	ChainStatusExtra     []byte                    `protobuf:"bytes,11,opt,name=chain_status_extra,json=chainStatusExtra,proto3" json:"chain_status_extra,omitempty"`
	CacheStatusExtra     []byte                    `protobuf:"bytes,12,opt,name=cache_status_extra,json=cacheStatusExtra,proto3" json:"cache_status_extra,omitempty"`
	Extra                []byte                    `protobuf:"bytes,13,opt,name=extra,proto3" json:"extra,omitempty"`
	Signature            []byte                    `protobuf:"bytes,14,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *StatusReport) Reset()         { *m = StatusReport{} }
func (m *StatusReport) String() string { return proto.CompactTextString(m) }
func (*StatusReport) ProtoMessage()    {}
func (*StatusReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca601dd99a24d24a, []int{1}
}

func (m *StatusReport) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusReport.Unmarshal(m, b)
}
func (m *StatusReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusReport.Marshal(b, m, deterministic)
}
func (m *StatusReport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusReport.Merge(m, src)
}
func (m *StatusReport) XXX_Size() int {
	return xxx_messageInfo_StatusReport.Size(m)
}
func (m *StatusReport) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusReport.DiscardUnknown(m)
}

var xxx_messageInfo_StatusReport proto.InternalMessageInfo

func (m *StatusReport) GetPeerId() []byte {
	if m != nil {
		return m.PeerId
	}
	return nil
}

func (m *StatusReport) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *StatusReport) GetServiceName() []byte {
	if m != nil {
		return m.ServiceName
	}
	return nil
}

func (m *StatusReport) GetStatus() StatusReport_HealthStatus {
	if m != nil {
		return m.Status
	}
	return StatusReport_WRONG
}

func (m *StatusReport) GetCurentTime() *timestamp.Timestamp {
	if m != nil {
		return m.CurentTime
	}
	return nil
}

func (m *StatusReport) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *StatusReport) GetGitHash() []byte {
	if m != nil {
		return m.GitHash
	}
	return nil
}

func (m *StatusReport) GetVersion() []byte {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *StatusReport) GetDbStatusExtra() []byte {
	if m != nil {
		return m.DbStatusExtra
	}
	return nil
}

func (m *StatusReport) GetQueueStatusExtra() []byte {
	if m != nil {
		return m.QueueStatusExtra
	}
	return nil
}

func (m *StatusReport) GetChainStatusExtra() []byte {
	if m != nil {
		return m.ChainStatusExtra
	}
	return nil
}

func (m *StatusReport) GetCacheStatusExtra() []byte {
	if m != nil {
		return m.CacheStatusExtra
	}
	return nil
}

func (m *StatusReport) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

func (m *StatusReport) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterEnum("shared.StatusReport_HealthStatus", StatusReport_HealthStatus_name, StatusReport_HealthStatus_value)
	proto.RegisterType((*StatusRequest)(nil), "shared.StatusRequest")
	proto.RegisterType((*StatusReport)(nil), "shared.StatusReport")
}

func init() { proto.RegisterFile("protos/shared/serverstatus.proto", fileDescriptor_ca601dd99a24d24a) }

var fileDescriptor_ca601dd99a24d24a = []byte{
	// 488 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x4d, 0x6f, 0x9b, 0x4c,
	0x10, 0xc7, 0x1f, 0x9c, 0x18, 0x87, 0x01, 0x3b, 0x68, 0xf5, 0x48, 0xa5, 0x51, 0xa5, 0x3a, 0x39,
	0xb4, 0x3e, 0xd4, 0xa0, 0xa6, 0xa7, 0xa8, 0x27, 0x47, 0xaa, 0x1c, 0x4b, 0x15, 0xb6, 0x30, 0x55,
	0xa5, 0x5e, 0xd0, 0x02, 0x1b, 0x40, 0x0a, 0xac, 0xb3, 0x2f, 0x55, 0xbf, 0x47, 0x3f, 0x6e, 0x2f,
	0x15, 0xbb, 0xa0, 0xda, 0x5c, 0xaa, 0xde, 0x3c, 0xff, 0xff, 0x6f, 0x76, 0x5e, 0x3c, 0xc0, 0xfc,
	0xc0, 0xa8, 0xa0, 0x3c, 0xe0, 0x25, 0x66, 0x24, 0x0f, 0x38, 0x61, 0xdf, 0x09, 0xe3, 0x02, 0x0b,
	0xc9, 0x7d, 0x65, 0x21, 0x53, 0x5b, 0x57, 0xaf, 0x0b, 0x4a, 0x8b, 0x27, 0x12, 0x28, 0x35, 0x95,
	0x8f, 0x81, 0xa8, 0x6a, 0xc2, 0x05, 0xae, 0x0f, 0x1a, 0xbc, 0xf9, 0x69, 0xc0, 0x74, 0xaf, 0x32,
	0x23, 0xf2, 0x2c, 0x09, 0x17, 0xe8, 0x2d, 0x5c, 0x32, 0xfd, 0x33, 0xc1, 0x79, 0xce, 0x08, 0xe7,
	0x9e, 0x31, 0x37, 0x16, 0x4e, 0x34, 0xeb, 0xe4, 0x95, 0x56, 0xd1, 0x47, 0xb0, 0x33, 0xc9, 0x48,
	0x23, 0x92, 0xf6, 0x51, 0x6f, 0x34, 0x37, 0x16, 0xf6, 0xed, 0x95, 0xaf, 0x2b, 0xfa, 0x7d, 0x45,
	0x3f, 0xee, 0x2b, 0x46, 0xa0, 0xf1, 0x56, 0x40, 0xaf, 0xc0, 0xe2, 0x55, 0xd1, 0x60, 0x21, 0x19,
	0xf1, 0xce, 0xd4, 0xfb, 0x7f, 0x84, 0x9b, 0x5f, 0xe7, 0xe0, 0xf4, 0x5d, 0x1d, 0x28, 0x13, 0xe8,
	0x05, 0x4c, 0x0e, 0x84, 0xb0, 0xa4, 0xca, 0xbb, 0x66, 0xcc, 0x36, 0xdc, 0xe4, 0xc8, 0x83, 0x49,
	0xdf, 0xe5, 0x48, 0x19, 0x7d, 0x88, 0xae, 0xc1, 0x69, 0x17, 0x53, 0x65, 0x24, 0x69, 0x70, 0xdd,
	0x17, 0xb1, 0x3b, 0x2d, 0xc4, 0x35, 0x41, 0x77, 0x60, 0xea, 0xad, 0x79, 0xe7, 0x73, 0x63, 0x31,
	0xbb, 0xbd, 0xf6, 0xf5, 0xda, 0xfc, 0xe3, 0xda, 0xfe, 0x03, 0xc1, 0x4f, 0xa2, 0xec, 0xa4, 0x2e,
	0x61, 0x38, 0xfc, 0xf8, 0x9f, 0x86, 0xbf, 0x03, 0xe0, 0x02, 0xb3, 0x2e, 0xd7, 0xfc, 0x6b, 0xae,
	0xa5, 0x68, 0x95, 0xfa, 0x12, 0x2e, 0x8a, 0x4a, 0x24, 0x25, 0xe6, 0xa5, 0x37, 0xd1, 0x03, 0x17,
	0x95, 0x78, 0xc0, 0xbc, 0x6c, 0x57, 0xd1, 0x9e, 0x41, 0x45, 0x1b, 0xef, 0x42, 0x3b, 0x5d, 0x88,
	0xde, 0xc0, 0x65, 0x9e, 0x26, 0xba, 0xf3, 0x84, 0xfc, 0x10, 0x0c, 0x7b, 0x96, 0x22, 0xa6, 0x79,
	0xaa, 0xe7, 0xfa, 0xd4, 0x8a, 0xe8, 0x1d, 0xa0, 0x67, 0x49, 0x24, 0x39, 0x45, 0x41, 0xa1, 0xae,
	0x72, 0x06, 0x74, 0x56, 0xe2, 0xaa, 0x39, 0xa5, 0x6d, 0x4d, 0x2b, 0x67, 0x48, 0xe3, 0xac, 0x1c,
	0xbc, 0xed, 0x74, 0x74, 0xeb, 0x1c, 0xd3, 0xff, 0xc3, 0x58, 0x03, 0x53, 0x05, 0xe8, 0xe0, 0xf4,
	0x68, 0x66, 0xc3, 0xa3, 0x59, 0x83, 0x73, 0xfc, 0x57, 0x21, 0x0b, 0xc6, 0x5f, 0xa3, 0x6d, 0xb8,
	0x76, 0xff, 0x43, 0x36, 0x4c, 0xa2, 0x2f, 0x61, 0xb8, 0x09, 0xd7, 0xae, 0x81, 0xa6, 0x60, 0xdd,
	0x6f, 0xb7, 0xf1, 0x3e, 0x8e, 0x56, 0x3b, 0x77, 0x84, 0x5c, 0x70, 0x76, 0xab, 0x28, 0xde, 0xac,
	0x3e, 0x27, 0xfb, 0x78, 0xbb, 0x73, 0xcf, 0xee, 0xdf, 0x7f, 0x0b, 0x8a, 0x4a, 0x94, 0x32, 0xf5,
	0x33, 0x5a, 0x07, 0x82, 0xd1, 0x66, 0x29, 0x79, 0x50, 0xd0, 0x65, 0x2a, 0x1e, 0xf9, 0x32, 0xa3,
	0x75, 0x4d, 0x9b, 0xe0, 0xe4, 0x13, 0x4c, 0x4d, 0x15, 0x7e, 0xf8, 0x1d, 0x00, 0x00, 0xff, 0xff,
	0xf0, 0x61, 0xc2, 0xed, 0x9a, 0x03, 0x00, 0x00,
}
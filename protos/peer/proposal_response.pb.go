// Code generated by protoc-gen-go.
// source: peer/proposal_response.proto
// DO NOT EDIT!

package peer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// A ProposalResponse is returned from an endorser to the proposal submitter.
// The idea is that this message contains the endorser's response to the
// request of a client to perform an action over a chaincode (or more
// generically on the ledger); the response might be success/error (conveyed in
// the Response field) together with a description of the action and a
// signature over it by that endorser.  If a sufficient number of distinct
// endorsers agree on the same action and produce signature to that effect, a
// transaction can be generated and sent for ordering.
type ProposalResponse struct {
	// Version indicates message protocol version
	Version int32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	// Timestamp is the time that the message
	// was created as  defined by the sender
	Timestamp *google_protobuf1.Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	// A response message indicating whether the
	// endorsement of the action was successful
	Response *Response `protobuf:"bytes,4,opt,name=response" json:"response,omitempty"`
	// The payload of response. It is the bytes of ProposalResponsePayload
	Payload []byte `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	// The endorsement of the proposal, basically
	// the endorser's signature over the payload
	Endorsement *Endorsement `protobuf:"bytes,6,opt,name=endorsement" json:"endorsement,omitempty"`
}

func (m *ProposalResponse) Reset()                    { *m = ProposalResponse{} }
func (m *ProposalResponse) String() string            { return proto.CompactTextString(m) }
func (*ProposalResponse) ProtoMessage()               {}
func (*ProposalResponse) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *ProposalResponse) GetTimestamp() *google_protobuf1.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ProposalResponse) GetResponse() *Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *ProposalResponse) GetEndorsement() *Endorsement {
	if m != nil {
		return m.Endorsement
	}
	return nil
}

// A response with a representation similar to an HTTP response that can
// be used within another message.
type Response struct {
	// A status code that should follow the HTTP status codes.
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	// A message associated with the response code.
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	// A payload that can be used to include metadata with this response.
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{1} }

// ProposalResponsePayload is the payload of a proposal response.  This message
// is the "bridge" between the client's request and the endorser's action in
// response to that request. Concretely, for chaincodes, it contains a hashed
// representation of the proposal (proposalHash) and a representation of the
// chaincode state changes and events inside the extension field.
type ProposalResponsePayload struct {
	// Hash of the proposal that triggered this response. The hash is used to
	// link a response with its proposal, both for bookeeping purposes on an
	// asynchronous system and for security reasons (accountability,
	// non-repudiation). The hash usually covers the entire Proposal message
	// (byte-by-byte). However this implies that the hash can only be verified
	// if the entire proposal message is available when ProposalResponsePayload is
	// included in a transaction or stored in the ledger. For confidentiality
	// reasons, with chaincodes it might be undesirable to store the proposal
	// payload in the ledger.  If the type is CHAINCODE, this is handled by
	// separating the proposal's header and
	// the payload: the header is always hashed in its entirety whereas the
	// payload can either be hashed fully, or only its hash may be hashed, or
	// nothing from the payload can be hashed. The PayloadVisibility field in the
	// Header's extension controls to which extent the proposal payload is
	// "visible" in the sense that was just explained.
	ProposalHash []byte `protobuf:"bytes,1,opt,name=proposal_hash,json=proposalHash,proto3" json:"proposal_hash,omitempty"`
	// Extension should be unmarshaled to a type-specific message. The type of
	// the extension in any proposal response depends on the type of the proposal
	// that the client selected when the proposal was initially sent out.  In
	// particular, this information is stored in the type field of a Header.  For
	// chaincode, it's a ChaincodeAction message
	Extension []byte `protobuf:"bytes,2,opt,name=extension,proto3" json:"extension,omitempty"`
}

func (m *ProposalResponsePayload) Reset()                    { *m = ProposalResponsePayload{} }
func (m *ProposalResponsePayload) String() string            { return proto.CompactTextString(m) }
func (*ProposalResponsePayload) ProtoMessage()               {}
func (*ProposalResponsePayload) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{2} }

// An endorsement is a signature of an endorser over a proposal response.  By
// producing an endorsement message, an endorser implicitly "approves" that
// proposal response and the actions contained therein. When enough
// endorsements have been collected, a transaction can be generated out of a
// set of proposal responses.  Note that this message only contains an identity
// and a signature but no signed payload. This is intentional because
// endorsements are supposed to be collected in a transaction, and they are all
// expected to endorse a single proposal response/action (many endorsements
// over a single proposal response)
type Endorsement struct {
	// Identity of the endorser (e.g. its certificate)
	Endorser []byte `protobuf:"bytes,1,opt,name=endorser,proto3" json:"endorser,omitempty"`
	// Signature of the payload included in ProposalResponse concatenated with
	// the endorser's certificate; ie, sign(ProposalResponse.payload + endorser)
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *Endorsement) Reset()                    { *m = Endorsement{} }
func (m *Endorsement) String() string            { return proto.CompactTextString(m) }
func (*Endorsement) ProtoMessage()               {}
func (*Endorsement) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{3} }

func init() {
	proto.RegisterType((*ProposalResponse)(nil), "protos.ProposalResponse")
	proto.RegisterType((*Response)(nil), "protos.Response")
	proto.RegisterType((*ProposalResponsePayload)(nil), "protos.ProposalResponsePayload")
	proto.RegisterType((*Endorsement)(nil), "protos.Endorsement")
}

func init() { proto.RegisterFile("peer/proposal_response.proto", fileDescriptor7) }

var fileDescriptor7 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x52, 0x5f, 0x4b, 0xfb, 0x30,
	0x14, 0xa5, 0xfb, 0xfd, 0x36, 0xb7, 0xbb, 0x09, 0xa3, 0x82, 0x96, 0x31, 0x70, 0xd4, 0x97, 0x89,
	0xd2, 0x82, 0x22, 0xf8, 0x2c, 0x88, 0x3e, 0x8e, 0x20, 0x3e, 0x88, 0x20, 0xe9, 0x76, 0xd7, 0x16,
	0xda, 0x26, 0xe4, 0xa6, 0xe2, 0x3e, 0xb0, 0xdf, 0x43, 0x9a, 0x26, 0xdd, 0xf4, 0xa9, 0x9c, 0xdb,
	0x93, 0xf3, 0x27, 0xb9, 0x30, 0x97, 0x88, 0x2a, 0x96, 0x4a, 0x48, 0x41, 0xbc, 0xf8, 0x50, 0x48,
	0x52, 0x54, 0x84, 0x91, 0x54, 0x42, 0x0b, 0x7f, 0x60, 0x3e, 0x34, 0x3b, 0x4f, 0x85, 0x48, 0x0b,
	0x8c, 0x0d, 0x4c, 0xea, 0x6d, 0xac, 0xf3, 0x12, 0x49, 0xf3, 0x52, 0xb6, 0xc4, 0xf0, 0xdb, 0x83,
	0xe9, 0xca, 0x8a, 0x30, 0xab, 0xe1, 0x07, 0x70, 0xf4, 0x89, 0x8a, 0x72, 0x51, 0x05, 0xde, 0xc2,
	0x5b, 0xf6, 0x99, 0x83, 0xfe, 0x3d, 0x8c, 0x3a, 0x85, 0xa0, 0xb7, 0xf0, 0x96, 0xe3, 0x9b, 0x59,
	0xd4, 0x7a, 0x44, 0xce, 0x23, 0x7a, 0x71, 0x0c, 0xb6, 0x27, 0xfb, 0xd7, 0x30, 0x74, 0x19, 0x83,
	0xff, 0xe6, 0xe0, 0xb4, 0x3d, 0x41, 0x91, 0xf3, 0x65, 0x1d, 0xa3, 0x49, 0x20, 0xf9, 0xae, 0x10,
	0x7c, 0x13, 0xf4, 0x17, 0xde, 0x72, 0xc2, 0x1c, 0xf4, 0xef, 0x60, 0x8c, 0xd5, 0x46, 0x28, 0xc2,
	0x12, 0x2b, 0x1d, 0x0c, 0x8c, 0xd4, 0x89, 0x93, 0x7a, 0xdc, 0xff, 0x62, 0x87, 0xbc, 0xf0, 0x15,
	0x86, 0x5d, 0xbd, 0x53, 0x18, 0x90, 0xe6, 0xba, 0x26, 0xdb, 0xce, 0xa2, 0xc6, 0xb4, 0x44, 0x22,
	0x9e, 0xa2, 0xa9, 0x36, 0x62, 0x0e, 0x1e, 0xc6, 0xf9, 0xf7, 0x2b, 0x4e, 0xf8, 0x0e, 0x67, 0x7f,
	0xaf, 0x6f, 0x65, 0x93, 0x5e, 0xc0, 0x71, 0xf7, 0x3c, 0x19, 0xa7, 0xcc, 0xb8, 0x4d, 0xd8, 0xc4,
	0x0d, 0x9f, 0x39, 0x65, 0xfe, 0x1c, 0x46, 0xf8, 0xa5, 0xb1, 0x32, 0x97, 0xdd, 0x33, 0x84, 0xfd,
	0x20, 0x7c, 0x82, 0xf1, 0x41, 0x23, 0x7f, 0x06, 0x43, 0xdb, 0x49, 0x59, 0xb1, 0x0e, 0x37, 0x42,
	0x94, 0xa7, 0x15, 0xd7, 0xb5, 0x42, 0x27, 0xd4, 0x0d, 0x1e, 0xae, 0xde, 0x2e, 0xd3, 0x5c, 0x67,
	0x75, 0x12, 0xad, 0x45, 0x19, 0x67, 0x3b, 0x89, 0xaa, 0xc0, 0x4d, 0x8a, 0x2a, 0xde, 0xf2, 0x44,
	0xe5, 0xeb, 0x76, 0x41, 0x28, 0x6e, 0x96, 0x2a, 0x69, 0x97, 0xe7, 0xf6, 0x27, 0x00, 0x00, 0xff,
	0xff, 0xe8, 0x57, 0x2f, 0xb2, 0x63, 0x02, 0x00, 0x00,
}
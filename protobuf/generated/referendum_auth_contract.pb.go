// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.11.4
// source: referendum_auth_contract.proto

package client

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReferendumOrganization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The threshold for releasing the proposal.
	ProposalReleaseThreshold *ProposalReleaseThreshold `protobuf:"bytes,1,opt,name=proposal_release_threshold,json=proposalReleaseThreshold,proto3" json:"proposal_release_threshold,omitempty"`
	// The token used during proposal operations.
	TokenSymbol string `protobuf:"bytes,2,opt,name=token_symbol,json=tokenSymbol,proto3" json:"token_symbol,omitempty"`
	// The address of organization.
	OrganizationAddress *Address `protobuf:"bytes,3,opt,name=organization_address,json=organizationAddress,proto3" json:"organization_address,omitempty"`
	// The organizations id.
	OrganizationHash *Hash `protobuf:"bytes,4,opt,name=organization_hash,json=organizationHash,proto3" json:"organization_hash,omitempty"`
	// The proposer whitelist.
	ProposerWhiteList *ProposerWhiteList `protobuf:"bytes,5,opt,name=proposer_white_list,json=proposerWhiteList,proto3" json:"proposer_white_list,omitempty"`
	// The creation token is for organization address generation.
	CreationToken *Hash `protobuf:"bytes,6,opt,name=creation_token,json=creationToken,proto3" json:"creation_token,omitempty"`
}

func (x *ReferendumOrganization) Reset() {
	*x = ReferendumOrganization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_referendum_auth_contract_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReferendumOrganization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReferendumOrganization) ProtoMessage() {}

func (x *ReferendumOrganization) ProtoReflect() protoreflect.Message {
	mi := &file_referendum_auth_contract_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReferendumOrganization.ProtoReflect.Descriptor instead.
func (*ReferendumOrganization) Descriptor() ([]byte, []int) {
	return file_referendum_auth_contract_proto_rawDescGZIP(), []int{0}
}

func (x *ReferendumOrganization) GetProposalReleaseThreshold() *ProposalReleaseThreshold {
	if x != nil {
		return x.ProposalReleaseThreshold
	}
	return nil
}

func (x *ReferendumOrganization) GetTokenSymbol() string {
	if x != nil {
		return x.TokenSymbol
	}
	return ""
}

func (x *ReferendumOrganization) GetOrganizationAddress() *Address {
	if x != nil {
		return x.OrganizationAddress
	}
	return nil
}

func (x *ReferendumOrganization) GetOrganizationHash() *Hash {
	if x != nil {
		return x.OrganizationHash
	}
	return nil
}

func (x *ReferendumOrganization) GetProposerWhiteList() *ProposerWhiteList {
	if x != nil {
		return x.ProposerWhiteList
	}
	return nil
}

func (x *ReferendumOrganization) GetCreationToken() *Hash {
	if x != nil {
		return x.CreationToken
	}
	return nil
}

var File_referendum_auth_contract_proto protoreflect.FileDescriptor

var file_referendum_auth_contract_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x64, 0x75, 0x6d, 0x5f, 0x61, 0x75, 0x74,
	0x68, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x1a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x03, 0x0a, 0x16, 0x52, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x64, 0x75, 0x6d, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x5e, 0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x5f, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x50,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x54, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x52, 0x18, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x42, 0x0a, 0x14, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x52, 0x13, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x39, 0x0a, 0x11, 0x6f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x48, 0x61, 0x73,
	0x68, 0x52, 0x10, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x49, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x5f,
	0x77, 0x68, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x65, 0x72, 0x57, 0x68, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x11, 0x70, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x65, 0x72, 0x57, 0x68, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x33,
	0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x48, 0x61, 0x73, 0x68, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_referendum_auth_contract_proto_rawDescOnce sync.Once
	file_referendum_auth_contract_proto_rawDescData = file_referendum_auth_contract_proto_rawDesc
)

func file_referendum_auth_contract_proto_rawDescGZIP() []byte {
	file_referendum_auth_contract_proto_rawDescOnce.Do(func() {
		file_referendum_auth_contract_proto_rawDescData = protoimpl.X.CompressGZIP(file_referendum_auth_contract_proto_rawDescData)
	})
	return file_referendum_auth_contract_proto_rawDescData
}

var file_referendum_auth_contract_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_referendum_auth_contract_proto_goTypes = []interface{}{
	(*ReferendumOrganization)(nil),   // 0: client.ReferendumOrganization
	(*ProposalReleaseThreshold)(nil), // 1: client.ProposalReleaseThreshold
	(*Address)(nil),                  // 2: client.Address
	(*Hash)(nil),                     // 3: client.Hash
	(*ProposerWhiteList)(nil),        // 4: client.ProposerWhiteList
}
var file_referendum_auth_contract_proto_depIdxs = []int32{
	1, // 0: client.ReferendumOrganization.proposal_release_threshold:type_name -> client.ProposalReleaseThreshold
	2, // 1: client.ReferendumOrganization.organization_address:type_name -> client.Address
	3, // 2: client.ReferendumOrganization.organization_hash:type_name -> client.Hash
	4, // 3: client.ReferendumOrganization.proposer_white_list:type_name -> client.ProposerWhiteList
	3, // 4: client.ReferendumOrganization.creation_token:type_name -> client.Hash
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_referendum_auth_contract_proto_init() }
func file_referendum_auth_contract_proto_init() {
	if File_referendum_auth_contract_proto != nil {
		return
	}
	file_client_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_referendum_auth_contract_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReferendumOrganization); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_referendum_auth_contract_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_referendum_auth_contract_proto_goTypes,
		DependencyIndexes: file_referendum_auth_contract_proto_depIdxs,
		MessageInfos:      file_referendum_auth_contract_proto_msgTypes,
	}.Build()
	File_referendum_auth_contract_proto = out.File
	file_referendum_auth_contract_proto_rawDesc = nil
	file_referendum_auth_contract_proto_goTypes = nil
	file_referendum_auth_contract_proto_depIdxs = nil
}
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: praefect.proto

package gitalypb

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

// SetReplicationFactorRequest sets the desired replication factor for a repository.
type SetReplicationFactorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// virtual_storage is the virtual storage the repository is located in
	VirtualStorage string `protobuf:"bytes,1,opt,name=virtual_storage,json=virtualStorage,proto3" json:"virtual_storage,omitempty"`
	// relative_path is the relative path of the repository
	RelativePath string `protobuf:"bytes,2,opt,name=relative_path,json=relativePath,proto3" json:"relative_path,omitempty"`
	// replication_factor is the desired replication factor. Replication must be equal or greater than 1.
	ReplicationFactor int32 `protobuf:"varint,3,opt,name=replication_factor,json=replicationFactor,proto3" json:"replication_factor,omitempty"`
}

func (x *SetReplicationFactorRequest) Reset() {
	*x = SetReplicationFactorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetReplicationFactorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetReplicationFactorRequest) ProtoMessage() {}

func (x *SetReplicationFactorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetReplicationFactorRequest.ProtoReflect.Descriptor instead.
func (*SetReplicationFactorRequest) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{0}
}

func (x *SetReplicationFactorRequest) GetVirtualStorage() string {
	if x != nil {
		return x.VirtualStorage
	}
	return ""
}

func (x *SetReplicationFactorRequest) GetRelativePath() string {
	if x != nil {
		return x.RelativePath
	}
	return ""
}

func (x *SetReplicationFactorRequest) GetReplicationFactor() int32 {
	if x != nil {
		return x.ReplicationFactor
	}
	return 0
}

// SetReplicationFactorResponse returns the assigned hosts after setting the desired replication factor.
type SetReplicationFactorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// storages are the storages assigned to host the repository.
	Storages []string `protobuf:"bytes,1,rep,name=storages,proto3" json:"storages,omitempty"`
}

func (x *SetReplicationFactorResponse) Reset() {
	*x = SetReplicationFactorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetReplicationFactorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetReplicationFactorResponse) ProtoMessage() {}

func (x *SetReplicationFactorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetReplicationFactorResponse.ProtoReflect.Descriptor instead.
func (*SetReplicationFactorResponse) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{1}
}

func (x *SetReplicationFactorResponse) GetStorages() []string {
	if x != nil {
		return x.Storages
	}
	return nil
}

type SetAuthoritativeStorageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VirtualStorage       string `protobuf:"bytes,1,opt,name=virtual_storage,json=virtualStorage,proto3" json:"virtual_storage,omitempty"`
	RelativePath         string `protobuf:"bytes,2,opt,name=relative_path,json=relativePath,proto3" json:"relative_path,omitempty"`
	AuthoritativeStorage string `protobuf:"bytes,3,opt,name=authoritative_storage,json=authoritativeStorage,proto3" json:"authoritative_storage,omitempty"`
}

func (x *SetAuthoritativeStorageRequest) Reset() {
	*x = SetAuthoritativeStorageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetAuthoritativeStorageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetAuthoritativeStorageRequest) ProtoMessage() {}

func (x *SetAuthoritativeStorageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetAuthoritativeStorageRequest.ProtoReflect.Descriptor instead.
func (*SetAuthoritativeStorageRequest) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{2}
}

func (x *SetAuthoritativeStorageRequest) GetVirtualStorage() string {
	if x != nil {
		return x.VirtualStorage
	}
	return ""
}

func (x *SetAuthoritativeStorageRequest) GetRelativePath() string {
	if x != nil {
		return x.RelativePath
	}
	return ""
}

func (x *SetAuthoritativeStorageRequest) GetAuthoritativeStorage() string {
	if x != nil {
		return x.AuthoritativeStorage
	}
	return ""
}

type SetAuthoritativeStorageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetAuthoritativeStorageResponse) Reset() {
	*x = SetAuthoritativeStorageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetAuthoritativeStorageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetAuthoritativeStorageResponse) ProtoMessage() {}

func (x *SetAuthoritativeStorageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetAuthoritativeStorageResponse.ProtoReflect.Descriptor instead.
func (*SetAuthoritativeStorageResponse) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{3}
}

type DatalossCheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VirtualStorage string `protobuf:"bytes,1,opt,name=virtual_storage,json=virtualStorage,proto3" json:"virtual_storage,omitempty"`
	// include_partially_unavailable indicates whether to include repositories which are available but
	// are unavailable on some assigned storages.
	IncludePartiallyReplicated bool `protobuf:"varint,2,opt,name=include_partially_replicated,json=includePartiallyReplicated,proto3" json:"include_partially_replicated,omitempty"`
}

func (x *DatalossCheckRequest) Reset() {
	*x = DatalossCheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatalossCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatalossCheckRequest) ProtoMessage() {}

func (x *DatalossCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatalossCheckRequest.ProtoReflect.Descriptor instead.
func (*DatalossCheckRequest) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{4}
}

func (x *DatalossCheckRequest) GetVirtualStorage() string {
	if x != nil {
		return x.VirtualStorage
	}
	return ""
}

func (x *DatalossCheckRequest) GetIncludePartiallyReplicated() bool {
	if x != nil {
		return x.IncludePartiallyReplicated
	}
	return false
}

type DatalossCheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// repositories with data loss
	Repositories []*DatalossCheckResponse_Repository `protobuf:"bytes,2,rep,name=repositories,proto3" json:"repositories,omitempty"`
}

func (x *DatalossCheckResponse) Reset() {
	*x = DatalossCheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatalossCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatalossCheckResponse) ProtoMessage() {}

func (x *DatalossCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatalossCheckResponse.ProtoReflect.Descriptor instead.
func (*DatalossCheckResponse) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{5}
}

func (x *DatalossCheckResponse) GetRepositories() []*DatalossCheckResponse_Repository {
	if x != nil {
		return x.Repositories
	}
	return nil
}

type RepositoryReplicasRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
}

func (x *RepositoryReplicasRequest) Reset() {
	*x = RepositoryReplicasRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepositoryReplicasRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepositoryReplicasRequest) ProtoMessage() {}

func (x *RepositoryReplicasRequest) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepositoryReplicasRequest.ProtoReflect.Descriptor instead.
func (*RepositoryReplicasRequest) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{6}
}

func (x *RepositoryReplicasRequest) GetRepository() *Repository {
	if x != nil {
		return x.Repository
	}
	return nil
}

type RepositoryReplicasResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Primary  *RepositoryReplicasResponse_RepositoryDetails   `protobuf:"bytes,1,opt,name=primary,proto3" json:"primary,omitempty"`
	Replicas []*RepositoryReplicasResponse_RepositoryDetails `protobuf:"bytes,2,rep,name=replicas,proto3" json:"replicas,omitempty"`
}

func (x *RepositoryReplicasResponse) Reset() {
	*x = RepositoryReplicasResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepositoryReplicasResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepositoryReplicasResponse) ProtoMessage() {}

func (x *RepositoryReplicasResponse) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepositoryReplicasResponse.ProtoReflect.Descriptor instead.
func (*RepositoryReplicasResponse) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{7}
}

func (x *RepositoryReplicasResponse) GetPrimary() *RepositoryReplicasResponse_RepositoryDetails {
	if x != nil {
		return x.Primary
	}
	return nil
}

func (x *RepositoryReplicasResponse) GetReplicas() []*RepositoryReplicasResponse_RepositoryDetails {
	if x != nil {
		return x.Replicas
	}
	return nil
}

type DatalossCheckResponse_Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// relative path of the repository with outdated replicas
	RelativePath string `protobuf:"bytes,1,opt,name=relative_path,json=relativePath,proto3" json:"relative_path,omitempty"`
	// storages on which the repository is outdated
	Storages []*DatalossCheckResponse_Repository_Storage `protobuf:"bytes,2,rep,name=storages,proto3" json:"storages,omitempty"`
	// unavailable indicates whether the repository is in unavailable.
	Unavailable bool `protobuf:"varint,3,opt,name=unavailable,proto3" json:"unavailable,omitempty"`
	// current primary storage of the repository
	Primary string `protobuf:"bytes,4,opt,name=primary,proto3" json:"primary,omitempty"`
}

func (x *DatalossCheckResponse_Repository) Reset() {
	*x = DatalossCheckResponse_Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatalossCheckResponse_Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatalossCheckResponse_Repository) ProtoMessage() {}

func (x *DatalossCheckResponse_Repository) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatalossCheckResponse_Repository.ProtoReflect.Descriptor instead.
func (*DatalossCheckResponse_Repository) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{5, 0}
}

func (x *DatalossCheckResponse_Repository) GetRelativePath() string {
	if x != nil {
		return x.RelativePath
	}
	return ""
}

func (x *DatalossCheckResponse_Repository) GetStorages() []*DatalossCheckResponse_Repository_Storage {
	if x != nil {
		return x.Storages
	}
	return nil
}

func (x *DatalossCheckResponse_Repository) GetUnavailable() bool {
	if x != nil {
		return x.Unavailable
	}
	return false
}

func (x *DatalossCheckResponse_Repository) GetPrimary() string {
	if x != nil {
		return x.Primary
	}
	return ""
}

type DatalossCheckResponse_Repository_Storage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name of the storage
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// behind_by indicates how many generations this storage is behind.
	BehindBy int64 `protobuf:"varint,2,opt,name=behind_by,json=behindBy,proto3" json:"behind_by,omitempty"`
	// assigned indicates whether the storage is assigned to host the repository.
	Assigned bool `protobuf:"varint,3,opt,name=assigned,proto3" json:"assigned,omitempty"`
	// healthy indicates whether the storage is considered healthy by the consensus of Praefect nodes.
	Healthy bool `protobuf:"varint,4,opt,name=healthy,proto3" json:"healthy,omitempty"`
	// valid_primary indicates whether the storage is ready to act as the primary if necessary.
	ValidPrimary bool `protobuf:"varint,5,opt,name=valid_primary,json=validPrimary,proto3" json:"valid_primary,omitempty"`
}

func (x *DatalossCheckResponse_Repository_Storage) Reset() {
	*x = DatalossCheckResponse_Repository_Storage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatalossCheckResponse_Repository_Storage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatalossCheckResponse_Repository_Storage) ProtoMessage() {}

func (x *DatalossCheckResponse_Repository_Storage) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatalossCheckResponse_Repository_Storage.ProtoReflect.Descriptor instead.
func (*DatalossCheckResponse_Repository_Storage) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{5, 0, 0}
}

func (x *DatalossCheckResponse_Repository_Storage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DatalossCheckResponse_Repository_Storage) GetBehindBy() int64 {
	if x != nil {
		return x.BehindBy
	}
	return 0
}

func (x *DatalossCheckResponse_Repository_Storage) GetAssigned() bool {
	if x != nil {
		return x.Assigned
	}
	return false
}

func (x *DatalossCheckResponse_Repository_Storage) GetHealthy() bool {
	if x != nil {
		return x.Healthy
	}
	return false
}

func (x *DatalossCheckResponse_Repository_Storage) GetValidPrimary() bool {
	if x != nil {
		return x.ValidPrimary
	}
	return false
}

type RepositoryReplicasResponse_RepositoryDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Checksum   string      `protobuf:"bytes,2,opt,name=checksum,proto3" json:"checksum,omitempty"`
}

func (x *RepositoryReplicasResponse_RepositoryDetails) Reset() {
	*x = RepositoryReplicasResponse_RepositoryDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_praefect_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepositoryReplicasResponse_RepositoryDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepositoryReplicasResponse_RepositoryDetails) ProtoMessage() {}

func (x *RepositoryReplicasResponse_RepositoryDetails) ProtoReflect() protoreflect.Message {
	mi := &file_praefect_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepositoryReplicasResponse_RepositoryDetails.ProtoReflect.Descriptor instead.
func (*RepositoryReplicasResponse_RepositoryDetails) Descriptor() ([]byte, []int) {
	return file_praefect_proto_rawDescGZIP(), []int{7, 0}
}

func (x *RepositoryReplicasResponse_RepositoryDetails) GetRepository() *Repository {
	if x != nil {
		return x.Repository
	}
	return nil
}

func (x *RepositoryReplicasResponse_RepositoryDetails) GetChecksum() string {
	if x != nil {
		return x.Checksum
	}
	return ""
}

var File_praefect_proto protoreflect.FileDescriptor

var file_praefect_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x61, 0x65, 0x66, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x1a, 0x0a, 0x6c, 0x69, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x1b, 0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x76, 0x69, 0x72,
	0x74, 0x75, 0x61, 0x6c, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x50, 0x61, 0x74, 0x68,
	0x12, 0x2d, 0x0a, 0x12, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x72, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x22,
	0x3a, 0x0a, 0x1c, 0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x22, 0xa3, 0x01, 0x0a, 0x1e,
	0x53, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27,
	0x0a, 0x0f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x33, 0x0a, 0x15,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x22, 0x21, 0x0a, 0x1f, 0x53, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x81, 0x01, 0x0a, 0x14, 0x44, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x73,
	0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x0f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x40, 0x0a, 0x1c, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64,
	0x65, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x6c, 0x79, 0x5f, 0x72, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x1a, 0x69, 0x6e,
	0x63, 0x6c, 0x75, 0x64, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x6c, 0x79, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x22, 0xbb, 0x03, 0x0a, 0x15, 0x44, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x73, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c,
	0x79, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x1a, 0xd3, 0x02, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x23, 0x0a, 0x0d, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x4c, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e,
	0x44, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x75, 0x6e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x6e, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x1a, 0x95,
	0x01, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x62, 0x65, 0x68, 0x69, 0x6e, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x62, 0x65, 0x68, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x79, 0x12, 0x23, 0x0a, 0x0d, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x6d, 0x61,
	0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x50,
	0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x4f, 0x0a, 0x19, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x72, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x22, 0xa3, 0x02, 0x0a, 0x1a, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x70,
	0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x50, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c,
	0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x08,
	0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x1a, 0x63, 0x0a, 0x11, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x32, 0x0a,
	0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x32, 0x95, 0x03,
	0x0a, 0x13, 0x50, 0x72, 0x61, 0x65, 0x66, 0x65, 0x63, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x12, 0x21, 0x2e, 0x67, 0x69,
	0x74, 0x61, 0x6c, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52,
	0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0d, 0x44, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x73, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x12, 0x1c, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x73, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x73, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x6a, 0x0a, 0x17, 0x53, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x26, 0x2e, 0x67, 0x69,
	0x74, 0x61, 0x6c, 0x79, 0x2e, 0x53, 0x65, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x53, 0x65, 0x74,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x14,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x12, 0x23, 0x2e, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x2e, 0x53, 0x65,
	0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x69, 0x74, 0x61,
	0x6c, 0x79, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a,
	0x04, 0xf0, 0x97, 0x28, 0x01, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x67,
	0x69, 0x74, 0x61, 0x6c, 0x79, 0x2f, 0x76, 0x31, 0x34, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x79, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_praefect_proto_rawDescOnce sync.Once
	file_praefect_proto_rawDescData = file_praefect_proto_rawDesc
)

func file_praefect_proto_rawDescGZIP() []byte {
	file_praefect_proto_rawDescOnce.Do(func() {
		file_praefect_proto_rawDescData = protoimpl.X.CompressGZIP(file_praefect_proto_rawDescData)
	})
	return file_praefect_proto_rawDescData
}

var file_praefect_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_praefect_proto_goTypes = []interface{}{
	(*SetReplicationFactorRequest)(nil),                  // 0: gitaly.SetReplicationFactorRequest
	(*SetReplicationFactorResponse)(nil),                 // 1: gitaly.SetReplicationFactorResponse
	(*SetAuthoritativeStorageRequest)(nil),               // 2: gitaly.SetAuthoritativeStorageRequest
	(*SetAuthoritativeStorageResponse)(nil),              // 3: gitaly.SetAuthoritativeStorageResponse
	(*DatalossCheckRequest)(nil),                         // 4: gitaly.DatalossCheckRequest
	(*DatalossCheckResponse)(nil),                        // 5: gitaly.DatalossCheckResponse
	(*RepositoryReplicasRequest)(nil),                    // 6: gitaly.RepositoryReplicasRequest
	(*RepositoryReplicasResponse)(nil),                   // 7: gitaly.RepositoryReplicasResponse
	(*DatalossCheckResponse_Repository)(nil),             // 8: gitaly.DatalossCheckResponse.Repository
	(*DatalossCheckResponse_Repository_Storage)(nil),     // 9: gitaly.DatalossCheckResponse.Repository.Storage
	(*RepositoryReplicasResponse_RepositoryDetails)(nil), // 10: gitaly.RepositoryReplicasResponse.RepositoryDetails
	(*Repository)(nil),                                   // 11: gitaly.Repository
}
var file_praefect_proto_depIdxs = []int32{
	8,  // 0: gitaly.DatalossCheckResponse.repositories:type_name -> gitaly.DatalossCheckResponse.Repository
	11, // 1: gitaly.RepositoryReplicasRequest.repository:type_name -> gitaly.Repository
	10, // 2: gitaly.RepositoryReplicasResponse.primary:type_name -> gitaly.RepositoryReplicasResponse.RepositoryDetails
	10, // 3: gitaly.RepositoryReplicasResponse.replicas:type_name -> gitaly.RepositoryReplicasResponse.RepositoryDetails
	9,  // 4: gitaly.DatalossCheckResponse.Repository.storages:type_name -> gitaly.DatalossCheckResponse.Repository.Storage
	11, // 5: gitaly.RepositoryReplicasResponse.RepositoryDetails.repository:type_name -> gitaly.Repository
	6,  // 6: gitaly.PraefectInfoService.RepositoryReplicas:input_type -> gitaly.RepositoryReplicasRequest
	4,  // 7: gitaly.PraefectInfoService.DatalossCheck:input_type -> gitaly.DatalossCheckRequest
	2,  // 8: gitaly.PraefectInfoService.SetAuthoritativeStorage:input_type -> gitaly.SetAuthoritativeStorageRequest
	0,  // 9: gitaly.PraefectInfoService.SetReplicationFactor:input_type -> gitaly.SetReplicationFactorRequest
	7,  // 10: gitaly.PraefectInfoService.RepositoryReplicas:output_type -> gitaly.RepositoryReplicasResponse
	5,  // 11: gitaly.PraefectInfoService.DatalossCheck:output_type -> gitaly.DatalossCheckResponse
	3,  // 12: gitaly.PraefectInfoService.SetAuthoritativeStorage:output_type -> gitaly.SetAuthoritativeStorageResponse
	1,  // 13: gitaly.PraefectInfoService.SetReplicationFactor:output_type -> gitaly.SetReplicationFactorResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_praefect_proto_init() }
func file_praefect_proto_init() {
	if File_praefect_proto != nil {
		return
	}
	file_lint_proto_init()
	file_shared_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_praefect_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetReplicationFactorRequest); i {
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
		file_praefect_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetReplicationFactorResponse); i {
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
		file_praefect_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetAuthoritativeStorageRequest); i {
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
		file_praefect_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetAuthoritativeStorageResponse); i {
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
		file_praefect_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatalossCheckRequest); i {
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
		file_praefect_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatalossCheckResponse); i {
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
		file_praefect_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepositoryReplicasRequest); i {
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
		file_praefect_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepositoryReplicasResponse); i {
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
		file_praefect_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatalossCheckResponse_Repository); i {
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
		file_praefect_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatalossCheckResponse_Repository_Storage); i {
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
		file_praefect_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepositoryReplicasResponse_RepositoryDetails); i {
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
			RawDescriptor: file_praefect_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_praefect_proto_goTypes,
		DependencyIndexes: file_praefect_proto_depIdxs,
		MessageInfos:      file_praefect_proto_msgTypes,
	}.Build()
	File_praefect_proto = out.File
	file_praefect_proto_rawDesc = nil
	file_praefect_proto_goTypes = nil
	file_praefect_proto_depIdxs = nil
}

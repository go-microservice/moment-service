// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.3
// source: api/moment/v1/like.proto

package v1

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

type Like struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId    int64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ObjType   int32 `protobuf:"varint,3,opt,name=obj_type,json=objType,proto3" json:"obj_type,omitempty"`
	ObjId     int64 `protobuf:"varint,4,opt,name=obj_id,json=objId,proto3" json:"obj_id,omitempty"`
	Status    int32 `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt int64 `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt int64 `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Like) Reset() {
	*x = Like{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Like) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Like) ProtoMessage() {}

func (x *Like) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Like.ProtoReflect.Descriptor instead.
func (*Like) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{0}
}

func (x *Like) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Like) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Like) GetObjType() int32 {
	if x != nil {
		return x.ObjType
	}
	return 0
}

func (x *Like) GetObjId() int64 {
	if x != nil {
		return x.ObjId
	}
	return 0
}

func (x *Like) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Like) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Like) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type CreateLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ObjType int32 `protobuf:"varint,2,opt,name=obj_type,json=objType,proto3" json:"obj_type,omitempty"`
	ObjId   int64 `protobuf:"varint,3,opt,name=obj_id,json=objId,proto3" json:"obj_id,omitempty"`
}

func (x *CreateLikeRequest) Reset() {
	*x = CreateLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLikeRequest) ProtoMessage() {}

func (x *CreateLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLikeRequest.ProtoReflect.Descriptor instead.
func (*CreateLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{1}
}

func (x *CreateLikeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateLikeRequest) GetObjType() int32 {
	if x != nil {
		return x.ObjType
	}
	return 0
}

func (x *CreateLikeRequest) GetObjId() int64 {
	if x != nil {
		return x.ObjId
	}
	return 0
}

type CreateLikeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateLikeReply) Reset() {
	*x = CreateLikeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLikeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLikeReply) ProtoMessage() {}

func (x *CreateLikeReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLikeReply.ProtoReflect.Descriptor instead.
func (*CreateLikeReply) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{2}
}

type UpdateLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateLikeRequest) Reset() {
	*x = UpdateLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLikeRequest) ProtoMessage() {}

func (x *UpdateLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLikeRequest.ProtoReflect.Descriptor instead.
func (*UpdateLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{3}
}

type UpdateLikeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateLikeReply) Reset() {
	*x = UpdateLikeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLikeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLikeReply) ProtoMessage() {}

func (x *UpdateLikeReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLikeReply.ProtoReflect.Descriptor instead.
func (*UpdateLikeReply) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{4}
}

type DeleteLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ObjType int32 `protobuf:"varint,2,opt,name=obj_type,json=objType,proto3" json:"obj_type,omitempty"`
	ObjId   int64 `protobuf:"varint,3,opt,name=obj_id,json=objId,proto3" json:"obj_id,omitempty"`
}

func (x *DeleteLikeRequest) Reset() {
	*x = DeleteLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteLikeRequest) ProtoMessage() {}

func (x *DeleteLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteLikeRequest.ProtoReflect.Descriptor instead.
func (*DeleteLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteLikeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DeleteLikeRequest) GetObjType() int32 {
	if x != nil {
		return x.ObjType
	}
	return 0
}

func (x *DeleteLikeRequest) GetObjId() int64 {
	if x != nil {
		return x.ObjId
	}
	return 0
}

type DeleteLikeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteLikeReply) Reset() {
	*x = DeleteLikeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteLikeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteLikeReply) ProtoMessage() {}

func (x *DeleteLikeReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteLikeReply.ProtoReflect.Descriptor instead.
func (*DeleteLikeReply) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{6}
}

type GetLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ObjType int32 `protobuf:"varint,2,opt,name=obj_type,json=objType,proto3" json:"obj_type,omitempty"`
	ObjId   int64 `protobuf:"varint,3,opt,name=obj_id,json=objId,proto3" json:"obj_id,omitempty"`
}

func (x *GetLikeRequest) Reset() {
	*x = GetLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLikeRequest) ProtoMessage() {}

func (x *GetLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLikeRequest.ProtoReflect.Descriptor instead.
func (*GetLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{7}
}

func (x *GetLikeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetLikeRequest) GetObjType() int32 {
	if x != nil {
		return x.ObjType
	}
	return 0
}

func (x *GetLikeRequest) GetObjId() int64 {
	if x != nil {
		return x.ObjId
	}
	return 0
}

type GetLikeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Like *Like `protobuf:"bytes,1,opt,name=like,proto3" json:"like,omitempty"`
}

func (x *GetLikeReply) Reset() {
	*x = GetLikeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLikeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLikeReply) ProtoMessage() {}

func (x *GetLikeReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLikeReply.ProtoReflect.Descriptor instead.
func (*GetLikeReply) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{8}
}

func (x *GetLikeReply) GetLike() *Like {
	if x != nil {
		return x.Like
	}
	return nil
}

type ListPostLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PostId int64 `protobuf:"varint,2,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	LastId int64 `protobuf:"varint,3,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
	Limit  int32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListPostLikeRequest) Reset() {
	*x = ListPostLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPostLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPostLikeRequest) ProtoMessage() {}

func (x *ListPostLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPostLikeRequest.ProtoReflect.Descriptor instead.
func (*ListPostLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{9}
}

func (x *ListPostLikeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ListPostLikeRequest) GetPostId() int64 {
	if x != nil {
		return x.PostId
	}
	return 0
}

func (x *ListPostLikeRequest) GetLastId() int64 {
	if x != nil {
		return x.LastId
	}
	return 0
}

func (x *ListPostLikeRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ListCommentLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CommentId int64 `protobuf:"varint,2,opt,name=comment_id,json=commentId,proto3" json:"comment_id,omitempty"`
	LastId    int64 `protobuf:"varint,3,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
	Limit     int32 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListCommentLikeRequest) Reset() {
	*x = ListCommentLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCommentLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommentLikeRequest) ProtoMessage() {}

func (x *ListCommentLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommentLikeRequest.ProtoReflect.Descriptor instead.
func (*ListCommentLikeRequest) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{10}
}

func (x *ListCommentLikeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ListCommentLikeRequest) GetCommentId() int64 {
	if x != nil {
		return x.CommentId
	}
	return 0
}

func (x *ListCommentLikeRequest) GetLastId() int64 {
	if x != nil {
		return x.LastId
	}
	return 0
}

func (x *ListCommentLikeRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ListLikeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items   []*Like `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Count   int64   `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	HasMore bool    `protobuf:"varint,3,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
	LastId  int64   `protobuf:"varint,4,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
}

func (x *ListLikeReply) Reset() {
	*x = ListLikeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_moment_v1_like_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLikeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLikeReply) ProtoMessage() {}

func (x *ListLikeReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_moment_v1_like_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLikeReply.ProtoReflect.Descriptor instead.
func (*ListLikeReply) Descriptor() ([]byte, []int) {
	return file_api_moment_v1_like_proto_rawDescGZIP(), []int{11}
}

func (x *ListLikeReply) GetItems() []*Like {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListLikeReply) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ListLikeReply) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

func (x *ListLikeReply) GetLastId() int64 {
	if x != nil {
		return x.LastId
	}
	return 0
}

var File_api_moment_v1_like_proto protoreflect.FileDescriptor

var file_api_moment_v1_like_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x6c, 0x69, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x61, 0x70, 0x69, 0x2e,
	0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22, 0xb7, 0x01, 0x0a, 0x04, 0x4c, 0x69,
	0x6b, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x62, 0x6a, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6f,
	0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x62, 0x6a, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x22, 0x5e, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6b,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x62, 0x6a, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06,
	0x6f, 0x62, 0x6a, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x62,
	0x6a, 0x49, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6b,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x13, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x5e,
	0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x62, 0x6a, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x6f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x62, 0x6a, 0x49, 0x64, 0x22, 0x11,
	0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x5b, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x6f, 0x62, 0x6a, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x6f, 0x62, 0x6a, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x62, 0x6a, 0x49, 0x64, 0x22, 0x37,
	0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x27,
	0x0a, 0x04, 0x6c, 0x69, 0x6b, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x6b,
	0x65, 0x52, 0x04, 0x6c, 0x69, 0x6b, 0x65, 0x22, 0x76, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x6f, 0x73, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22,
	0x7f, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x22, 0x84, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f, 0x72, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x32, 0xee, 0x03, 0x0a, 0x0b, 0x4c, 0x69, 0x6b, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4e, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4e, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6b, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69,
	0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4c, 0x69,
	0x6b, 0x65, 0x12, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x50,
	0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x22,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x56, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c,
	0x69, 0x6b, 0x65, 0x12, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c,
	0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x69, 0x6b, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x4d, 0x0a, 0x0d, 0x61, 0x70, 0x69, 0x2e,
	0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x3a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x6f, 0x6d, 0x65, 0x6e,
	0x74, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_moment_v1_like_proto_rawDescOnce sync.Once
	file_api_moment_v1_like_proto_rawDescData = file_api_moment_v1_like_proto_rawDesc
)

func file_api_moment_v1_like_proto_rawDescGZIP() []byte {
	file_api_moment_v1_like_proto_rawDescOnce.Do(func() {
		file_api_moment_v1_like_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_moment_v1_like_proto_rawDescData)
	})
	return file_api_moment_v1_like_proto_rawDescData
}

var file_api_moment_v1_like_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_api_moment_v1_like_proto_goTypes = []interface{}{
	(*Like)(nil),                   // 0: api.moment.v1.Like
	(*CreateLikeRequest)(nil),      // 1: api.moment.v1.CreateLikeRequest
	(*CreateLikeReply)(nil),        // 2: api.moment.v1.CreateLikeReply
	(*UpdateLikeRequest)(nil),      // 3: api.moment.v1.UpdateLikeRequest
	(*UpdateLikeReply)(nil),        // 4: api.moment.v1.UpdateLikeReply
	(*DeleteLikeRequest)(nil),      // 5: api.moment.v1.DeleteLikeRequest
	(*DeleteLikeReply)(nil),        // 6: api.moment.v1.DeleteLikeReply
	(*GetLikeRequest)(nil),         // 7: api.moment.v1.GetLikeRequest
	(*GetLikeReply)(nil),           // 8: api.moment.v1.GetLikeReply
	(*ListPostLikeRequest)(nil),    // 9: api.moment.v1.ListPostLikeRequest
	(*ListCommentLikeRequest)(nil), // 10: api.moment.v1.ListCommentLikeRequest
	(*ListLikeReply)(nil),          // 11: api.moment.v1.ListLikeReply
}
var file_api_moment_v1_like_proto_depIdxs = []int32{
	0,  // 0: api.moment.v1.GetLikeReply.like:type_name -> api.moment.v1.Like
	0,  // 1: api.moment.v1.ListLikeReply.items:type_name -> api.moment.v1.Like
	1,  // 2: api.moment.v1.LikeService.CreateLike:input_type -> api.moment.v1.CreateLikeRequest
	3,  // 3: api.moment.v1.LikeService.UpdateLike:input_type -> api.moment.v1.UpdateLikeRequest
	5,  // 4: api.moment.v1.LikeService.DeleteLike:input_type -> api.moment.v1.DeleteLikeRequest
	7,  // 5: api.moment.v1.LikeService.GetLike:input_type -> api.moment.v1.GetLikeRequest
	9,  // 6: api.moment.v1.LikeService.ListPostLike:input_type -> api.moment.v1.ListPostLikeRequest
	10, // 7: api.moment.v1.LikeService.ListCommentLike:input_type -> api.moment.v1.ListCommentLikeRequest
	2,  // 8: api.moment.v1.LikeService.CreateLike:output_type -> api.moment.v1.CreateLikeReply
	4,  // 9: api.moment.v1.LikeService.UpdateLike:output_type -> api.moment.v1.UpdateLikeReply
	6,  // 10: api.moment.v1.LikeService.DeleteLike:output_type -> api.moment.v1.DeleteLikeReply
	8,  // 11: api.moment.v1.LikeService.GetLike:output_type -> api.moment.v1.GetLikeReply
	11, // 12: api.moment.v1.LikeService.ListPostLike:output_type -> api.moment.v1.ListLikeReply
	11, // 13: api.moment.v1.LikeService.ListCommentLike:output_type -> api.moment.v1.ListLikeReply
	8,  // [8:14] is the sub-list for method output_type
	2,  // [2:8] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_api_moment_v1_like_proto_init() }
func file_api_moment_v1_like_proto_init() {
	if File_api_moment_v1_like_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_moment_v1_like_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Like); i {
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
		file_api_moment_v1_like_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLikeReply); i {
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
		file_api_moment_v1_like_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateLikeReply); i {
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
		file_api_moment_v1_like_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteLikeReply); i {
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
		file_api_moment_v1_like_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLikeReply); i {
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
		file_api_moment_v1_like_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPostLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCommentLikeRequest); i {
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
		file_api_moment_v1_like_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLikeReply); i {
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
			RawDescriptor: file_api_moment_v1_like_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_moment_v1_like_proto_goTypes,
		DependencyIndexes: file_api_moment_v1_like_proto_depIdxs,
		MessageInfos:      file_api_moment_v1_like_proto_msgTypes,
	}.Build()
	File_api_moment_v1_like_proto = out.File
	file_api_moment_v1_like_proto_rawDesc = nil
	file_api_moment_v1_like_proto_goTypes = nil
	file_api_moment_v1_like_proto_depIdxs = nil
}

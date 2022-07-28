// Package model golang model from handler gRPC
package model

// Company model
type Company struct {
	ID   string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

// Price model
type Price struct {
	Company *Company `protobuf:"bytes,1,opt,name=Company,proto3" json:"Company,omitempty"`
	Ask     uint32   `protobuf:"varint,2,opt,name=Ask,proto3" json:"Ask,omitempty"`
	Bid     uint32   `protobuf:"varint,3,opt,name=Bid,proto3" json:"Bid,omitempty"`
	Time    string   `protobuf:"bytes,4,opt,name=Time,proto3" json:"Time,omitempty"`
}

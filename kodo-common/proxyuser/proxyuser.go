package proxyuser

import "github.com/qianjin/kodo-security/kodokey"

const (
	UserType_Users   = uint32(524316)
	UserType_Sudoers = uint32(193)
)

var (
	// dev users

	ProxyUser_Dev_general_storage_011 = ProxyUser{Uid: kodokey.Dev_UID_general_torage_011, Utype: UserType_Users}
	ProxyUser_Dev_general_storage_002 = ProxyUser{Uid: kodokey.Dev_UID_general_torage_002, Utype: UserType_Users}

	// dev sudoers

	ProxySudoer_Dev_admin = ProxyUser{Uid: kodokey.Dev_Uid_admin, Utype: UserType_Sudoers}
)

type ProxyUser struct {
	Uid       uint32 `json:"uid"`
	IamUid    uint32 `json:"iuid,omitempty"`
	Sudoer    uint32 `json:"suid,omitempty"`
	Utype     uint32 `json:"ut"`
	UtypeSu   uint32 `json:"sut,omitempty"`
	Devid     uint32 `json:"dev,omitempty"`
	Appid     uint32 `json:"app,omitempty"`
	Expires   uint32 `json:"e,omitempty"`
	AccessKey string `json:"ak,omitempty"`
	// extra properties
	SignType string `json:"signType,omitempty"`
}

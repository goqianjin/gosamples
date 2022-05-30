package bucketmodel

import "encoding/xml"

func NewCreateOption() *CreateOption {
	return &CreateOption{}
}

// ----

type CreateOption struct {
	Region string // z0, z1
}

func (opt *CreateOption) WithRegion(region string) *CreateOption {
	opt.Region = region
	return opt
}

// ----- helper -----

type BucketEntry struct {
	Id                string `json:"id" bson:"id"`
	Tbl               string `json:"tbl" bson:"tbl"`
	Uid               uint32 `json:"uid" bson:"uid"`
	Itbl              uint32 `json:"itbl" bson:"itbl"`
	PhyTbl            string `json:"phy" bson:"phy"`
	Ctime             int64  `json:"ctime" bson:"ctime"`
	DropTime          int64  `json:"drop" bson:"drop"` // !=0时，表示该条目被删除
	DropType          int64  `bson:"drop_type" json:"drop_type"`
	Region            string `json:"region" bson:"region"`
	Product           string `bson:"product,omitempty" json:"product,omitempty"`
	Zone              string `json:"zone" bson:"zone"`
	Global            bool   `json:"global" bson:"global"`
	Line              bool   `json:"line" bson:"line"`
	Versioning        bool   `bson:"versioning" json:"versioning"`
	EncryptionEnabled bool   `bson:"encryption_enabled" json:"encryption_enabled"`
	AllowNullKey      bool   `bson:"allow_nullkey" json:"allow_nullkey"`

	Ouid                uint32 `json:"ouid" bson:"ouid,omitempty"`
	Oitbl               uint32 `json:"oitbl" bson:"oitbl,omitempty"`
	Otbl                string `json:"otbl" bson:"otbl,omitempty"`
	Oid                 string `json:"oid" bson:"oid,omitempty"`
	Perm                uint32 `json:"perm" bson:"perm,omitempty"`
	NotAllowAccessByTbl bool   `bson:"not_allow_access_by_tbl,omitempty" json:"not_allow_access_by_tbl,omitempty"` //默认为false，准许以tbl访问s3api
	RsDbClusterId       string `bson:"rs_db_cluster_id,omitempty" json:"rs_db_cluster_id,omitempty"`

	Val        string       `json:"val" bson:"val,omitempty"`                 // bucket 被删除时，(uc_)val 会被保留
	DomainInfo []DomainInfo `json:"domain_info" bson:"domain_info,omitempty"` // bucket 被删除时，domain_info 会被清除
	Tags       []Tag        `json:"-" bson:"tags,omitempty"`                  //不需要输出json

	ObjectLock             Config   `json:"object_lock,omitempty" bson:"object_lock,omitempty"`
	SysTags                []string `json:"systags,omitempty" bson:"systags,omitempty"`
	MultiregionEnabled     bool     `json:"multiregion_enabled,omitempty" bson:"multiregion_enabled,omitempty"`
	MultiregionEverEnabled bool     `json:"multiregion_ever_enabled,omitempty" bson:"multiregion_ever_enabled,omitempty"`
}

type DomainInfo struct { //数据库字段
	Domain     string    `json:"domain" bson:"domain"`
	Refresh    bool      `json:"refresh" bson:"refresh"`
	Global     bool      `json:"global" bson:"global"`
	Ctime      int64     `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Utime      int64     `json:"utime,omitempty" bson:"utime,omitempty"`
	AntiLeech  AntiLeech `json:"antileech,omitempty" bson:"antileech,omitempty"`
	DomainType int       `json:"domaintype,omitempty" bson:"domaintype,omitempty"`
	ApiScope   int       `json:"apiscope,omitempty" bson:"apiscope,omitempty"`
}

type Tag struct {
	Key   string `json:"Key" bson:"k" xml:"Key"`
	Value string `json:"Value" bson:"v" xml:"Value"`
}

type Config struct {
	XMLNS             string    `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty" bson:"xmlns,omitempty"`
	XMLName           *xml.Name `xml:"ObjectLockConfiguration" json:"xml_name,omitempty" bson:"xml_name,omitempty"`
	ObjectLockEnabled string    `xml:"ObjectLockEnabled" json:"object_lock_enabled,omitempty" bson:"object_lock_enabled,omitempty"`
	Rule              *Rule     `xml:"Rule,omitempty" json:"rule,omitempty" bson:"rule,omitempty"`
}

type Rule struct {
	DefaultRetention *DefaultRetention `xml:"DefaultRetention" json:"default_retention" bson:"default_retention"`
}

type DefaultRetention struct {
	XMLName *xml.Name `xml:"DefaultRetention" json:"xmlns,omitempty" bson:"xmlns,omitempty"`
	Mode    string    `xml:"Mode" json:"mode" bson:"mode"`
	Days    *int64    `xml:"Days,omitempty" json:"days,omitempty" bson:"days,omitempty"`
	Years   *int64    `xml:"Years,omitempty" json:"years,omitempty" bson:"years,omitempty"`
}

type AntiLeech struct {
	ReferWhiteList []string `json:"refer_wl,omitempty" bson:"refer_wl"`
	ReferBlackList []string `json:"refer_bl,omitempty" bson:"refer_bl"`
	ReferNoRefer   bool     `json:"no_refer" bson:"no_refer"`
	AntiLeechMode  int      `json:"anti_leech_mode" bson:"anti_leech_mode"` // 0:off,1:wl,2:bl
	AntiLeechUsed  bool     `json:"anti_leech_used" bson:"anti_leech_used"` // 表示是否设置过,只要设置过了就应该一直为true
	SourceEnabled  bool     `json:"source_enabled" bson:"source_enabled"`
}

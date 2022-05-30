package rsmodel

type DeleteObjectReq struct {
	Bucket string
	Key    string
}

type DeleteObjectResp struct {
}

type GetObjectStatReq struct {
	Bucket    string
	Key       string
	NeedParts bool
}

type GetObjectStatResp struct {
	Hash     string  `json:"hash"`     // 文件的HASH值，使用hash值算法计算。
	Fsize    int64   `json:"fsize"`    // 资源内容的大小，单位：字节。
	PutTime  int64   `json:"putTime"`  // 上传时间，单位：100纳秒，其值去掉低七位即为Unix时间戳。
	MimeType string  `json:"mimeType"` // 资源的 MIME 类型。
	Type     int     `json:"type"`     // 资源的存储类型，0表示标准存储，1 表示低频存储，2 表示归档存储，3 表示深度归档存储
	Status   int     `json:"status"`   // 文件的存储状态，即禁用状态和启用状态间的的互相转换，0表示启用，1表示禁用，请参考：文件状态。
	Parts    []int64 `json:"parts"`    // 分拣分片信息，可能为空
}

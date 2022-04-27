package upmodel

type InitPartsReq struct {
	Bucket string
	Key    string
	//hasKey bool
}

type InitPartsResp struct {
	UploadID string `json:"uploadId"`
}

type UploadPartsReq struct {
	Bucket string
	Key    string
	//hasKey bool
	UploadId   string
	PartNumber int64
	Size       int
}

type UploadPartsResp struct {
	UploadPartsRet
}
type CompletePartsReq struct {
	Bucket string
	Key    string
	//hasKey bool
	UploadId string

	Extra *RputV2Extra
}

type CompletePartsResp struct {
	PutRet
}

// ----- helper ----

type UploadPartInfo struct {
	Etag       string `json:"etag"`
	PartNumber int64  `json:"partNumber"`
	partSize   int
	fileOffset int64
}

// RputV2Extra 表示分片上传 v2 额外可以指定的参数
type RputV2Extra struct {
	Recorder   Recorder          // 可选。上传进度记录
	Metadata   map[string]string // 可选。用户自定义文件 metadata 信息
	CustomVars map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost     string
	MimeType   string                                      // 可选。
	PartSize   int64                                       // 可选。每次上传的块大小
	TryTimes   int                                         // 可选。尝试次数
	Progresses []UploadPartInfo                            // 上传进度
	Notify     func(partNumber int64, ret *UploadPartsRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(partNumber int64, err error)
}

// UploadPartsRet 表示分片上传 v2 每个片上传完毕的返回值
type UploadPartsRet struct {
	Etag string `json:"etag"`
	MD5  string `json:"md5"`
}

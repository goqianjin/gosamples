package upmodel

type FormUploadReq struct {
	//ResourceKey string `json:resource_key`
	//CustomName  string `json:custom_name`
	//CostomValue string `json:custom_value`
	//CRC32       string `json:crc32`
	UploadToken string `json:upload_token`
	Accept      string `json:accept`
	//Data        io.Reader
	Key      string `json:key`
	FileName string `json:file_name`
	Extra    *PutExtra
	FileSize int64
}

type FormUploadResp struct {
	PutRet
}

// --------- helper -------

// PutExtra 为表单上传的额外可选项
type PutExtra struct {
	// 可选，用户自定义参数，必须以 "x:" 开头。若不以x:开头，则忽略。
	Params map[string]string

	UpHost string

	// 可选，当为 "" 时候，服务端自动判断。
	MimeType string

	// 上传事件：进度通知。这个事件的回调函数应该尽可能快地结束。
	OnProgress func(fsize, uploaded int64)
}

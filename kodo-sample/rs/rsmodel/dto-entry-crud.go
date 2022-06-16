package rsmodel

type GetEntryInfoReq struct {
	Itbl   uint32
	Bucket string
	Key    string
}

type GetEntryInfoResp struct {
	EntryInfo
}

type EntryInfo struct {
	EncodedFh                     string            `json:"fh"`
	Hash                          string            `json:"hash"`
	MimeType                      string            `json:"mimeType"`
	EndUser                       string            `json:"endUser"`
	Fsize                         int64             `json:"fsize"`
	PutTime                       int64             `json:"putTime"`
	Idc                           uint16            `json:"idc"`
	IP                            string            `json:"ip"`
	EncryptionKey                 string            `json:"encryptionKey"`
	XMeta                         map[string]string `json:"x-qn-meta"`
	Type                          uint32            `json:"type"`
	Version                       string            `json:"version"`
	LccDel                        int64             `json:"del"`
	LccToLine                     int64             `json:"line"`
	LccToArchive                  int64             `json:"ar"`
	LccToDeepArchive              int64             `json:"dar"`
	LccFreezeArchiveOrDeepArchive int64             `json:"frAr"`
}

package rsmodel

type EntryInfoResp struct {
	EncodedFh  string `	json:"fh"`
	EncodedFh2 string `json:"fh2,omitempty"`
	Hash       string `json:"hash"`
	MimeType   string `json:"mimeType"`
	EndUser    string `json:"endUser"`
	Fsize      int64  `json:"fsize"`
	PutTime    int64  `json:"putTime"`
	Idc        uint16 `json:"idc"`
}

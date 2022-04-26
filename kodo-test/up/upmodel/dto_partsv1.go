package upmodel

import "os"

type MkblkReq struct {
	BlockSize  int64
	BodyLength int64
}

type MkblkResp struct {
	BlkputRet
}
type BputReq struct {
	Ctx    string `json:"ctx"`
	Offset uint32 `json:"offset"`
}

type BputResp struct {
	BlkputRet
}
type MkfileReq struct {
	Fsize int64
	//hasKey bool
	Key string

	Extra *RputExtra
}

type MkfileResp struct {
	PutRet
}

// ----- helper ----

// BlkputRet 表示分片上传每个片上传完毕的返回值
type BlkputRet struct {
	Ctx        string `json:"ctx"`
	Checksum   string `json:"checksum"`
	Crc32      uint32 `json:"crc32"`
	Offset     uint32 `json:"offset"`
	Host       string `json:"host"`
	ExpiredAt  int64  `json:"expired_at"`
	chunkSize  int
	fileOffset int64
	blkIdx     int
}

// RputExtra 表示分片上传额外可以指定的参数
type RputExtra struct {
	Recorder   Recorder          // 可选。上传进度记录
	Params     map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost     string
	MimeType   string                                        // 可选。
	ChunkSize  int                                           // 可选。每次上传的Chunk大小
	TryTimes   int                                           // 可选。尝试次数
	Progresses []BlkputRet                                   // 可选。上传进度
	Notify     func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(blkIdx int, blkSize int, err error)
}

type Recorder interface {
	// 新建或更新文件分片上传的进度
	Set(key string, data []byte) error

	// 获取文件分片上传的进度信息
	Get(key string) ([]byte, error)

	// 删除文件分片上传的进度文件
	Delete(key string) error

	// 根据给定的文件信息生成持久化纪录的 key
	GenerateRecorderKey(keyInfos []string, sourceFileInfo os.FileInfo) string
}

// PutRet 为七牛标准的上传回复内容。
// 如果使用了上传回调或者自定义了returnBody，那么需要根据实际情况，自己自定义一个返回值结构体
type PutRet struct {
	Hash         string `json:"hash"`
	PersistentID string `json:"persistentId"`
	Key          string `json:"key"`
}

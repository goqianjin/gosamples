package sisyphusmodel

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CreateReq struct {
	Name             string   `json:"name"`
	Uid              uint32   `json:"uid"`
	SrcBkt           string   `json:"src_bkt"`
	DstBkt           string   `json:"dst_bkt"`
	Bkts             []string `json:"bkts"` //双向同步时使用
	IsSync           bool     `json:"is_sync"`
	Prefix           string   `json:"prefix"`
	ConflictStrategy int      `json:"conflict_strategy"`
}

type CreateResp struct {
	TaskId
}

type QueryReq struct {
	TaskId
}
type QueryResp struct {
	Task
}

type StopReq struct {
	TaskId
}

type StartReq struct {
	TaskId
}

type DeleteReq struct {
	TaskId
}

type CreateDualTaskResp struct {
	TaskId
}

type CreateDualTaskReq struct {
	CreateReq
}

type QueryDualTaskReq struct {
	TaskId
}
type QueryDualTaskResp struct {
	DualTask
}
type StopDualTaskReq struct {
	TaskId
}

type StartDualTaskReq struct {
	TaskId
}

type DeleteDualTaskReq struct {
	TaskId
}

// ---- entity ---

type TaskId struct {
	Id string `json:"id"`
}

type Task struct {
	Id              bson.ObjectId `json:"id" bson:"_id"`
	Name            string        `json:"name" bson:"name"`
	Source          BucketInfo    `json:"source" bson:"source"`
	Target          BucketInfo    `json:"target" bson:"target"`
	Option          Option        `json:"option" bson:"option"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
	DeletedAt       time.Time     `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	Status          string        `json:"status" bson:"status"`
	JobDone         int64         `json:"job_done" bson:"job_done"`
	JobFailed       int64         `json:"job_failed" bson:"job_failed"`
	JobSkipped      int64         `json:"job_skipped" bson:"job_skipped"`
	JobCount        int64         `json:"job_count" bson:"job_count"`
	FileDoneSize    int64         `json:"file_done_size" bson:"file_done_size"`
	LastJobPutTime  time.Time     `json:"last_job_put_time" bson:"last_job_put_time"`
	LastJobDoneTime time.Time     `json:"last_job_done_time" bson:"last_job_done_time"`
	Marker          string        `json:"marker" bson:"marker"`
	DisPatched      bool          `json:"dispatched" bson:"dispatched"`
	Type            int           `json:"type" bson:"type"`
}

type DualTask struct {
	Id      bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	Uid     uint32          `json:"uid,omitempty" bson:"uid,omitempty"`
	TaskIds []bson.ObjectId `json:"-" bson:"task_ids"`
	Tasks   []Task          `json:"tasks,omitempty" bson:"-"`
}

type BucketInfo struct {
	Uid    uint32      `json:"uid" bson:"uid"`
	Bucket string      `json:"bucket" bson:"bucket"`
	Zone   int         `json:"zone" bson:"zone"`
	Region string      `json:"region" bson:"region"`
	Args   *OQiniuArgs `json:"args,omitempty" bson:"args,omitempty"`
}

type OQiniuArgs struct {
	//若指定以下参数的全部指定则认为是跨账号体系的迁移任务，否则认为是跨区域同步任务
	Ak      string `json:"ak" bson:"ak"`             //可指定账号AK,用于跨账号体系的迁移
	Sk      string `json:"sk" bson:"sk"`             //可指定账号SK,用于跨账号体系的迁移
	BktAddr string `json:"bkt_addr" bson:"bkt_addr"` //可指定Bucket地址,用于跨账号体系的迁移
	RsAddr  string `json:"rs_addr" bson:"rs_addr"`   //可指定Rs地址,用于跨账号体系的迁移
	UpAddr  string `json:"up_addr" bson:"up_addr"`   //可指定Up地址,用于跨账号体系的迁移
}

type Option struct {
	IsSync           bool         `json:"is_sync" bson:"is_sync"`
	Prefix           string       `json:"prefix" bson:"prefix"`
	ConflictStrategy int          `json:"conflict_strategy" bson:"conflict_strategy"`
	SyncDelArgs      *SyncDelArgs `json:"sync_del_args,omitempty" bson:"sync_del_args,omitempty"`
}

type SyncDelArgs struct {
	SyncDel   bool  `json:"sync_del"`   //开启删除同步
	DelBefore int64 `json:"del_before"` //仅删除删除时间减去上传时间大于DelBefore阈值的文件，（秒）
}

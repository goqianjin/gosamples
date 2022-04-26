package model

type RsConfig struct {
	Hosts             []string `json:"hosts"`
	FailRetryInterval int64    `json:"fail_retry_interval"`
	TryTimes          uint32   `json:"try_times"`
	AccessKey         string   `json:"access_key"`
	SecretKey         string   `json:"secret_key"`
}


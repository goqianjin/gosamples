package main

import (
	"fmt"
	"strings"
	"time"
)

func generateTaskName() string {
	RFC3339Nano := "20060102150405.999999999"
	bucketPattern := "qianjin-task-%s"
	return fmt.Sprintf(bucketPattern, strings.ReplaceAll(time.Now().Format(RFC3339Nano), ".", ""))
}

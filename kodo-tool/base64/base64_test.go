package base64

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	entry := "userwork-hb" + ":" + "68d8365b9034f202ef06991d4f1493fb.mp3"
	result := base64.URLEncoding.EncodeToString([]byte(entry))
	fmt.Println(result)
}

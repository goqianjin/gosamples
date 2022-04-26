package rspub

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/qianjin/kodo-security/kodokey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/shenqianjin/rs_sample/ref/uri"
)

func TestStat_dev(t *testing.T) {
	bucket := "kodoimport-multirs-src"
	key := "test01.txt"
	body := ""
	//var Hosts = []string{"http://10.200.20.23:9433", "http://10.200.20.23:9433"}

	cli := client.NewClientWithHost("10.200.20.23:9433").
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "/stat/"+uri.Encode(bucket+":"+key)).
		RawQuery("").
		AddHeader("Host", "10.200.20.23:9433").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body).
		AddHeader("X-Qiniu-Date", strconv.FormatInt(time.Now().Add(20*time.Minute).Unix(), 10))
	resp := cli.Call(req)
	fmt.Printf("Resp: %+v\n", resp)
}

func TestStat_多用例_dev(t *testing.T) {
	bucket := "kodoimport-multirs-src"
	key := "test01.txt"
	body := ""
	fmt.Println("bucket: " + bucket + ", key: " + key)
	cli := client.NewClientWithHost("10.200.20.23:9433").
		WithAuthKey(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, "/stat/"+uri.Encode(bucket+":"+key)).
		RawQuery("").
		AddHeader("Host", "10.200.20.23:9433").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body)
	//http://10.200.20.25:10221/mkbucketv3/kodo-bucket-qiniuauth-expires-1j85ylvf/region/z0/nodomain/true

	type UserCase struct {
		subject        string
		headerXDate    string
		headerXExpires string
		expectedCode   int
		expectedErr    string
		queryXDate     string
		queryXExpires  string
	}
	/*time10MinsFromNow := strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)
	time20MinsFromNow := strconv.FormatInt(time.Now().Add(20*time.Minute).Unix(), 10)
	time10MinsBeforeNow := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	time20MinsBeforeNow := strconv.FormatInt(time.Now().Add(-20*time.Minute).Unix(), 10)*/
	var RFC3339TimeInSecondPattern = "20060102T150405Z"
	time10MinsFromNow := time.Now().Add(10 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time20MinsFromNow := time.Now().Add(20 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time10MinsBeforeNow := time.Now().Add(-10 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	time20MinsBeforeNow := time.Now().Add(-20 * time.Minute).UTC().Format(RFC3339TimeInSecondPattern)
	cases := []UserCase{
		{subject: "无X-Qiniu-Date & 无X-Qiniu-Expires",
			headerXDate: "", headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往前不超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time10MinsBeforeNow, headerXExpires: "", expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time20MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "含X-Qiniu-Date (无效:now往前超过15分钟) & 不包含X-Qiniu-Expires",
			headerXDate: time20MinsBeforeNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed"},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (有效)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (无效)",
			headerXDate: time20MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},

		{subject: "含X-Qiniu-Date (格式不合法) & 不包含X-Qiniu-Expires",
			headerXDate: "not_number", headerXExpires: "", expectedCode: 400, expectedErr: "invalid X-Qiniu-Date"},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (有效)",
			headerXDate: "not_number", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: ""},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (无效)",
			headerXDate: "not_number", headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired"},
		{subject: "不包含X-Qiniu-Date & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (有效:now往后不超过15分钟) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time10MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (无效:now往后超过15分钟) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: time20MinsBeforeNow, headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},
		{subject: "含X-Qiniu-Date (格式不合法) & 包含X-Qiniu-Expires (格式不合法)",
			headerXDate: "not_number", headerXExpires: "not_number", expectedCode: 400, expectedErr: "invalid X-Qiniu-Expires"},

		{subject: "Header含X-Qiniu-Date (有效:now往后不超过15分钟) & Query含X-Qiniu-Date (有效:now往后不超过15分钟)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效:now往后不超过15分钟) & Query含X-Qiniu-Date (无效:now往后超过15分钟)",
			headerXDate: time10MinsFromNow, headerXExpires: "", expectedCode: 403, expectedErr: "request time is too skewed", queryXDate: time20MinsFromNow},
		{subject: "Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Expires (有效)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXExpires: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Expires (无效)",
			headerXDate: "", headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXExpires: time10MinsBeforeNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Date (有效) & Query含X-Qiniu-Expires (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 200, expectedErr: "", queryXDate: time10MinsFromNow, queryXExpires: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (有效) & Query含X-Qiniu-Date (有效) & Query含X-Qiniu-Expires (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsFromNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow, queryXExpires: time10MinsBeforeNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (无效) & Query含X-Qiniu-Date (有效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time10MinsFromNow},
		{subject: "Header含X-Qiniu-Date (有效) & Header含X-Qiniu-Expires (无效) & Query含X-Qiniu-Date (无效)",
			headerXDate: time10MinsFromNow, headerXExpires: time10MinsBeforeNow, expectedCode: 403, expectedErr: "token is expired", queryXDate: time20MinsFromNow},
	}
	for _, c := range cases {
		copiedReq := req.DeepClone()
		if c.headerXDate != "" {
			copiedReq.SetHeader("X-Qiniu-Date", c.headerXDate)
		}
		if c.headerXExpires != "" {
			copiedReq.SetHeader("X-Qiniu-Expires", c.headerXExpires)
		}
		// query
		query := copiedReq.GetRawQuery()
		if c.queryXDate != "" {
			if query != "" {
				query = query + "&"
			}
			query = "X-Qiniu-Date=" + c.queryXDate
		}
		if c.queryXExpires != "" {
			if query != "" {
				query = query + "&"
			}
			query = query + "X-Qiniu-Expires=" + c.queryXExpires
		}
		if query != "" {
			copiedReq.RawQuery(query)
		}
		resp := cli.Call(copiedReq)
		passed := resp.StatusCode == c.expectedCode
		if c.expectedErr != "" {
			passed = strings.Contains(fmt.Sprintf("%v", resp.Body), c.expectedErr)
		}
		prefix := "【Failed】"
		if passed {
			prefix = "【Passed】"
		}
		errmsg := ""
		if resp.Err != nil {
			errmsg = resp.Err.Error()
		}
		fmt.Println(prefix + "[" + c.subject + "] responseCode: " + strconv.Itoa(resp.StatusCode) + ", body: " + fmt.Sprintf("%v", resp.Body) + ", err: " + errmsg)
	}
}

func TestTimeFormat(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	now := time.Now()
	fmt.Println("now format: " + now.Format(http.TimeFormat))
	fmt.Println("now format with UTC: " + now.UTC().Format(http.TimeFormat))
	//fmt.Printf("%v\n", now)
	pt1, _ := time.Parse(http.TimeFormat, "Fri, 14 Jan 2022 15:26:39 GMT")
	fmt.Printf("**** now format: %v\n", pt1)
	pt2, _ := time.Parse(http.TimeFormat, "Fri, 14 Jan 2022 07:26:39 GMT")
	fmt.Printf("**** now format with UTC: %v - now: %v, duration to now: %v\n", pt2, time.Now(), time.Now().Sub(pt2))
	pt3, _ := time.ParseInLocation(http.TimeFormat, "Fri, 14 Jan 2022 07:26:39 GMT", time.UTC)
	fmt.Printf("**** now format with UTC: %v\n", pt3)
	pt4, _ := time.ParseInLocation(http.TimeFormat, "Fri, 14 Jan 2022 07:26:39 GMT", loc)
	fmt.Printf("**** now format with UTC: %v --> %v\n", pt4, pt4.Format(http.TimeFormat))

	fmt.Println("---------------------------")
	fmt.Println("RFC3339 format: " + now.Format(time.RFC3339))
	fmt.Println("RFC3339 format (UTC): " + now.UTC().Format(time.RFC3339))
	rfcpt1, _ := time.Parse(time.RFC3339, "2022-01-14T18:08:04+08:00")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt1, rfcpt1.Format(http.TimeFormat), rfcpt1.Location())
	rfcpt2, _ := time.Parse(time.RFC3339, "2022-01-14T17:08:04+07:00")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt2, rfcpt2.Format(http.TimeFormat), rfcpt2.Location())
	rfcpt2, _ = time.Parse(time.RFC3339, "2022-01-14T10:08:04Z")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt2, rfcpt2.Format(http.TimeFormat), rfcpt2.Location())
	rfcpt2, _ = time.Parse("2006-01-02T15:04", "2022-01-14T10:08:04Z")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt2, rfcpt2.Format(http.TimeFormat), rfcpt2.Location())
	rfcpt2, _ = time.Parse("2006-01-02T15:04", "2022-01-14T10:08:04+08:00")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt2, rfcpt2.Format(http.TimeFormat), rfcpt2.Location())

}

func TestAcc(tt *testing.T) {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-09-20 15:39:06 +0800")
	fmt.Println(t)
	t, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", "2018-09-20 15:39:06 +0000 CST")
	fmt.Println(t)
	t, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "2018-09-20 15:39:06 +0800 CST")
	fmt.Println(t)
	t, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "2018-09-20 15:39:06 Z GMT")
	fmt.Println(t)
	t, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "2018-09-20 15:39:06 +0000 GMT")
	fmt.Println(t)

	var rfcpt2 time.Time
	rfcpt2, _ = time.Parse("2006-01-02T15:04:05Z07:00", "2022-01-14T10:08:04+01:00")
	fmt.Printf("**** now format with UTC: %v --> %v --> location: %v\n", rfcpt2, rfcpt2.Format(http.TimeFormat), rfcpt2.Location())

	pattern := "20060102T150405Z"
	now := time.Now()
	s := now.Format(pattern)
	fmt.Println(s)
	parsedNow, err := time.Parse(pattern, s)
	fmt.Println(parsedNow)
	fmt.Println(err)
	parsedNow, err = time.ParseInLocation(pattern, s, time.Local)
	fmt.Println(parsedNow)
	fmt.Println(err)

	utcnow := now.UTC()
	s = utcnow.Format(pattern)
	fmt.Println(s)
	parsedNow, err = time.Parse(pattern, s)
	fmt.Println(parsedNow)
	now = time.Unix(now.Unix(), 0)
	fmt.Println(parsedNow.Equal(now))
	fmt.Println(parsedNow.Before(time.Now()))
	fmt.Println(parsedNow.Local())
	fmt.Println(parsedNow.Local().Equal(now))
	fmt.Println(parsedNow.Local().Before(time.Now()))
	fmt.Println(err)
}

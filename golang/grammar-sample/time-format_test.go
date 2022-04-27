package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

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

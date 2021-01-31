package chouseisan

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type Candidate struct {
	year         int
	month        time.Month
	day          int
	hour, minute int
}

func NewCandidate(year int, month time.Month, day, hour, minute int) Candidate {
	return Candidate{
		year:   year,
		month:  month,
		day:    day,
		hour:   hour,
		minute: minute,
	}
}

func Create(ctx context.Context, name, comment string, candidates ...Candidate) (string, error) {
	values := url.Values{}
	values.Add("name", name)
	values.Add("comment", comment)
	values.Add("kouho", buildKouho(candidates))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://chouseisan.com/schedule/newEvent/create", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	location, err := resp.Location()
	if err != nil {
		return "", err
	}

	return "https://chouseisan.com/s?h=" + location.Query().Get("h"), nil
}

func buildKouho(candidates []Candidate) string {
	ss := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		ss = append(ss, candidate.String())
	}
	return strings.Join(ss, "\r\n")
}

func (candidate Candidate) String() string {
	var weekdayStr string
	switch zellersCongruence(candidate.year, candidate.month, candidate.day) {
	case 0:
		weekdayStr = "土"
	case 1:
		weekdayStr = "日"
	case 2:
		weekdayStr = "月"
	case 3:
		weekdayStr = "火"
	case 4:
		weekdayStr = "水"
	case 5:
		weekdayStr = "木"
	case 6:
		weekdayStr = "金"
	}

	return fmt.Sprintf("%d/%d(%s) %02d:%02d〜", candidate.month, candidate.day, weekdayStr, candidate.hour, candidate.minute)
}

func floor(x int) int {
	return int(math.Floor(float64(x)))
}

func zellersCongruence(year int, month time.Month, day int) int {
	if month == time.January || month == time.February {
		year--
		month += 12
	}
	c := year / 100
	y := year % 100

	weekday := day + (26 * (int(month) + 1) / 10) + y + (y / 4) + 5*c + (c / 4)

	return weekday % 7
}

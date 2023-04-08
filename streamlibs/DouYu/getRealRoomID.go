package DouYu

import (
	"errors"
	"fmt"
	"io"
	"regexp"
)

func (l *Link) getRealRoomID(rfid string) (rid string, err error) {
	ridRegex := regexp.MustCompile(`rid":(\d{1,8}),"vipId`)
	resp, err := l.client.Get(fmt.Sprintf("https://m.douyu.com/%s", rfid))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	l.res = string(body)
	rids := ridRegex.FindStringSubmatch(l.res)
	if len(rids) != 2 {
		return "", errors.New("房间号错误")
	}
	return rids[1], nil
}

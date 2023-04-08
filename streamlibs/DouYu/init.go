package DouYu

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Link struct {
	rid    string
	did    string
	t10    string
	t13    string
	res    string
	client *http.Client
}

func NewDouyuLink(rid string, proxy string) (*Link, error) {
	var (
		err          error
		proxyURL     *url.URL
		preErrorCode int64
	)
	dy := new(Link)
	dy.t10 = strconv.Itoa(int(time.Now().Unix()))
	dy.t13 = strconv.Itoa(int(time.Now().UnixMilli()))
	if len(proxy) > 0 {
		proxyURL, err = url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		dy.client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		dy.client = &http.Client{}
	}
	dy.rid, err = dy.getRealRoomID(rid)
	if err != nil {
		return nil, err
	}
	dy.did, err = dy.getDeviceID()
	if err != nil {
		return nil, err
	}
	preErrorCode, err = dy.getPreData()
	switch preErrorCode {
	case 102:
		return nil, errors.New("房间不存在")
	case 104:
		return nil, errors.New("房间未开播")
	}
	return dy, nil
}

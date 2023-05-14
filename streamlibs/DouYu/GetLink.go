package DouYu

import (
	"bytes"
	"encoding/json"
	"net/url"
)

func (l *Link) GetLink() (string, error) {
	data, err := l.getRateStream()
	if err != nil {
		return "", err
	}
	u, err := url.Parse(data.Get("data.url").String())
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	link := map[string]string{
		"m3u8": u.String(),
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	err = enc.Encode(link)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

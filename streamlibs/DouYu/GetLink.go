package DouYu

import (
	"bytes"
	"encoding/json"
)

func (l *Link) GetLink() (string, error) {
	data, err := l.getRateStream()
	if err != nil {
		return "", err
	}
	link := map[string]string{
		"m3u8": data.Get("data.url").String(),
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	err = enc.Encode(link)
	if err != nil {
		return "", err
	}
	return string(buf.String()), nil
}

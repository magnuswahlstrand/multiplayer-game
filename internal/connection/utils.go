package connection

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
)

func Encode(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	c, err := zip(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(c), nil
}

func Decode(data string, v interface{}) error {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	u, err := unzip(b)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(u, v); err != nil {
		return err
	}
	return nil
}

func zip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(in)
	if err != nil {
		return nil, nil
	}
	if err := w.Close(); err != nil {
		return nil, nil
	}
	return b.Bytes(), nil
}

func unzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return nil, nil
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return nil, nil
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil
	}
	return out, nil
}

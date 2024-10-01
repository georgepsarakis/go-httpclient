package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func DeserializeJSON(resp *http.Response, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("pointer required, got %T", target)
	}
	b, err := InterceptResponseBody(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func InterceptResponseBody(r *http.Response) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func MustInterceptResponseBody(r *http.Response) []byte {
	b, err := InterceptResponseBody(r)
	if err != nil {
		panic(err)
	}
	return b
}

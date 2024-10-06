package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// DeserializeJSON unmarshals the response body payload to the object referenced by the `target` pointer.
// If `target` is not a pointer, an error is returned.
// The body stream is restored as a NopCloser, so subsequent calls to `Body.Close()` will never fail.
// Note that the above behavior may have impact on memory requirements since memory must be reserved
// for the full lifecycle of the http.Response object.
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

// InterceptResponseBody will read the full contents of the http.Response.Body stream and release any resources
// associated with the Response object while allowing the Body stream to be accessed again.
// The Body stream is restored as a NopCloser, so subsequent calls to `Body.Close()` will never fail.
// Note that the above behavior may have impact on memory requirements since memory must be reserved
// for the full lifecycle of the http.Response object.
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

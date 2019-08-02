package test_result

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func DrainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func CopyRequest(r *http.Request) (*http.Request, error) {
	r2 := new(http.Request)
	*r2 = *r

	if r.URL != nil {
		r2URL := new(url.URL)
		*r2URL = *r.URL
		r2.URL = r2URL
	}

	if r.Body != nil {
		var err error
		r.Body, r2.Body, err = DrainBody(r.Body)
		if err != nil {
			return nil, err
		}
	}
	return r2, nil
}

package test_result

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func StartTestServerWithResponseMap(stubResponses map[string]*http.Response) (
	serverUrl string,
	client *http.Client,
	closeFn func(),
	incomingRequestsFn func() []*http.Request,
) {
	requests := make([]*http.Request, 0)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		copiedReq, _ := CopyRequest(req)
		requests = append(requests, copiedReq)

		if resp, exists := stubResponses[req.URL.String()]; exists {
			rw.WriteHeader(resp.StatusCode)
			rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			if resp.Body != nil {
				var body io.ReadCloser
				body, resp.Body, _ = DrainBody(resp.Body)
				_, _ = io.Copy(rw, body)
			}
			return
		}

		rw.WriteHeader(http.StatusNotFound)
	}))

	return server.URL,
		server.Client(),
		func() { server.Close() },
		func() []*http.Request { return requests }
}

func MockResponseOk(body string) *http.Response {
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(body))}
}

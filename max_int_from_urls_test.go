package test_result

import (
	"net/http"
	"testing"
)

func TestMaxIntFromUrls(t *testing.T) {
	/*
	   http://www.mocky.io/v2/5d1f3c7f3100001589ebebdc
	   http://www.mocky.io/v2/5d1f3c91310000bc86ebebde
	   http://www.mocky.io/v2/5d1f3c9f310000ce88ebebe0
	   http://www.mocky.io/v2/5d1f3cab310000bc86ebebe2
	*/

	serverUrl, client, closeServerFn, _ := StartTestServerWithResponseMap(map[string]*http.Response{
		"/v2/5d1f3c7f3100001589ebebdc": MockResponseOk("1"),
		"/v2/5d1f3c91310000bc86ebebde": MockResponseOk("1000"),
		"/v2/5d1f3c9f310000ce88ebebe0": MockResponseOk("123"),
		"/v2/5d1f3cab310000bc86ebebe2": MockResponseOk("10"),
	})
	defer closeServerFn()

	testCases := map[string]struct {
		list        []string
		expectedMax int
	}{
		"success case": {
			[]string{
				serverUrl + "/v2/5d1f3c7f3100001589ebebdc",
				serverUrl + "/v2/5d1f3c91310000bc86ebebde",
				serverUrl + "/v2/5d1f3c9f310000ce88ebebe0",
				serverUrl + "/v2/5d1f3cab310000bc86ebebe2",
			},
			1000,
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual, err := MaxIntFromUrls(client, testData.list)
			if err != nil {
				t.Errorf("got error %v", err)
			} else if actual != testData.expectedMax {
				t.Errorf("expect %d, but got %d", testData.expectedMax, actual)
			}
		})
	}
}

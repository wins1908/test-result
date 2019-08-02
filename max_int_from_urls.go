package test_result

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type result struct {
	onUrl     string
	gotErr    error
	gotNumber int
}

type multipleErrs map[string]error

func (e multipleErrs) Error() string { return fmt.Sprint(e) }

func MaxIntFromUrls(client *http.Client, urls []string) (int, error) {
	ch := make(chan result, len(urls))
	defer close(ch)

	for _, oneUrl := range urls {
		go func(url string) {
			rs := result{onUrl: url}
			resp, err := client.Get(url)
			if err != nil {
				rs.gotErr = err
				ch <- rs
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				rs.gotErr = fmt.Errorf("response status code is %d", resp.StatusCode)
				ch <- rs
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				rs.gotErr = err
				ch <- rs
				return
			}
			var number int
			if err := json.Unmarshal(body, &number); err != nil {
				rs.gotErr = err
				ch <- rs
				return
			}

			rs.gotNumber = number
			ch <- rs
		}(oneUrl)
	}

	max := 0
	gotErrs := make(map[string]error)
	for i := 0; i < len(urls); i++ {
		rs := <-ch
		if rs.gotErr != nil {
			gotErrs[rs.onUrl] = rs.gotErr
		} else if rs.gotNumber > max {
			max = rs.gotNumber
		}
	}

	if len(gotErrs) > 0 {
		return 0, multipleErrs(gotErrs)
	}
	return max, nil
}

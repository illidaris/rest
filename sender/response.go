package sender

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ParseResponse(resp *http.Response) ([]byte, error) {
	respbs := []byte{}
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	if resp.Body != nil {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		respbs = bs
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status is [%d]%s", resp.StatusCode, string(respbs))
	}
	return respbs, nil
}

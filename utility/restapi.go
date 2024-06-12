package utility

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpError struct {
	httpresultcode int
}

func (e *HttpError) Error() string {
	return "HTTP error " + strconv.Itoa(e.httpresultcode)
}
func CallRestAPIWithMethod(param interface{}, url, method string) (body string, err error, statuscode int) {
	jsonValue, _ := json.Marshal(param)

	// create request structure
	req, err := http.NewRequest(method, url, strings.NewReader(string(jsonValue)))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "application/json")
	req.Close = true // this is required to prevent too many files open

	// Create HTTP Connection
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(15) * time.Second,
	}

	// Now hit to destionation endpoint
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	buff := new(bytes.Buffer)
	buff.ReadFrom(res.Body)
	body = buff.String()
	statuscode = res.StatusCode

	if res.StatusCode != 200 {
		err = &HttpError{res.StatusCode}
	}

	return
}

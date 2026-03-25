package utils

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type RestyRequest struct {
	Method  string      `json:"method"`
	Url     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body"`
}

type RestyResponse struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
	ReceivedAt time.Time   `json:"received_at"`
}

type RestyLog struct {
	Request  RestyRequest  `json:"request"`
	Response RestyResponse `json:"response"`
}

func GetRestyLog(resp *resty.Response) RestyLog {
	reqHeaders := map[string][]string(resp.Request.Header)
	delete(reqHeaders, "User-Agent")

	return RestyLog{
		Request: RestyRequest{
			Method:  resp.Request.Method,
			Url:     resp.Request.URL,
			Headers: reqHeaders,
			Body:    resp.Request.Body,
		},
		Response: RestyResponse{
			StatusCode: resp.StatusCode(),
			Status:     resp.Status(),
			Headers:    resp.Header(),
			Body:       resp.String(),
			ReceivedAt: resp.ReceivedAt(),
		},
	}
}

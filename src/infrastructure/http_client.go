package infrastructure

import "net/http"

type HttpClient interface {
	Get(url string) (response *http.Response, err error)
}

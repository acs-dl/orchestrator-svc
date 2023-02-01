package requests

import (
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetRequestsRequest struct {
	FromUserId *int64 `filter:"fromUserId"`
	ToUserId   *int64 `filter:"toUserId"`
}

func NewGetRequestsRequest(r *http.Request) (GetRequestsRequest, error) {
	var request GetRequestsRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}

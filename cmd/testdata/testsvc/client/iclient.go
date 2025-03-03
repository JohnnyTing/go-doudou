package client

import (
	"context"
	"testsvc/vo"

	"github.com/go-resty/resty/v2"
)

type ITestsvcClient interface {
	PageUsers(ctx context.Context, _headers map[string]string, query vo.PageQuery) (_resp *resty.Response, code int, data vo.PageRet, err error)
}

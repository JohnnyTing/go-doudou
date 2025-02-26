package client

import (
	"context"
	"mime/multipart"
	"os"
	"testdata/vo"

	"github.com/go-resty/resty/v2"
	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
)

type IUsersvcClient interface {
	PageUsers(ctx context.Context, _headers map[string]string, query vo.PageQuery) (_resp *resty.Response, code int, data vo.PageRet, msg error)
	GetUser(ctx context.Context, _headers map[string]string, userId string, photo string) (_resp *resty.Response, code int, data string, msg error)
	SignUp(ctx context.Context, _headers map[string]string, username string, password int, actived bool, score []int) (_resp *resty.Response, code int, data string, msg error)
	UploadAvatar(ctx context.Context, _headers map[string]string, pf []v3.FileModel, ps string, pf2 v3.FileModel, pf3 *multipart.FileHeader, pf4 []*multipart.FileHeader) (_resp *resty.Response, ri int, ri2 interface{}, re error)
	DownloadAvatar(ctx context.Context, _headers map[string]string, userId interface{}, data []byte, userAttrs ...string) (_resp *resty.Response, rf *os.File, re error)
}

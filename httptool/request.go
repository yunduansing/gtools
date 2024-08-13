package httptool

import "context"

type IRequest interface {
	Get(ctx context.Context, url string, reqData map[string]interface{}, headers map[string]string) (res []byte, statusCode int, err error)
	PostJson(ctx context.Context, url string, reqData interface{}, headers map[string]string) (res []byte, statusCode int, err error)
}

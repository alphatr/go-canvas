package canvas

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
)

// GroupFunc 调用函数
type GroupFunc func() error

// Download 下载文件
func Download(url string) (io.ReadCloser, *http.Response, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	return res.Body, res, nil
}

// NewParallel 创建
func NewParallel(ctx context.Context, handles ...GroupFunc) error {
	errChan := make(chan error)
	doneChan := make(chan *struct{})

	for _, handle := range handles {
		currentHandle := handle

		go func() {
			if err := currentHandle(); err != nil {
				errChan <- err
			}

			doneChan <- nil
		}()
	}

	count := len(handles)

	for {
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}

		case err := <-errChan:
			{
				return err
			}

		case <-doneChan:
			{
				count--
				if count <= 0 {
					return nil
				}
			}
		}
	}
}

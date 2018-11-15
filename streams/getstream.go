package streams

import (
	"io"
	"net/http"
	"fmt"
)

type GetStream struct {
	reader io.Reader
}

func NewGetStream(objNameWithAddr string) (*GetStream, error) {
	resp, err := http.Get("http://" + objNameWithAddr)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("data server returned status code %d", resp.StatusCode)
	}

	return &GetStream{resp.Body}, nil
}

func (gs *GetStream) Read(p []byte) (n int, err error) {
	return gs.reader.Read(p)
}

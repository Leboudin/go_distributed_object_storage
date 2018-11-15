package streams

import (
	"io"
	"net/http"
	"fmt"
)

type PutStream struct {
	writer *io.PipeWriter
	errorC chan error
}

// Put objNameWithAddr in a goroutine and returns a PutStream struct
func NewPutStream(objNameWithAddr string) *PutStream {
	reader, writer := io.Pipe()
	errorC := make(chan error)

	go func() {
		req, _ := http.NewRequest("PUT", "http://"+objNameWithAddr, reader)
		client := http.Client{}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("data server returned status code %d", resp.StatusCode)
		}

		errorC <- err
	}()

	return &PutStream{writer, errorC}
}

// Implements the Write method
func (ps *PutStream) Write(data []byte) (n int, err error) {
	return ps.writer.Write(data)
}

// Implements the Close method, return any error during http request
func (ps *PutStream) Close() error {
	ps.writer.Close()
	return <-ps.errorC
}

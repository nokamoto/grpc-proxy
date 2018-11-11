package observe

import (
	"fmt"
	"github.com/nokamoto/grpc-proxy/yaml"
	"google.golang.org/grpc/codes"
	"os"
	"time"
)

type Log interface {
	Write(string, codes.Code, int, int, time.Duration) (int, error)
}

func NewLog(c yaml.Log) (Log, error) {
	if c.File == "/dev/stdout" {
		return &log{file: os.Stdout}, nil
	}
	if c.File == "/dev/stderr" {
		return &log{file: os.Stderr}, nil
	}

	file, err := os.OpenFile(c.File, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &log{file: file}, nil
}

type log struct {
	file *os.File
}

func (l *log) Write(method string, code codes.Code, req, res int, nanos time.Duration) (int, error) {
	return fmt.Fprintf(
		l.file,
		`{"ts":"%s","method":"%s","code":"%s","req":%d,"res":%d,"nanos":%d}%c`,
		time.Now(),
		method,
		code,
		req,
		res,
		nanos,
		'\n',
	)
}

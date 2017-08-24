package jsonstream

import (
	"bufio"
	"encoding/json"
	"io"
)

type Reader struct {
	r io.Reader
	s *bufio.Scanner
}

func NewReader(r io.Reader) *Reader {
	reader := &Reader{
		r: r,
		s: bufio.NewScanner(r),
	}

	return reader
}

func (r *Reader) Next() (Message, error) {
	for r.s.Scan() {
		if len(r.s.Bytes()) == 0 {
			continue
		}

		var msg partialMsg
		err := json.Unmarshal(r.s.Bytes(), &msg)
		if err != nil {
			return nil, err
		}

		switch msg.Type {
		case "log":
			var logMsg LogMsg
			err = json.Unmarshal(msg.Data, &logMsg)
			if err != nil {
				return nil, err
			}

			if logMsg.Attrs == nil {
				logMsg.Attrs = make(map[string]interface{})
			}

			return &logMsg, nil
		case "header":
			var headerMsg HeaderMsg
			err := json.Unmarshal(msg.Data, &headerMsg)
			if err != nil {
				return nil, err
			}

			if headerMsg.Headers == nil {
				headerMsg.Headers = make(map[string]interface{})
			}

			return &headerMsg, nil
		default:
			return nil, ErrUnknownType
		}
	}

	if r.s.Err() != nil {
		return nil, r.s.Err()
	}

	return nil, io.EOF
}

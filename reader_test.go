package jsonstream

import (
	"io"
	"os"
	"time"

	"github.com/aphistic/gomol"
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ReaderSuite struct{}

func (s *ReaderSuite) TestV1NoHeaders(t sweet.T) {
	f, err := os.Open("data/v1_noheader.json")
	Expect(err).To(BeNil())
	defer f.Close()

	rdr := NewReader(f)

	msg, err := rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(10, 0).UTC(),
		Msg:       "Message 1",
		Level:     gomol.LevelDebug,
		Attrs:     map[string]interface{}{},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(11, 0).UTC(),
		Msg:       "Message 2",
		Level:     gomol.LevelInfo,
		Attrs: map[string]interface{}{
			"attr1": float64(4321),
			"attr2": "val2",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(Equal(io.EOF))
	Expect(msg).To(BeNil())
}

func (s *ReaderSuite) TestV1OneHeader(t sweet.T) {
	f, err := os.Open("data/v1_oneheader.json")
	Expect(err).To(BeNil())
	defer f.Close()

	rdr := NewReader(f)

	msg, err := rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&HeaderMsg{
		Timestamp: time.Unix(5, 0).UTC(),
		Headers: map[string]interface{}{
			"header1": float64(1234),
			"header2": "headerVal2",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(10, 0).UTC(),
		Msg:       "Message 1",
		Level:     gomol.LevelDebug,
		Attrs:     map[string]interface{}{},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(11, 0).UTC(),
		Msg:       "Message 2",
		Level:     gomol.LevelInfo,
		Attrs: map[string]interface{}{
			"attr1": float64(4321),
			"attr2": "val2",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(Equal(io.EOF))
	Expect(msg).To(BeNil())
}

func (s *ReaderSuite) TestV1MultiHeaders(t sweet.T) {
	f, err := os.Open("data/v1_multiheader.json")
	Expect(err).To(BeNil())
	defer f.Close()

	rdr := NewReader(f)

	msg, err := rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&HeaderMsg{
		Timestamp: time.Unix(5, 0).UTC(),
		Headers: map[string]interface{}{
			"header1": float64(1234),
			"header2": "headerVal2",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(10, 0).UTC(),
		Msg:       "Message 1",
		Level:     gomol.LevelDebug,
		Attrs:     map[string]interface{}{},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&HeaderMsg{
		Timestamp: time.Unix(10, 0).UTC(),
		Headers: map[string]interface{}{
			"header1": "val1",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(BeNil())
	Expect(msg).To(Equal(&LogMsg{
		Timestamp: time.Unix(11, 0).UTC(),
		Msg:       "Message 2",
		Level:     gomol.LevelInfo,
		Attrs: map[string]interface{}{
			"attr1": float64(4321),
			"attr2": "val2",
		},
	}))

	msg, err = rdr.Next()
	Expect(err).To(Equal(io.EOF))
	Expect(msg).To(BeNil())
}

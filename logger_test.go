package jsonstream

import (
	"bytes"
	"time"

	"github.com/aphistic/gomol"
	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type LoggerSuite struct{}

func (s *LoggerSuite) TestUnsetHeaders(t sweet.T) {
	buf := bytes.NewBuffer([]byte{})

	cfg := NewConfig(buf)
	l, err := NewLogger(cfg)
	Expect(err).To(BeNil())
	Expect(l).ToNot(BeNil())

	err = l.InitLogger()
	Expect(err).To(BeNil())
	defer l.ShutdownLogger()

	err = l.Logm(
		time.Unix(10, 0).UTC(),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 1234,
			"attr2": "val2",
		},
		"Message 1",
	)
	Expect(err).To(BeNil())

	err = l.Flush()
	Expect(err).To(BeNil())

	Expect(buf.String()).To(Equal(`{` +
		`"type":"log",` +
		`"data":{` +
		`"time":"1970-01-01T00:00:10Z",` +
		`"msg":"Message 1",` +
		`"level":"debug",` +
		`"attrs":{"attr1":1234,"attr2":"val2"}` +
		`}}` + "\n",
	))
}

func (s *LoggerSuite) TestSetHeaders(t sweet.T) {
	buf := bytes.NewBuffer([]byte{})
	clock := glock.NewMockClock()
	clock.SetCurrent(time.Unix(5, 0).UTC())

	cfg := NewConfig(
		buf,
		WithHeaders(map[string]interface{}{
			"header1": 4321,
			"header2": "headerVal2",
		}),
		withClock(clock),
	)
	l, err := NewLogger(cfg)
	Expect(err).To(BeNil())
	Expect(l).ToNot(BeNil())

	err = l.InitLogger()
	Expect(err).To(BeNil())
	defer l.ShutdownLogger()

	err = l.Logm(
		time.Unix(10, 0).UTC(),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 1234,
			"attr2": "val2",
		},
		"Message 1",
	)
	Expect(err).To(BeNil())

	err = l.Flush()
	Expect(err).To(BeNil())

	Expect(buf.String()).To(Equal(`{` +
		`"type":"header",` +
		`"data":{` +
		`"time":"1970-01-01T00:00:05Z",` +
		`"headers":{"header1":4321,"header2":"headerVal2"}` +
		`}}` + "\n" +
		`{` +
		`"type":"log",` +
		`"data":{` +
		`"time":"1970-01-01T00:00:10Z",` +
		`"msg":"Message 1",` +
		`"level":"debug",` +
		`"attrs":{"attr1":1234,"attr2":"val2"}` +
		`}}` + "\n",
	))
}

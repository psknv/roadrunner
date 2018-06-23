package roadrunner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrBuffer_Write_Len(t *testing.T) {
	buf := newErrBuffer()

	buf.Write([]byte("hello"))
	assert.Equal(t, 5, buf.Len())
	assert.Equal(t, "hello", buf.String())
}

func TestErrBuffer_Write_Event(t *testing.T) {
	buf := newErrBuffer()

	tr := make(chan interface{})
	buf.Listen(func(event int, ctx interface{}) {
		assert.Equal(t, EventStderrOutput, event)
		assert.Equal(t, []byte("hello\n\n"), ctx)
		close(tr)
	})

	buf.Write([]byte("hello\n\n"))

	<-tr

	// messages are read
	assert.Equal(t, 0, buf.Len())
	assert.Equal(t, "", buf.String())
}

func TestErrBuffer_Write_Event_Separated(t *testing.T) {
	buf := newErrBuffer()

	tr := make(chan interface{})
	buf.Listen(func(event int, ctx interface{}) {
		assert.Equal(t, EventStderrOutput, event)
		assert.Equal(t, []byte("hello\n\n"), ctx)
		close(tr)
	})

	buf.Write([]byte("hel"))
	buf.Write([]byte("lo\n\n"))
	buf.Write([]byte("ending"))

	<-tr
	assert.Equal(t, 6, buf.Len())
	assert.Equal(t, "ending", buf.String())
}

func TestErrBuffer_Write_Remaining(t *testing.T) {
	buf := newErrBuffer()

	buf.Write([]byte("hel"))

	assert.Equal(t, 3, buf.Len())
	assert.Equal(t, "hel", buf.String())
}

//+build ignore-coverage

package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/radiofrance/kolibri/log"
)

func TestNew(t *testing.T) {
	f := New()
	assert.NotNil(t, f)
	assert.Equal(t, "", f.(*Fake).string)
}

func TestFake_Attributes(t *testing.T) {
	f := New()

	f = f.Named("test")
	assert.Equal(t, "/test", f.(*Fake).string)
	f = f.Named("test2")
	assert.Equal(t, "/test/test2", f.(*Fake).string)

	f = f.With(log.Bool("bool", true))
	assert.Equal(t, "/test/test2", f.(*Fake).string)
}

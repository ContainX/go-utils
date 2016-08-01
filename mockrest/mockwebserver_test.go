package mockrest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDynamicURLPresent(t *testing.T) {
	s := New()
	defer s.Stop()

	assert.Empty(t, s.URL, "Expecting empty URL")

	s.Start()
	assert.NotEmpty(t, s.URL, "Expecting URL to be defined")
}

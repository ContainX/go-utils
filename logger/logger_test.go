package logger

import (
	"testing"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	log := GetLogger("test")
	assert.Equal(t, logrus.InfoLevel, log.Level, "Expecting INFO level")
	SetLevel(logrus.FatalLevel, "test")
	assert.Equal(t, logrus.FatalLevel, log.Level, "Expecting FATAL level")
}

func TestSameInstance(t *testing.T) {
	log := GetLogger("a")
	log.Level = logrus.FatalLevel

	log2 := GetLogger("a")
	assert.Equal(t, log2.Level, logrus.FatalLevel, "Expecting FATAL level")
}

func TestCategoryPresentInLogger(t *testing.T) {
	log := GetLogger("abc")
	assert.Equal(t, "abc", log.category, "Expecting 'abc'")
}

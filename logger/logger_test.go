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

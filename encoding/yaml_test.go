package encoding

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"reflect"
)

func TestYAMLMarshal(t *testing.T) {

	encoder := newYAMLEncoder()
	for _, test := range personTests {
		str, err := encoder.Marshal(test)
		assert.NoError(t, err, "Error was found")
		for _, ss := range strings.Split(test.yaml, "\n") {
			assert.True(t, strings.Contains(str, ss), fmt.Sprintf("Was expecting '%s'", ss))
		}
	}
}

func TestYAMLUnMarshal(t *testing.T) {

	encoder := newYAMLEncoder()
	for _, test := range personTests {
		val := personStruct{}
		err := encoder.UnMarshalStr(test.yaml, &val)
		assert.NoError(t, err, "Error was found")

		// since we don't marshal/unmarshal this is needed for equality
		val.yaml = test.yaml
		val.json = test.json
		assert.True(t, reflect.DeepEqual(val, test))
	}
}


package encoding

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"reflect"
)

func TestJSONMarshal(t *testing.T) {

	encoder := newJSONEncoder()
	for _, test := range personTests {
		str, err := encoder.Marshal(test)
		assert.NoError(t, err, "Error was found")
		assert.True(t, strings.Contains(str, test.json), fmt.Sprintf("Was expecting '%s'", test.json))
	}
}

func TestJSONUnMarshal(t *testing.T) {

	encoder := newJSONEncoder()
	for _, test := range personTests {
		val := personStruct{}
		err := encoder.UnMarshalStr(test.json, &val)
		assert.NoError(t, err, "Error was found")

		// since we don't marshal/unmarshal this is needed for equality
		val.yaml = test.yaml
		val.json = test.json
		assert.True(t, reflect.DeepEqual(val, test))
	}
}
package httpclient

import (
	"testing"
	"github.com/ContainX/go-utils/mockrest"
	"github.com/stretchr/testify/assert"
)

type personStruct struct {
	Name   string   `json:"name,omitempty"`
	Age    int      `json:"age,omitempty"`
	Score  float64  `json:"score,omitempty"`
	Colors []string `json:"colors,omitempty"`
}

const (
	TestDataDir = "testdata/"
)

func TestGET(t *testing.T) {
	s := mockrest.StartNewWithFile(TestDataDir + "GET-response.json")
	defer s.Stop()

	person := &personStruct{}
	resp := Get(s.URL, person)

	assert.Nil(t, resp.Error, "Error response was not expected")
	assert.Equal(t, "John Doe", person.Name, "Expected name of John Doe")
	assert.Equal(t, 22, person.Age, "Expected age of 22")
}

func TestGET_404(t *testing.T) {
	s := mockrest.StartNewWithStatusCode(404)
	defer s.Stop()

	resp := Get(s.URL, nil)
	assert.Error(t, resp.Error, "Error was expected")
	assert.Equal(t, 404, resp.Status, "Status of 404 was expected")
}




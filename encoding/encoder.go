// encoding package for YAML and JSON encoding types
package encoding

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Encoder type represents either JSON or YAML
type EncoderType int

const (
	// JSON encoder type
	JSON EncoderType = 1 + iota
	// YAML encoder type
	YAML
)

var ErrorInvalidExtension = errors.New("File extension must be [.json | .yml | .yaml]")

type Encoder interface {

	// MarshalIndent is like Marshal but applies Indent to format the output.
	MarshalIndent(v interface{}) (string, error)

	// Marshal returns the JSON or YAML encoding of v.
	Marshal(v interface{}) (string, error)

	// Unmarshal parses the JSON-encoded or YAML-encoded data and stores the result in the value pointed to by v.
	UnMarshal(r io.Reader, v interface{}) error

	// UnMarshalStr is like UnMarshal but handles a raw string vs. a reader interface
	UnMarshalStr(data string, result interface{}) error
}

// NewEncoder creates an encoder based by the value pointed to by encoder
func NewEncoder(encoder EncoderType) (Encoder, error) {
	switch encoder {
	case JSON:
		return newJSONEncoder(), nil
	case YAML:
		return newYAMLEncoder(), nil
	default:
		panic(fmt.Errorf("Unsupported encoder type: %d", encoder))
	}
}

// NewEncoderFromFileExt creates an encoder based on the filename extension
func NewEncoderFromFileExt(filename string) (Encoder, error) {

	if et, err := EncoderTypeFromExt(filename); err != nil {
		return nil, err
	} else {
		return NewEncoder(et)
	}
}

// EncoderTypeFromExt simply returns the encoder type based on the filename extension
func EncoderTypeFromExt(filename string) (EncoderType, error) {
	switch filepath.Ext(filename) {
	case ".yml", ".yaml":
		return YAML, nil
	case ".json":
		return JSON, nil
	}
	return JSON, ErrorInvalidExtension

}

// ConvertFile converts from one encoding type of infile and writes to another encoding type as the outfile.
// dataType is the concrete type used for both infile and outfile
func ConvertFile(infile, outfile string, dataType interface{}) error {
	var fromEnc, toEnc Encoder
	var encErr error

	if fromEnc, encErr = NewEncoderFromFileExt(infile); encErr != nil {
		return encErr
	}

	if toEnc, encErr = NewEncoderFromFileExt(outfile); encErr != nil {
		return encErr
	}

	file, err := os.Open(infile)
	if err != nil {
		return err
	}

	if err := fromEnc.UnMarshal(file, dataType); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outfile), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if data, err := toEnc.MarshalIndent(dataType); err != nil {
		return err
	} else {
		f.WriteString(data)
	}
	return nil
}

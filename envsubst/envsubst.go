// envsubst provides substitution and expansion of environment variable place holders with a file/string
package envsubst

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"unicode"
	"github.com/ContainX/go-utils/logger"
)

var log = logger.Logger()

type runeReader interface {
	ReadRune() (rune, int, error)
}

type runeWriter interface {
	WriteRune(rune) (int, error)
}

type state int

const (
	initial state = iota
	readingVarName
	readingBracedVarName
)

type varNameTokenStatus int

const (
	complete varNameTokenStatus = iota
	incomplete
)

type undefinedVariableBehavior int

const (
	remove undefinedVariableBehavior = iota
	preserve
)

type envsubst struct {
	state             state
	buffer            bytes.Buffer
	target            runeWriter
	undefinedBehavior undefinedVariableBehavior
	resolver          func(string) string
}

func isVarNameCharacter(char rune, isFirstLetter bool) bool {
	if !isFirstLetter && unicode.IsDigit(char) {
		return true
	}
	return unicode.IsLetter(char) || char == '_'
}

func standaloneDollarString(varNameTokenStatus varNameTokenStatus, state state) string {
	switch {
	case state == readingVarName:
		return "$"
	case varNameTokenStatus == incomplete:
		return "${"
	}

	return "${}"
}

func writeString(s string, target runeWriter) error {
	for _, char := range s {
		if err := writeRune(char, target); err != nil {
			return err
		}
	}

	return nil
}

func writeRune(char rune, target runeWriter) error {
	_, err := target.WriteRune(char)
	return err
}

func substituteVariableReferences(source runeReader, target runeWriter, undefinedBehavior undefinedVariableBehavior, resolver func(string) string) error {
	et := envsubst{
		target:            target,
		undefinedBehavior: undefinedBehavior,
		resolver:          resolver,
	}

	for char, size, _ := source.ReadRune(); size != 0; char, size, _ = source.ReadRune() {
		if err := et.processRune(char); err != nil {
			return err
		}
	}

	return et.endOfInput()
}

func (et *envsubst) processRune(char rune) error {
	switch et.state {
	case initial:
		switch {
		case char == '$':
			et.state = readingVarName
		default:
			return writeRune(char, et.target)
		}
	case readingVarName:
		switch {
		//		case isVarNameCharacter(char, et.buffer.Len() == 0):
		//			return writeRune(char, &et.buffer)
		case char == '{' && et.buffer.Len() == 0:
			et.state = readingBracedVarName
		default:
			return et.flushBufferAndProcessNextRune(complete, char)
		}
	case readingBracedVarName:
		switch {
		case isVarNameCharacter(char, et.buffer.Len() == 0):
			return writeRune(char, &et.buffer)
		case char == '}':
			return et.flushBuffer(complete)
		default:
			return et.flushBufferAndProcessNextRune(incomplete, char)
		}
	}

	return nil
}

func (et *envsubst) endOfInput() error {
	if et.state != initial {
		return et.flushBuffer(incomplete)
	}

	return nil
}

func (et *envsubst) flushBufferAndProcessNextRune(bufferStatus varNameTokenStatus, nextChar rune) error {
	if err := et.flushBuffer(bufferStatus); err != nil {
		return err
	}

	return et.processRune(nextChar)
}

func (et *envsubst) flushBuffer(bufferStatus varNameTokenStatus) error {
	var err error

	switch {
	case et.buffer.Len() == 0:
		err = writeString(standaloneDollarString(bufferStatus, et.state), et.target)
	case et.state == readingBracedVarName && bufferStatus == incomplete:
		err = writeString("${"+et.buffer.String(), et.target)
	default:
		err = writeString(et.resolve(et.buffer.String()), et.target)
	}

	et.state = initial
	et.buffer.Reset()

	return err
}

func (et *envsubst) resolve(variableName string) string {
	resolvedValue := et.resolver(variableName)
	if len(resolvedValue) == 0 && et.undefinedBehavior == preserve {
		if et.state == readingBracedVarName {
			return "${" + variableName + "}"
		}
		return "$" + variableName
	}
	return resolvedValue
}

// Substitute will substitute ${TOKEN} tokens from the in reader and replaces the value
// using the resolver func.  The preserveUndef flag can be used to leave the original
// ${TOKEN} as is when scenarios arrise and the TOKEN value cannot be found against the
// resolver
func Substitute(in io.Reader, preserveUndef bool, resolver func(string) string) string {
	undefinedBehavior := remove
	if preserveUndef {
		undefinedBehavior = preserve
	}

	buf := new(bytes.Buffer)
	if err := substituteVariableReferences(bufio.NewReader(in), buf, undefinedBehavior, resolver); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

// SubstFileTokens provides feedback if any of the file tokens are missing within the specified
// in reader.
func SubstFileTokens(in io.Reader, params map[string]string) (parsed string, missing bool) {
	parsed = Substitute(in, true, func(s string) string {
		if params != nil && params[s] != "" {
			return params[s]
		}
		if os.Getenv(s) == "" {
			log.Warning("Cannot find a value for varible ${%s} in template", s)
			missing = true
		}
		return os.Getenv(s)
	})
	return parsed, missing
}

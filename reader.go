package pc

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

//---

// TextPosition is the value type that represents the position of a text.
type TextPosition struct {
	Line   int
	Column int
}

func (p TextPosition) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// IsValid checks the validity of the line and column.
func (p TextPosition) IsValid() bool {
	return p.Line > 0 && p.Column > 0
}

// TextPosition values.
var (
	TextPositionZero = TextPosition{}
	// TextPositionStart is the start position of any text.
	TextPositionStart = TextPosition{1, 1}
)

// TextRange is the value type that represents the range of a text.
type TextRange struct {
	Start TextPosition
	End   TextPosition
}

func (r TextRange) String() string {
	return fmt.Sprintf("[%s,%s]", r.Start, r.End)
}

//---

// Transactional has methods to make transaction.
type Transactional interface {
	Begin()
	Commit() error
	Rollback() error
}

//---

// Transaction is a helper for a Transactional type.
type Transaction struct {
	t    Transactional
	done bool
}

// NewTransaction creates a Transaction instance.
func NewTransaction(t Transactional) *Transaction {
	txn := &Transaction{t, false}
	txn.t.Begin()
	return txn
}

// Guard is assumed to be used with deffer.
func (t *Transaction) Guard() {
	if !t.done {
		_ = t.t.Rollback()
	}
}

// Commit calls the Commit method of the transactional instance.
func (t *Transaction) Commit() error {
	t.done = true
	return t.t.Commit()
}

// Rollback calls the Rollback method of the transactional instance.
// If using Guard, Rollback doesn't need to be called explicitly.
func (t *Transaction) Rollback() error {
	t.done = true
	return t.t.Rollback()
}

//---

// Reader is the interface for parser.
type Reader interface {
	io.RuneReader
	Transactional
	// CurrentPosition returns the current position of the cursor.
	CurrentPosition() TextPosition
}

//---

// LineBreak is returned ...
const LineBreak rune = '\n'

// InMemoryReader is an implementation of Reader.
type InMemoryReader struct {
	lines    [][]rune
	pos      TextPosition
	posStack *list.List
}

// NewInMemoryReader creates an InMemoryReader instance as Reader.
func NewInMemoryReader(buf []byte) Reader {
	numLines := bytes.Count(buf, []byte{'\n'}) // for less reallocation
	r := &InMemoryReader{
		lines:    make([][]rune, 0, numLines),
		pos:      TextPositionStart,
		posStack: list.New(),
	}
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		r.lines = append(r.lines, []rune(scanner.Text()))
	}
	return r
}

// ReadRune is a method of Reader interface.
func (r *InMemoryReader) ReadRune() (rune, int, error) {
	if !r.pos.IsValid() || len(r.lines) <= r.pos.Line-1 {
		return 0, 0, io.EOF
	}
	line := r.lines[r.pos.Line-1]
	var rune rune
	if len(line) > r.pos.Column-1 {
		rune = line[r.pos.Column-1]
		r.pos.Column++
	} else {
		rune = LineBreak
		r.pos.Line++
		r.pos.Column = 1
	}
	return rune, utf8.RuneLen(rune), nil
}

// CurrentPosition is a method of Reader interface.
func (r *InMemoryReader) CurrentPosition() TextPosition {
	return r.pos
}

// Begin is a method of Reader interface.
func (r *InMemoryReader) Begin() {
	r.posStack.PushBack(r.pos)
}

// Commit is a method of Reader interface.
func (r *InMemoryReader) Commit() error {
	back := r.posStack.Back()
	if back == nil {
		return errors.New("no active transaction")
	}
	r.posStack.Remove(back)
	return nil
}

// Rollback is a method of Reader interface.
func (r *InMemoryReader) Rollback() error {
	back := r.posStack.Back()
	if back == nil {
		return errors.New("no active transaction")
	}
	pos := r.posStack.Remove(back).(TextPosition)
	r.pos = pos
	return nil
}

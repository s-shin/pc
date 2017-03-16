package pc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextPositionString(t *testing.T) {
	assert.Equal(t, "1:2", TextPosition{1, 2}.String())
}

func TestTextRangeString(t *testing.T) {
	assert.Equal(t, "[1:2,3:4]", TextRange{TextPosition{1, 2}, TextPosition{3, 4}}.String())
}

func TestInMemoryReader(t *testing.T) {
	r := NewInMemoryReader([]byte("ab\nc\n\nd\n"))

	//    1  2  3
	// 1: a  b  \n
	// 2: c  \n
	// 3: \n
	// 4: d  \n

	type Expected struct {
		BeforePos   TextPosition
		Rune        rune
		IsReadError bool
	}

	for _, exp := range []Expected{
		{TextPosition{1, 1}, 'a', false},
		{TextPosition{1, 2}, 'b', false},
		{TextPosition{1, 3}, '\n', false},
		{TextPosition{2, 1}, 'c', false},
		{TextPosition{2, 2}, '\n', false},
		{TextPosition{3, 1}, '\n', false},
		{TextPosition{4, 1}, 'd', false},
		{TextPosition{4, 2}, '\n', false},
		{TextPosition{5, 1}, 0, true},
	} {
		assert.Equal(t, exp.BeforePos, r.CurrentPosition())
		rune, _, err := r.ReadRune()
		assert.Equal(t, string(exp.Rune), string(rune))
		assert.Equal(t, exp.IsReadError, err != nil)
	}
}

func TestInMemoryReaderTransaction(t *testing.T) {
	r := NewInMemoryReader([]byte("123456789"))

	r.ReadRune()
	assert.Equal(t, TextPosition{1, 2}, r.CurrentPosition())

	r.Begin()
	r.ReadRune()
	r.ReadRune()
	r.Rollback()
	assert.Equal(t, TextPosition{1, 2}, r.CurrentPosition())

	r.Begin()
	r.ReadRune()
	r.ReadRune()
	r.Commit()
	assert.Equal(t, TextPosition{1, 4}, r.CurrentPosition())

	func() {
		txn := NewTransaction(r)
		defer txn.Guard()
		r.ReadRune()
		r.ReadRune()
	}()
	assert.Equal(t, TextPosition{1, 4}, r.CurrentPosition())

	func() {
		txn := NewTransaction(r)
		defer txn.Guard()
		r.ReadRune()
		r.ReadRune()
		txn.Rollback()
	}()
	assert.Equal(t, TextPosition{1, 4}, r.CurrentPosition())

	func() {
		txn := NewTransaction(r)
		defer txn.Guard()
		r.ReadRune()
		r.ReadRune()
		txn.Commit()
	}()
	assert.Equal(t, TextPosition{1, 6}, r.CurrentPosition())

	func() {
		txn := NewTransaction(r)
		defer txn.Guard()
		r.ReadRune()
		assert.Equal(t, TextPosition{1, 7}, r.CurrentPosition())
		r.Begin()
		r.ReadRune()
		assert.Equal(t, TextPosition{1, 8}, r.CurrentPosition())
		r.Rollback()
		assert.Equal(t, TextPosition{1, 7}, r.CurrentPosition())
	}()
	assert.Equal(t, TextPosition{1, 6}, r.CurrentPosition())

	assert.Error(t, r.Commit())
	assert.Error(t, r.Rollback())
}

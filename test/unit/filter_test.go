package unit_test

import (
	"testing"

	"kickbase/internal/handler"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter_Operators(t *testing.T) {
	t.Run("EQ Operator", func(t *testing.T) {
		f := handler.ParseFilter("EQ:Jakarta", handler.OpCT)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpEQ, f.Op)
		assert.Equal(t, "Jakarta", f.Value)
	})

	t.Run("CT Operator (Contains)", func(t *testing.T) {
		f := handler.ParseFilter("CT:Pers", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpCT, f.Op)
		assert.Equal(t, "Pers", f.Value)
	})

	t.Run("IN Operator", func(t *testing.T) {
		f := handler.ParseFilter("IN:CF,SS,AMF", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpIN, f.Op)
		assert.Equal(t, []string{"CF", "SS", "AMF"}, f.Values)
	})

	t.Run("NI Operator (Not In)", func(t *testing.T) {
		f := handler.ParseFilter("NI:GK,CB", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpNI, f.Op)
		assert.Equal(t, []string{"GK", "CB"}, f.Values)
	})

	t.Run("GT Operator", func(t *testing.T) {
		f := handler.ParseFilter("GT:1930", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpGT, f.Op)
		assert.Equal(t, "1930", f.Value)
	})

	t.Run("GTE Operator", func(t *testing.T) {
		f := handler.ParseFilter("GTE:180", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpGTE, f.Op)
		assert.Equal(t, "180", f.Value)
	})

	t.Run("LT Operator", func(t *testing.T) {
		f := handler.ParseFilter("LT:70", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpLT, f.Op)
		assert.Equal(t, "70", f.Value)
	})

	t.Run("LTE Operator", func(t *testing.T) {
		f := handler.ParseFilter("LTE:75", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpLTE, f.Op)
		assert.Equal(t, "75", f.Value)
	})

	t.Run("BT Operator (Between)", func(t *testing.T) {
		f := handler.ParseFilter("BT:1920,1940", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpBT, f.Op)
		assert.Equal(t, []string{"1920", "1940"}, f.Values)
	})

	t.Run("Fallback Default Operator", func(t *testing.T) {
		f := handler.ParseFilter("Jakarta", handler.OpEQ)
		assert.NotNil(t, f)
		assert.Equal(t, handler.OpEQ, f.Op)
		assert.Equal(t, "Jakarta", f.Value)
	})
}

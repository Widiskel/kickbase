package handler

import (
	"strings"

	"gorm.io/gorm"
)

// FilterOp defines the operation type
type FilterOp string

const (
	OpEQ  FilterOp = "EQ"  // Equals
	OpCT  FilterOp = "CT"  // Contains (case-insensitive substring)
	OpIN  FilterOp = "IN"  // In list (comma separated)
	OpNI  FilterOp = "NI"  // Not in list (comma separated)
	OpGT  FilterOp = "GT"  // Greater than
	OpGTE FilterOp = "GTE" // Greater than or equal
	OpLT  FilterOp = "LT"  // Less than
	OpLTE FilterOp = "LTE" // Less than or equal
	OpBT  FilterOp = "BT"  // Between (comma separated: val1,val2)
)

// FilterCriterion holds parsed filter details
type FilterCriterion struct {
	Op    FilterOp
	Value string
	Values []string
}

// ParseFilter parses an input string into a FilterCriterion.
// Examples:
// - "CT:Pers" -> Op: CT, Value: "Pers"
// - "IN:CF,SS,AMF" -> Op: IN, Values: ["CF", "SS", "AMF"]
// - "BT:170,185" -> Op: BT, Values: ["170", "185"]
// - "Jakarta" -> Op: EQ, Value: "Jakarta"
func ParseFilter(input string, defaultOp FilterOp) *FilterCriterion {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil
	}

	parts := strings.SplitN(trimmed, ":", 2)
	if len(parts) == 2 {
		opCandidate := FilterOp(strings.ToUpper(strings.TrimSpace(parts[0])))
		switch opCandidate {
		case OpEQ, OpCT, OpGT, OpGTE, OpLT, OpLTE:
			return &FilterCriterion{
				Op:    opCandidate,
				Value: strings.TrimSpace(parts[1]),
			}
		case OpIN, OpNI:
			rawVals := strings.Split(parts[1], ",")
			var vals []string
			for _, v := range rawVals {
				vTrim := strings.TrimSpace(v)
				if vTrim != "" {
					vals = append(vals, vTrim)
				}
			}
			return &FilterCriterion{
				Op:     opCandidate,
				Values: vals,
			}
		case OpBT:
			rawVals := strings.Split(parts[1], ",")
			var vals []string
			for _, v := range rawVals {
				vTrim := strings.TrimSpace(v)
				if vTrim != "" {
					vals = append(vals, vTrim)
				}
			}
			if len(vals) >= 2 {
				return &FilterCriterion{
					Op:     OpBT,
					Values: []string{vals[0], vals[1]},
				}
			}
		}
	}

	// Fallback to default operator
	if defaultOp == OpCT {
		return &FilterCriterion{
			Op:    OpCT,
			Value: trimmed,
		}
	}

	return &FilterCriterion{
		Op:    OpEQ,
		Value: trimmed,
	}
}

// ApplyFilterToQuery applies the FilterCriterion to a GORM query securely
func ApplyFilterToQuery(db *gorm.DB, column string, fc *FilterCriterion) *gorm.DB {
	if fc == nil {
		return db
	}

	switch fc.Op {
	case OpEQ:
		return db.Where(column+" = ?", fc.Value)
	case OpCT:
		return db.Where("LOWER("+column+") LIKE ?", "%"+strings.ToLower(fc.Value)+"%")
	case OpIN:
		if len(fc.Values) > 0 {
			return db.Where(column+" IN ?", fc.Values)
		}
	case OpNI:
		if len(fc.Values) > 0 {
			return db.Where(column+" NOT IN ?", fc.Values)
		}
	case OpGT:
		return db.Where(column+" > ?", fc.Value)
	case OpGTE:
		return db.Where(column+" >= ?", fc.Value)
	case OpLT:
		return db.Where(column+" < ?", fc.Value)
	case OpLTE:
		return db.Where(column+" <= ?", fc.Value)
	case OpBT:
		if len(fc.Values) >= 2 {
			return db.Where(column+" BETWEEN ? AND ?", fc.Values[0], fc.Values[1])
		}
	}

	return db
}

// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
)

func TestPlus(t *testing.T) {
	var testCases = []struct {
		name        string
		left, right float64
		expected    string
	}{
		{"1 + 1", 1, 1, "2"},
		{"-1 + 1", -1, 1, "0"},
		{"0 + 0", 0, 0, "0"},
		{"0.14159 + 3.0", 0.14159, 3.0, "3.14159"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			result, err := NewPlus(
				NewLiteral(tt.left, sql.Float64),
				NewLiteral(tt.right, sql.Float64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(err)
			r, ok := result.(decimal.Decimal)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
		})
	}

	require := require.New(t)
	result, err := NewPlus(NewLiteral("2", sql.LongText), NewLiteral(3, sql.Float64)).
		Eval(sql.NewEmptyContext(), sql.NewRow())
	require.NoError(err)
	r, ok := result.(decimal.Decimal)
	require.True(ok)
	require.Equal("5", r.StringFixed(r.Exponent()*-1))
}

func TestPlusInterval(t *testing.T) {
	require := require.New(t)

	expected := time.Date(2018, time.May, 2, 0, 0, 0, 0, time.UTC)
	op := NewPlus(
		NewLiteral("2018-05-01", sql.LongText),
		NewInterval(NewLiteral(int64(1), sql.Int64), "DAY"),
	)

	result, err := op.Eval(sql.NewEmptyContext(), nil)
	require.NoError(err)
	require.Equal(expected, result)

	op = NewPlus(
		NewInterval(NewLiteral(int64(1), sql.Int64), "DAY"),
		NewLiteral("2018-05-01", sql.LongText),
	)

	result, err = op.Eval(sql.NewEmptyContext(), nil)
	require.NoError(err)
	require.Equal(expected, result)
}

func TestMinus(t *testing.T) {
	var testCases = []struct {
		name        string
		left, right float64
		expected    string
	}{
		{"1 - 1", 1, 1, "0"},
		{"1 - -1", 1, -1, "2"},
		{"0 - 0", 0, 0, "0"},
		{"3.14159 - 3.0", 3.14159, 3.0, "0.14159"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			result, err := NewMinus(
				NewLiteral(tt.left, sql.Float64),
				NewLiteral(tt.right, sql.Float64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(err)
			r, ok := result.(decimal.Decimal)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
		})
	}

	require := require.New(t)
	result, err := NewMinus(NewLiteral("10", sql.LongText), NewLiteral(10, sql.Int64)).
		Eval(sql.NewEmptyContext(), sql.NewRow())
	require.NoError(err)
	r, ok := result.(decimal.Decimal)
	require.True(ok)
	require.Equal("0", r.StringFixed(r.Exponent()*-1))
}

func TestMinusInterval(t *testing.T) {
	require := require.New(t)

	expected := time.Date(2018, time.May, 1, 0, 0, 0, 0, time.UTC)
	op := NewMinus(
		NewLiteral("2018-05-02", sql.LongText),
		NewInterval(NewLiteral(int64(1), sql.Int64), "DAY"),
	)

	result, err := op.Eval(sql.NewEmptyContext(), nil)
	require.NoError(err)
	require.Equal(expected, result)
}

func TestMult(t *testing.T) {
	var testCases = []struct {
		name        string
		left, right float64
		expected    string
	}{
		{"1 * 1", 1, 1, "1"},
		{"-1 * 1", -1, 1, "-1"},
		{"0 * 0", 0, 0, "0"},
		{"3.14159 * 3.0", 3.14159, 3.0, "9.42477"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			result, err := NewMult(
				NewLiteral(tt.left, sql.Float64),
				NewLiteral(tt.right, sql.Float64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(err)
			r, ok := result.(decimal.Decimal)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
		})
	}

	require := require.New(t)
	result, err := NewMult(NewLiteral("10", sql.LongText), NewLiteral("10", sql.LongText)).
		Eval(sql.NewEmptyContext(), sql.NewRow())
	require.NoError(err)
	r, ok := result.(decimal.Decimal)
	require.True(ok)
	require.Equal("100", r.StringFixed(r.Exponent()*-1))
}

func TestIntDiv(t *testing.T) {
	var testCases = []struct {
		name        string
		left, right int64
		expected    int64
		null        bool
	}{
		{"1 div 1", 1, 1, 1, false},
		{"8 div 3", 8, 3, 2, false},
		{"1 div 3", 1, 3, 0, false},
		{"0 div -1024", 0, -1024, 0, false},
		{"1 div 0", 1, 0, 0, true},
		{"0 div 0", 1, 0, 0, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			result, err := NewIntDiv(
				NewLiteral(tt.left, sql.Int64),
				NewLiteral(tt.right, sql.Int64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(err)
			if tt.null {
				assert.Equal(t, nil, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMod(t *testing.T) {
	var testCases = []struct {
		name        string
		left, right int64
		expected    string
		null        bool
	}{
		{"1 % 1", 1, 1, "0", false},
		{"8 % 3", 8, 3, "2", false},
		{"1 % 3", 1, 3, "1", false},
		{"0 % -1024", 0, -1024, "0", false},
		{"-1 % 2", -1, 2, "-1", false},
		{"1 % -2", 1, -2, "1", false},
		{"-1 % -2", -1, -2, "-1", false},
		{"1 % 0", 1, 0, "0", true},
		{"0 % 0", 0, 0, "0", true},
		{"0.5 % 0.24", 0, 0, "0.02", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			result, err := NewMod(
				NewLiteral(tt.left, sql.Int64),
				NewLiteral(tt.right, sql.Int64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(err)
			if tt.null {
				require.Nil(result)
			} else {
				r, ok := result.(decimal.Decimal)
				require.True(ok)
				require.Equal(tt.expected, r.StringFixed(r.Exponent()*-1))
			}
		})
	}
}

func TestAllFloat64(t *testing.T) {
	var testCases = []struct {
		op       string
		value    float64
		expected string
	}{
		// The value here are given with decimal place to force the value type to float, but the interpreted values
		// will not have 0 scale, so the mult is 3.0000 * 0 = 0.0000 instead of 3.0000 * 0.0 = 0.00000
		{"+", 1.0, "1"},
		{"-", -8.0, "9"},
		{"/", 3.0, "3.0000"},
		{"*", 4.0, "12.0000"},
		{"%", 11, "1.0000"},
	}

	// ((((0 + 1) - (-8)) / 3) * 4) % 11 == 1
	lval := NewLiteral(float64(0.0), sql.Float64)
	for _, tt := range testCases {
		t.Run(tt.op, func(t *testing.T) {
			require := require.New(t)
			var result interface{}
			var err error
			if tt.op == "/" {
				result, err = NewDiv(lval,
					NewLiteral(tt.value, sql.Float64),
				).Eval(sql.NewEmptyContext(), sql.NewRow())
			} else if tt.op == "%" {
				result, err = NewMod(lval,
					NewLiteral(tt.value, sql.Float64),
				).Eval(sql.NewEmptyContext(), sql.NewRow())
			} else {
				result, err = NewArithmetic(lval,
					NewLiteral(tt.value, sql.Float64), tt.op,
				).Eval(sql.NewEmptyContext(), sql.NewRow())
			}
			require.NoError(err)
			if r, ok := result.(decimal.Decimal); ok {
				assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
			} else {
				assert.Equal(t, tt.expected, result)
			}

			lval = NewLiteral(result, sql.Float64)
		})
	}
}

func TestUnaryMinus(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		typ      sql.Type
		expected interface{}
	}{
		{"int32", int32(1), sql.Int32, int32(-1)},
		{"uint32", uint32(1), sql.Uint32, int32(-1)},
		{"int64", int64(1), sql.Int64, int64(-1)},
		{"uint64", uint64(1), sql.Uint64, int64(-1)},
		{"float32", float32(1), sql.Float32, float32(-1)},
		{"float64", float64(1), sql.Float64, float64(-1)},
		{"int text", "1", sql.LongText, "-1"},
		{"float text", "1.2", sql.LongText, "-1.2"},
		{"nil", nil, sql.LongText, nil},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := NewUnaryMinus(NewLiteral(tt.input, tt.typ))
			result, err := f.Eval(sql.NewEmptyContext(), nil)
			require.NoError(t, err)
			if dt, ok := result.(decimal.Decimal); ok {
				require.Equal(t, tt.expected, dt.StringFixed(dt.Exponent()*-1))
			} else {
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

package envconf

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoadInt(t *testing.T) {
	var dest struct {
		Int   int
		Int8  int8
		Int16 int16
		Int32 int32
		Int64 int64
	}
	settings := map[string]string{
		"INT":   "99",
		"INT8":  "8",
		"INT16": "16",
		"INT32": "32",
		"INT64": "64",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.Int != 99 || dest.Int8 != 8 || dest.Int16 != 16 || dest.Int32 != 32 || dest.Int64 != 64 {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestLoadFloat(t *testing.T) {
	var dest struct {
		Fl32 float64
		Fl64 float64
	}
	settings := map[string]string{
		"FL32": "32.32",
		"FL64": "64.64",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.Fl32 != 32.32 || dest.Fl64 != 64.64 {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestLoadBool(t *testing.T) {
	type conf struct {
		Val bool
	}

	cases := map[string]bool{
		"1":    true,
		"t":    true,
		"T":    true,
		"true": true,
		"True": true,
		"TRUE": true,

		"0":     false,
		"f":     false,
		"F":     false,
		"false": false,
		"False": false,
		"FALSE": false,
	}

	for raw, exp := range cases {
		var dest conf
		if err := Load(&dest, map[string]string{"VAL": raw}); err != nil {
			t.Fatalf("cannot load: %s", err)
		}
		if dest.Val != exp {
			t.Errorf("%s -> %v: invalid result: %#v", raw, exp, dest)
		}
	}
}

func TestLoadString(t *testing.T) {
	var dest struct {
		Val string
		N   string
	}
	settings := map[string]string{
		"VAL": "foo bar",
		"N":   "321",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.Val != "foo bar" || dest.N != "321" {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestLoadBytes(t *testing.T) {
	var dest struct {
		Val []byte
	}
	settings := map[string]string{
		"VAL": "foo bar",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if string(dest.Val) != "foo bar" {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestLoadSliceOfString(t *testing.T) {
	var dest struct {
		S1 []string
		S2 []string
	}
	settings := map[string]string{
		"S1": "foo,bar",
		"S2": "xx yy",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if want, got := []string{"foo", "bar"}, dest.S1; !reflect.DeepEqual(want, got) {
		t.Errorf("S1: want %v, got %v", want, got)
	}
	if want, got := []string{"xx yy"}, dest.S2; !reflect.DeepEqual(want, got) {
		t.Errorf("S2: want %v, got %v", want, got)
	}
}

func TestLoadSliceOfInt(t *testing.T) {
	var dest struct {
		I1 []int
		I2 []int
	}
	settings := map[string]string{
		"I1": "11,22",
		"I2": "104",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if want, got := []int{11, 22}, dest.I1; !reflect.DeepEqual(want, got) {
		t.Errorf("I1: want %v, got %v", want, got)
	}
	if want, got := []int{104}, dest.I2; !reflect.DeepEqual(want, got) {
		t.Errorf("I2: want %v, got %v", want, got)
	}
}

func TestLoadSliceOfBool(t *testing.T) {
	var dest struct {
		B1 []bool
		B2 []bool
	}
	settings := map[string]string{
		"B1": "0,t,true,false",
		"B2": "f",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if want, got := []bool{false, true, true, false}, dest.B1; !reflect.DeepEqual(want, got) {
		t.Errorf("B1: want %v, got %v", want, got)
	}
	if want, got := []bool{false}, dest.B2; !reflect.DeepEqual(want, got) {
		t.Errorf("B2: want %v, got %v", want, got)
	}
}

func TestLoadSliceOfFloat(t *testing.T) {
	var dest struct {
		F1 []float32
		F2 []float32
	}
	settings := map[string]string{
		"F1": "11.32, 44, 12",
		"F2": "22.321",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if want, got := []float32{11.32, 44, 12}, dest.F1; !reflect.DeepEqual(want, got) {
		t.Errorf("F1: want %v, got %v", want, got)
	}
	if want, got := []float32{22.321}, dest.F2; !reflect.DeepEqual(want, got) {
		t.Errorf("F2: want %v, got %v", want, got)
	}
}

func TestLoadWithTagName(t *testing.T) {
	var dest struct {
		X string `envconf:"value"`
	}
	settings := map[string]string{
		"x":     "x",
		"X":     "X",
		"VALUE": "VALUE",
		"value": "value", // as defined by tag
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.X != "value" {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestLoadWithMissingRequired(t *testing.T) {
	var dest struct {
		Val string `envconf:",required"`
	}

	settings := map[string]string{
		"Val": "Val",
		"val": "val",
		// "VAL": "",  missing
	}
	err := Load(&dest, settings)
	if err == nil {
		t.Fatal("want error")
	}
	errs := err.(ParseErrors)
	if len(errs) != 1 {
		t.Errorf("want 1 error, got %d", len(errs))
	}
	e := errs[0]
	if e.Field != "Val" || e.Err != errRequired {
		t.Errorf("invalid error: %#v", e)
	}
}

func TestLoadWithEmptyRequired(t *testing.T) {
	var dest struct {
		Val string `envconf:",required"`
	}

	settings := map[string]string{
		"Val": "Val",
		"val": "val",
		"VAL": "", // empty value != missing
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.Val != "" {
		t.Errorf("invalid result: %#v", dest)
	}
}

func TestCustomFieldType(t *testing.T) {
	var dest struct {
		Val CustomField
	}

	settings := map[string]string{
		"VAL": "x",
	}
	if err := Load(&dest, settings); err != nil {
		t.Fatalf("cannot load: %s", err)
	}
	if dest.Val.val != "custom: x" {
		t.Errorf("invalid result: %#v", dest)
	}
}

type CustomField struct {
	val string
}

func (cf *CustomField) UnmarshalText(b []byte) error {
	cf.val = "custom: " + string(b)
	return nil
}

func (cf *CustomField) MarshalText() ([]byte, error) {
	return []byte(cf.val), nil
}

func TestDescribe(t *testing.T) {
	cases := map[string]struct {
		dest    interface{}
		want    string
		wantErr error
	}{
		"empty": {
			dest: &struct{}{},
			want: "",
		},
		"with_defaults": {
			dest: &struct {
				Str  string
				Int  int32   `envconf:"integer"`
				Bool bool    `envconf:",required"`
				Fl   float64 `envconf:"age"`
			}{
				Str:  "bob",
				Int:  32,
				Bool: true,
				Fl:   0.2,
			},
			want: `
STR      string   "bob"
integer  int32    "32"
BOOL     bool     "true"
age      float64  "0.2"
                        `,
		},
		"no_default": {
			dest: &struct {
				Str    string
				Int    int32       `envconf:"integer"`
				Bool   bool        `envconf:",required"`
				Custom CustomField `envconf:"custom,required"`
				Fl     float64     `envconf:"age"`
				StrArr []string
				IntArr []int16
			}{},
			want: `
STR      string
integer  int32
BOOL     bool         (required)
custom   CustomField  (required)
age      float64
STR_ARR  string list
INT_ARR  int16 list
                        `,
		},
		"bytes": {
			dest: &struct {
				Raw   []byte
				Bytes []byte `envconf:",required"`
			}{},
			want: `
RAW    bytes
BYTES  bytes  (required)
                        `,
		},
	}

	for tname, tc := range cases {
		desc, err := Describe(tc.dest)
		if err != nil {
			if tc.wantErr != err {
				t.Errorf("%s: want %q error, got %q", tname, tc.wantErr, err)
			}
			continue
		}

		if want, got := strings.TrimSpace(tc.want), strings.TrimSpace(desc); want != got {
			t.Errorf(`%s:
-- want --
%s
-- got --
%s

`, tname, want, got)
		}
	}
}

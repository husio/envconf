package envconf

import (
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

func TestLoadSliceOfString(t *testing.T) {
}

func TestLoadSliceOfInt(t *testing.T) {
}

func TestLoadSliceOfBool(t *testing.T) {
}

func TestLoadSliceOfFloat(t *testing.T) {
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

func (cf *CustomField) ParseConf(s string) error {
	cf.val = "custom: " + s
	return nil
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

package log

import "testing"

func TestSymbols_String(t *testing.T) {

	tests := []struct {
		name    string
		symbols Symbols
		want    string
	}{
		// TODO: Add test cases.
		{"base", Symbols{}, "Debug:  Info:  Success:  Warning:  Error:  Fatal: "},
		{"base", Symbols{Debug: "Debug"}, "Debug: Debug Info:  Success:  Warning:  Error:  Fatal: "},
		{"base", Symbols{
			Debug:   Symbol("λ"),
			Info:    Symbol("ℹ"),
			Success: Symbol("✔"),
			Warning: Symbol("⚠"),
			Error:   Symbol("!!"),
			Fatal:   Symbol("✖"),
		}, "Debug: λ Info: ℹ Success: ✔ Warning: ⚠ Error: !! Fatal: ✖"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Symbols{
				Debug:   tt.symbols.Debug,
				Info:    tt.symbols.Info,
				Warning: tt.symbols.Warning,
				Warn:    tt.symbols.Warn,
				Error:   tt.symbols.Error,
				Fatal:   tt.symbols.Fatal,
				Success: tt.symbols.Success,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("\nSymbols.String() = %v, \nwant = %v", got, tt.want)
			}
		})
	}
}

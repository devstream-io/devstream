package log

import (
	"fmt"
)

var colorOn bool

type Symbol string

// Symbols struct contains all symbols
type Symbols struct {
	Debug   Symbol
	Info    Symbol
	Warning Symbol
	Warn    Symbol
	Error   Symbol
	Fatal   Symbol
	Success Symbol
}

var osBaseSymbols Symbols

var normal = Symbols{
	Debug:   Symbol("ℹ"),
	Info:    Symbol("ℹ"),
	Success: Symbol("✔"),
	Warning: Symbol("⚠"),
	Warn:    Symbol("⚠"),
	Error:   Symbol("✖"),
	Fatal:   Symbol("✖"),
}

var fallback = Symbols{
	Debug:   Symbol("i"),
	Info:    Symbol("i"),
	Success: Symbol("√"),
	Warning: Symbol("‼"),
	Warn:    Symbol("‼"),
	Error:   Symbol("×"),
}

// String returns a printable representation of Symbols struct
func (s Symbols) String() string {
	return fmt.Sprintf("Info: %s Success: %s Warning: %s Error: %s", s.Info, s.Success, s.Warning, s.Error)
}

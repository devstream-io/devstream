package log

import (
	"fmt"
)

var colorOn bool

type Symbol string

//var (
//	// Debug represents the debug symbol
//	debug Symbol
//	// Info represents the information symbol
//	info Symbol
//	// Warning represents the warning symbol
//	warning Symbol
//	// Warn alias of Warning
//	warn Symbol
//	// Error represents the error symbol
//	error Symbol
//	// Fatal represents the fatal symbol
//	fatal Symbol
//	// Success represents the success symbol
//	success Symbol
//)

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

package main

type Plugin interface {
	Health() (bool, error)
	Name() string
}

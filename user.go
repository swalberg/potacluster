package main

const (
	NotLoggedIn = iota
	LoggedIn
)

type user struct {
	CallSign string
	State    int
}

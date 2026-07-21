package main

type ResErr struct {
	Error string `json:"error"`
}

type ResVersion struct {
	Version []int `json:"version"`
}

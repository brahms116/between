package main

type lspState struct {
	documents map[string]string
}

func newLspState() *lspState {
	return &lspState{
		documents: make(map[string]string),
	}
}

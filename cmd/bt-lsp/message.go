package main

type RpcRequest[T any] struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params T      `json:"params"`
}

type RpcResponse[T any] struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Result *T     `json:"result,omitempty"`
}

type RpcNotification[T any] struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
	Params T      `json:"params"`
}

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}

package model

type Response struct {
	Result Result
	Method string
}

type Result struct {
	Orders []Order
}

type Order struct {
	Type int
}

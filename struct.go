package main

const (
	TypeString = iota
	TypeInt
	TypeInt64
	TypeBool
	UnsupportedType
)

type dataNode struct {
	Type    int
	Data interface{}
}


package main

import "unsafe"

const (
	TypeString = iota
	TypeInt
	TypeInt64
	TypeBool
	UnsupportedType
)

type dataNode struct {
	Type    int
	Pointer unsafe.Pointer
}

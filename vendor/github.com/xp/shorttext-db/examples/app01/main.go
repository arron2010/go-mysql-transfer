package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

//实验字节数组与整型的转换
//实验GO指针的用法

func testPointer() {
	var data uint32
	data = 0x12345678
	var buf []byte
	buf = make([]byte, 4, 4)
	a := unsafe.Pointer(&buf[0])
	b := (*uint32)(a)
	*b = data
	fmt.Println("指针方式转换：", getHexForBytes(buf))
	buf2 := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buf2, data)
	fmt.Println("小端方式：", getHexForBytes(buf2))

	buf3 := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(buf3, data)
	fmt.Println("大端方式：", getHexForBytes(buf3))
}
func getHexForBytes(buf []byte) string {
	var str string
	for i := 0; i < len(buf); i++ {
		str = str + fmt.Sprintf("%X ", buf[i])
	}
	return str
}

type IDemo1 interface {
	Print1()
}

type IDemo2 interface {
	Print2()
}

type Demo struct {
	data int
}

func (d *Demo) Print1() {
	fmt.Println("IDemo1---->")
}

func (d1 Demo) Print2() {
	d1.data = 100
	fmt.Println("IDemo2---->")
}

func testInterface() {
	var d Demo
	//var dd IDemo1
	//dd =d
	d.Print2()
	fmt.Println(d.data)
}

type options struct {
	width int
	size  int
	name  string
}

func WithWidth(width int) Option {
	return func(o *options) {
		o.width = width
	}
}

type Option func(*options)

type Shape struct {
	opts options
}

func WithSize(size int) Option {
	return func(o *options) {
		fmt.Println("WithSize-->", size)
		o.size = size
	}
}

func testOptions() {
	size := 100
	opts := &options{}
	opt := WithSize(size)
	opt(opts)
	fmt.Println("Shape-->", opts.size)
}

func NewShape(opts ...Option) *Shape {
	instance := &Shape{}
	for _, optHandle := range opts {
		optHandle(&instance.opts)
	}
	return instance
}
func main() {
	testOptions()
	//fmt.Printf("%X\n",15)
	//testPointer()
	//testInterface()
}

//type Demo1 struct {
//	Demo
//	data int
//}

//func  (d *Demo1)Print3(){
//	d.Print1()
//	fmt.Println("IDemo1---->Demo1")
//}
//
//func (d *Demo1)Print1(){
//	//d.Demo.Print1()
//	fmt.Println("IDemo2---->")
//}

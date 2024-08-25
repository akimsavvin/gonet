// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package main

import (
	"fmt"
	"github.com/akimsavvin/gonet/di"
)

type MyInterface interface {
	Method1() string
}

type MyClass struct {
	Huy string
}

func NewMyClass() *MyClass {
	return &MyClass{
		Huy: "Huy",
	}
}

func (c *MyClass) Method1() string {
	return "Hello World"
}

type MyClass2 struct {
	Her string
}

func NewMyClass2() *MyClass2 {
	return &MyClass2{
		Her: "Her",
	}
}

func (c *MyClass2) Method1() string {
	return "Hello World 2"
}

type HerInterface interface {
	Huynya() string
}

type HerService struct {
	MyI  MyInterface
	MyI2 MyInterface
}

func NewHerService(my []MyInterface) *HerService {
	return &HerService{
		MyI:  my[0],
		MyI2: my[1],
	}
}

func (h *HerService) Huynya() string {
	return h.MyI.Method1() + "+" + h.MyI2.Method1()
}

type Huy struct {
	Hello string
}

func main() {
	di.AddSingleton[MyInterface](NewMyClass)
	di.AddSingleton[MyInterface](NewMyClass2)
	di.AddTransient[HerInterface](NewHerService)

	di.Build()

	her := di.GetRequiredService[HerInterface]()
	fmt.Println(her.Huynya())
}

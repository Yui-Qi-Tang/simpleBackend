package testpkg

import (
	"fmt"
	"reflect"
)

// IntAdd sum of x and y
func IntAdd(x int, y int) int {
	return x + y
}

func checkDataType(data interface{}) {
	fmt.Println(reflect.TypeOf(data))
}

type car1 struct {
	model string
	price int
	speed int //km/hr
}

func (c *car1) backward() {
	fmt.Println(c.model, " backward")
}

func (c *car1) forward() {
	fmt.Println(c.model, " forward")
}

func (c *car1) showInfo() {
	fmt.Println("Car Model: ", c.model)
	fmt.Println("Price: ", c.price, "$")
	fmt.Println("Speed: ", c.speed, "Km/hr")
}

type carTestPlan1 interface {
	carInfo
	backward()
}

type carTestPlan2 interface {
	carInfo
	forward()
}

type carInfo interface {
	showInfo()
}

func carTest(ct carTestPlan1) {
	ct.showInfo()
	ct.backward()
}

func carTest2(ct carTestPlan2) {
	ct.showInfo()
	ct.forward()
}

func carTestTemplate() {
	car1 := &car1{
		model: "car1",
		price: 100,
		speed: 100,
	}

	fmt.Println("First car testing...")
	carTest(car1)

	fmt.Println("Second car testing...")
	carTest2(car1)
}

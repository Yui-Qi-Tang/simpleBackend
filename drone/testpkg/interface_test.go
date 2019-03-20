package testpkg

import (
	"testing"
)

type person struct {
	name string
	ages int
}

func (p person) getNameWithCopy() string {
	return p.name
}

func (p *person) setName(newName string) {
	p.name = newName
}

func (p person) setName2(newName string) {
	p.name = newName
}

func (p *person) getName() string {
	return p.name
}

func (p *person) getPAdress(t *testing.T) {
	t.Logf("address of pointer receiver: %p", p)
}

func (p person) getOAdress(t *testing.T) {
	t.Logf("address of object receiver: %p", &p)
}

type personGetNameWithCopyTest interface {
	getNameWithCopy() string
	setName(name string)
	setName2(name string)
}

type getNameWithCopyTest interface {
	getNameWithCopy() string
	setName2(name string)
}

func setNewName(i personGetNameWithCopyTest, newName string, t *testing.T) {
	t.Log("before set new name:", i.getNameWithCopy())
	i.setName(newName)
	t.Log("after set new name:", i.getNameWithCopy())
}

func setNewName2(i getNameWithCopyTest, newName string, t *testing.T) {
	t.Log("before set new name:", i.getNameWithCopy())
	i.setName2(newName)
	t.Log("after set new name:", i.getNameWithCopy())
}

func setNewNameViaInstanceMethod(i getNameWithCopyTest, newName string, t *testing.T, p person) {
	t.Log("before set new name:", i.getNameWithCopy())
	p.setName(newName)
	t.Log("real p", p.getNameWithCopy())
	t.Log("after set new name:", i.getNameWithCopy())
}

func TestNormal(t *testing.T) {
	// carTestTemplate()
	/*
		var p1 person
		var p1i personGetNameWithCopyTest
		p1.ages = 12
		p1.name = "p1"
		p1i = p1
		t.Log("before set new name:", p1i.getNameWithCopy())
		p1.setName("new p1 name")
		t.Log("after set new name:", p1i.getNameWithCopy())
		t.Log("Check out real p1 name:", p1.getNameWithCopy())
	*/
	var p1 person
	p1.name = "p1"
	t.Log(p1.getName())
	p1.setName("New p1")
	t.Log(p1.getName())
	p1.getPAdress(t)
	p1.getOAdress(t)
	t.Logf("Real p1 address: %p", &p1)

	// setNewNameViaInstanceMethod(p1i, "p1 new name", t, p1)
	//p1.setName("latest p1  name")
	t.Log("Set p1 PASS", p1.getName())

	var p2 *person
	p2 = &person{
		name: "p2",
		ages: 19,
	}
	var p2i personGetNameWithCopyTest
	p2i = p2
	setNewName(p2i, "p2 new name", t)
	//t.Log("before set new name:", p2i.getNameWithCopy())

	t.Log("Set p2 PASS", p2.name)
}

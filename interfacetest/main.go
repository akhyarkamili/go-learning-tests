package main

import "fmt"

type creature interface {
	Breathe() error
}

type person struct {
}

func (person) Breathe() error {
	fmt.Println("Breathing")
	return nil
}

func resuscitate(creatures []creature) {
	for _, c := range creatures {
		_ = c.Breathe()
	}
}

func main() {
	people := []person{person{}, person{}}

	//resuscitate(people) // does NOT compile!

	var creatures []creature
	for _, p := range people {
		creatures = append(creatures, p)
	}
	resuscitate(creatures) // compiles
}

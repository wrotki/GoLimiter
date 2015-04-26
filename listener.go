package main

import (
	"container/list"
	"fmt"
	"time"
)

type Message interface{}

type Login struct {
	Id       string
	Password string
}

func (m Login) String() string {
	return fmt.Sprintf("Login Id: %v, Password: %v", m.Id, m.Password)
}

type Logout struct{}

func (m Logout) String() string {
	return "Logout"
}

type Done struct{}

func service(c chan Message){
	
}

func proxy(counter InFlightCounter, c chan Message) {
	for {
		var msg = <-c
		switch msg.(type) {
		case Login:
			fmt.Println(msg)
		case Logout:
			fmt.Println("I'm outta here\n")
		case Done:
			counter.Count--
		}
	}
}

const Max = 3

type InFlightCounter struct {
	Count int
	Queue *list.List
}

func main() {
	c := make(chan Message)
	go proxy(InFlightCounter{Queue: list.New()}, c)
	for i := 0; i < 5; i++ {
		c <- Login{Id: "Me", Password: "Foo"}
		c <- Done{}
		c <- nil
		time.Sleep(300 * time.Millisecond)
	}
}

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

func service(c,s chan Message) {
	for {
		var msg = <-s
		switch msg.(type) {
		case Login:
			fmt.Println(msg)
		case Logout:
			fmt.Println("I'm outta here\n")
		}
		c <- Done{}
	}
}

func proxy(counter InFlightCounter, c,s chan Message) {
	for {
		cq := counter.Queue
		cc := &counter.Count
		var msg = <-c
		if msg == nil {
			return
		}
		switch msg.(type) {
		case Done:
			if counter.Queue.Len() > 0 {
				e := cq.Back()
				cq.Remove(e)
				s <- e.Value
			} else {
				*cc--
			}
			fmt.Println("Done: count: ",*cc)
		default:
			if *cc > Max {
				cq.PushFront(msg)
			} else {
				s <- msg
				*cc++
			}
			fmt.Println("Message: count: ",*cc)
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
	s := make(chan Message)
	go service(c,s)
	go proxy(InFlightCounter{Queue: list.New()}, c, s)
	for i := 0; i < 5; i++ {
		c <- Login{Id: "Me", Password: "Foo"}
		c <- Logout{}
		c <- nil
		time.Sleep(300 * time.Millisecond)
	}
}

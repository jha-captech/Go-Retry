package main

import (
	"errors"
	"log"
	"time"
)

func main() {
	words := "these are words"

	stuff, err := retryWithReturn(5, 5*time.Second, func() (Stuff, error) {
		return newStuff(words)
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = retry(5, 5*time.Second, func() error {
		return stuff.doThingsWithStuff()
	})
	if err != nil {
		log.Fatalln(err)
	}
}

// test obj and funcs

type Stuff struct {
	words string
}

func newStuff(words string) (Stuff, error) {
	stuff := Stuff{
		words: words,
	}
	return stuff, errors.New("test error")
}

func (s Stuff) doThingsWithStuff() error {
	return errors.New("error in doThingsWithStuff")
}

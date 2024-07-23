package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	words := "these are words"

	stuff, err := RetryResult(ctx, 5*time.Second, func() (Stuff, error) {
		return newStuff(words)
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = Retry(ctx, 5*time.Second, func() error {
		return stuff.doThingsWithStuff()
	})
	if err != nil {
		log.Fatalln(err)
	}
}

// demo obj and funcs

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

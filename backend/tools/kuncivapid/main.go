package main

import (
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	privat, publik, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("VAPID_KUNCI_PUBLIK=%s\nVAPID_KUNCI_PRIVAT=%s\n", publik, privat)
}

package main

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

type Deck struct {
	Cards []*Card
}

var cardImages = []string{
	"2C.png", "2D.png", "2H.png", "2S.png",
	"3C.png", "3D.png", "3H.png", "3S.png",
	"4C.png", "4D.png", "4H.png", "4S.png",
	"5C.png", "5D.png", "5H.png", "5S.png",
	"6C.png", "6D.png", "6H.png", "6S.png",
	"7C.png", "7D.png", "7H.png", "7S.png",
	"8C.png", "8D.png", "8H.png", "8S.png",
	"9C.png", "9D.png", "9H.png", "9S.png",
	"10C.png", "10D.png", "10H.png", "10S.png",
	"AC.png", "AD.png", "AH.png", "AS.png",
	"JC.png", "JD.png", "JH.png", "JS.png",
	"KC.png", "KD.png", "KH.png", "KS.png",
	"QC.png", "QD.png", "QH.png", "QS.png",
}

func (d *Deck) init() {
	for _, image := range cardImages {
		card := &Card{
			Image: image,
		}
		d.Cards = append(d.Cards, card)
	}
}

func (d *Deck) shuffle() {
	var cards []*Card
	rand.Seed(time.Now().UnixNano())
	total := len(cardImages)
	for {
		index := rand.Intn(total)
		if d.Cards[index] != nil {
			cards = append(cards, d.Cards[index])
			d.Cards[index] = nil
			if len(cards) == total {
				break
			}
		}
	}
	d.Cards = cards
}

func (d *Deck) String() string {
	s := ""
	for i, card := range d.Cards {
		s += strings.Split(card.Image, ".")[0] + " "
		if i == 6 || i == 13 {
			s += " "
		}
	}
	return s
}

// Place selected row in the middle
func (d *Deck) deal(row int) {
	log.Printf("before: %s", d)
	var cards []*Card
	switch {
	case row == 0:
		cards = append(cards, d.Cards[7:14]...)
		cards = append(cards, d.Cards[:7]...)
		cards = append(cards, d.Cards[14:]...)

	case row == 1:
		cards = d.Cards

	case row == 2:
		cards = append(cards, d.Cards[:7]...)
		cards = append(cards, d.Cards[14:]...)
		cards = append(cards, d.Cards[7:14]...)
	}
	d.Cards = cards
	log.Printf(" after: %s", d)
}

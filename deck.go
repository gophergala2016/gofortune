package main

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

type Deck struct {
	Cards []*Card
	Row1  []*Card
	Row2  []*Card
	Row3  []*Card
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
	if d.Row1 == nil {
		d.Row1 = d.Cards[:7]
		d.Row2 = d.Cards[7:14]
		d.Row3 = d.Cards[14:]
	}
	s := ""
	for _, card := range d.Row1 {
		s += strings.Split(card.Image, ".")[0] + " "
	}
	s += " "
	for _, card := range d.Row2 {
		s += strings.Split(card.Image, ".")[0] + " "
	}
	s += " "
	for _, card := range d.Row3 {
		s += strings.Split(card.Image, ".")[0] + " "
	}
	return s
}

// Deal deck in three row
func (d *Deck) deal() {
	log.Printf("deal: %s\n", d)
	var row1 []*Card
	var row2 []*Card
	var row3 []*Card
	var cards []*Card
	cards = append(cards, d.Row1...)
	cards = append(cards, d.Row2...)
	cards = append(cards, d.Row3...)
	for i, card := range cards {
		switch i % 3 {
		case 0:
			row1 = append(row1, card)
		case 1:
			row2 = append(row2, card)
		case 2:
			row3 = append(row3, card)
		}
	}
	d.Row1 = row1
	d.Row2 = row2
	d.Row3 = row3
	log.Printf("deal: %s\n\n", d)
}

// Place selected row in the middle
func (d *Deck) placeMiddle(row int) {
	log.Printf("placeMiddle: Selected Row: %d\n", row)
	log.Printf("placeMiddle: %s\n", d)
	switch {
	case row == 1:
		row1 := d.Row1
		d.Row1 = d.Row2
		d.Row2 = row1

	case row == 3:
		row2 := d.Row2
		d.Row2 = d.Row3
		d.Row3 = row2
	}
	log.Printf("placeMiddle: %s\n\n", d)
}

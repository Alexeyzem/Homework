package main

import (
	"strings"
)

type Player struct {
	Inventory   []string
	Backpack    bool
	Target      string
	CurrentRoom Room
}

func (p *Player) PutOn(thing string) string {
	for i, value := range p.CurrentRoom.Furnitures {
		for j, temp := range value.ThingsForPutOn {
			if temp == thing {
				simply := &p.CurrentRoom.Furnitures[i]
				simply.ThingsForPutOn = append(simply.ThingsForPutOn[:j], simply.ThingsForPutOn[j+1:]...)
				if thing == "рюкзак" {
					p.Backpack = true
					p.Target = "надо идти в универ. "
				}
				out := "вы надели: " + thing
				return out
			}
		}
	}
	out := "нельзя надеть: " + thing
	return out
}

func (p *Player) LookAround() string {
	kol := 0
	out := p.CurrentRoom.Message
	for i, value := range p.CurrentRoom.Furnitures {
		if len(value.ThingsForPutOn) == 0 && len(value.ThingsForTake) == 0 {
			kol++
		} else {
			if i == 0 {
				out = out + "на " + value.Decor + ": "
			} else {
				out = out + ", на " + value.Decor + ": "
			}
			for j, temp := range value.ThingsForTake {
				if 0 != j || 0 != i {
					out = out + ", " + temp
				} else {
					out = out + temp
				}
			}
			for j, thing := range value.ThingsForPutOn {
				if j != 0 || len(value.ThingsForTake) != 0 {
					out = out + ", " + thing
				} else {
					out = out + thing
				}
			}
		}
	}

	if kol == len(p.CurrentRoom.Furnitures) {
		out = out + "пустая комната. "
	} else if p.CurrentRoom.NameOfRoom != "кухня" {
		out = out + ". "
	} else {
		out = out + ", " + p.Target
	}
	out = out + "можно пройти - "
	for i, value := range p.CurrentRoom.AdjacentRooms {
		if i != len(p.CurrentRoom.AdjacentRooms)-1 {
			out = out + value.NameOfRoom + ", "
		} else {
			out = out + value.NameOfRoom
		}
	}
	return out
}

func (p *Player) Take(thing string) string {
	for i, value := range p.CurrentRoom.Furnitures {
		for j, temp := range value.ThingsForTake {
			if thing == temp && p.Backpack {
				p.Inventory = append(p.Inventory, thing)
				simply := &p.CurrentRoom.Furnitures[i]
				simply.ThingsForTake = append(simply.ThingsForTake[:j], simply.ThingsForTake[j+1:]...)
				return "предмет добавлен в инвентарь: " + thing
			} else if !p.Backpack {
				//fmt.Println("некуда класть")
				return "некуда класть"
			}
		}
	}
	return "нет такого"
}

func (p *Player) Use(ThingFromInv, ThingAround string) string {
	var out string
	for _, value := range p.Inventory {
		if ThingFromInv == value {
			switch ThingFromInv {
			case "ключи":
				{
					switch ThingAround {
					case "дверь":
						{
							p.CurrentRoom.OpenDoor = true
							for _, temp := range p.CurrentRoom.AdjacentRooms {
								if !temp.OpenDoor {
									temp.OpenDoor = !temp.OpenDoor
									break
								}
							}
							out = out + "дверь открыта"
							return out
						}
					default:
						out = out + "не к чему применить"
						return out
					}
				}

			}
		}
	}
	out = out + "нет предмета в инвентаре - " + ThingFromInv
	return out
}
func (p *Player) Go(GoTo string) string {
	var out string
	for i, value := range p.CurrentRoom.AdjacentRooms {
		if value.NameOfRoom == GoTo {
			if p.CurrentRoom.OpenDoor || p.CurrentRoom.AdjacentRooms[i].OpenDoor {
				p.CurrentRoom = *p.CurrentRoom.AdjacentRooms[i]
				out = out + p.CurrentRoom.MessageForGo + "можно пройти - "
				for j, temp := range p.CurrentRoom.AdjacentRooms {
					if j != len(p.CurrentRoom.AdjacentRooms)-1 {
						out = out + temp.NameOfRoom + ", "
					} else {
						out = out + temp.NameOfRoom
					}
				}
				return out
			} else {
				out = out + "дверь закрыта"
				return out
			}
		}
	}
	out = out + "нет пути в " + GoTo
	return out
}

type Room struct {
	NameOfRoom    string
	Furnitures    []Furniture
	AdjacentRooms []*Room
	Message       string
	OpenDoor      bool
	MessageForGo  string
}

type Furniture struct {
	Decor          string
	ThingsForPutOn []string
	ThingsForTake  []string
}

var P Player
var (
	RHouse   Room
	RStreet  Room
	RMyRoom  Room
	RHallway Room
	RKitchen Room
)

func initGame() {

	var neighborsMyRoom []*Room
	neighborsMyRoom = append(neighborsMyRoom, &RHallway)
	var MyRoomFurnit = []Furniture{
		{"столе", nil, []string{"ключи", "конспекты"}},
		{"стуле", []string{"рюкзак"}, nil},
	}
	RMyRoom = Room{
		Furnitures:    MyRoomFurnit,
		AdjacentRooms: neighborsMyRoom,
		NameOfRoom:    "комната",
		OpenDoor:      true,
		MessageForGo:  "ты в своей комнате. ",
		Message:       "",
	}
	var neighborsKitchen []*Room
	neighborsKitchen = append(neighborsKitchen, &RHallway)
	var kitchenFurnit = []Furniture{
		{"столе", []string{"чай"}, nil},
	}

	RKitchen = Room{
		Furnitures:    kitchenFurnit,
		AdjacentRooms: neighborsKitchen,
		NameOfRoom:    "кухня",
		OpenDoor:      true,
		MessageForGo:  "кухня, ничего интересного. ",
		Message:       "ты находишься на кухне, ",
	}
	var neighborsHallway = []*Room{&RKitchen, &RMyRoom, &RStreet}
	RHallway = Room{
		Furnitures:    nil,
		AdjacentRooms: neighborsHallway,
		NameOfRoom:    "коридор",
		OpenDoor:      false,
		MessageForGo:  "ничего интересного. ",
		Message:       "ты находишься в коридоре, ",
	}
	var neighborsStreet = []*Room{&RHouse}
	RStreet = Room{
		NameOfRoom:    "улица",
		MessageForGo:  "на улице весна. ",
		OpenDoor:      false,
		AdjacentRooms: neighborsStreet,
	}
	RHouse = RHallway
	RHouse.NameOfRoom = "домой"
	P = Player{
		Inventory:   nil,
		Backpack:    false,
		Target:      "надо собрать рюкзак и идти в универ. ",
		CurrentRoom: RKitchen,
	}
}
func handleCommand(str string) string {
	strM := strings.Split(str, " ")
	switch strM[0] {
	case "осмотреться":
		return P.LookAround()
	case "взять":
		return P.Take(strM[1])
	case "идти":
		return P.Go(strM[1])
	case "надеть":
		return P.PutOn(strM[1])
	case "применить":
		return P.Use(strM[1], strM[2])
	case "посмотреть":
	default:
		return "неизвестная команда"
	}
	return ""
}

func main() {
}

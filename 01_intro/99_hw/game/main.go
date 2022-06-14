package main

import (
	"fmt"
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
				p.CurrentRoom.Furnitures[i].ThingsForPutOn = append(p.CurrentRoom.Furnitures[i].ThingsForPutOn[:j], p.CurrentRoom.Furnitures[i].ThingsForPutOn[j+1:]...)
				if thing == "рюкзак" {
					p.Backpack = true
					p.Target = "надо идти в универ. "
				}
				out := "вы надели: " + thing
				//fmt.Println("Вы надели:", thing)
				return out
			}
		}
	}
	out := "нельзя надеть: " + thing
	//fmt.Println("нельзя надеть: ", thing)
	return out
}

func (p *Player) LookAround() string {
	kol := 0
	out := p.CurrentRoom.Message
	//fmt.Print(p.CurrentRoom.Message)
	for i, value := range p.CurrentRoom.Furnitures {
		if len(value.ThingsForPutOn) == 0 && len(value.ThingsForTake) == 0 {
			kol++
		} else {
			if i == 0 {
				out = out + "на " + value.Decor + ": "
				//fmt.Print("на ", value.Decor, ": ")
			} else {
				out = out + ", на " + value.Decor + ": "
				//fmt.Print(", на ", value.Decor, ": ")
			}
			for j, temp := range value.ThingsForTake {
				if 0 != j || 0 != i {
					out = out + ", " + temp
					//fmt.Print(", ", temp)
				} else {
					out = out + temp
					//fmt.Print(temp)
				}
			}
			for j, thing := range value.ThingsForPutOn {
				if j != 0 || len(value.ThingsForTake) != 0 {
					out = out + ", " + thing
					//fmt.Print(", ", thing)
				} else {
					out = out + thing
					//fmt.Print(thing)
				}
			}
		}
	}

	if kol == len(p.CurrentRoom.Furnitures) {
		out = out + "пустая комната. "
		//fmt.Print("пустая комната. ")
	} else if p.CurrentRoom.NameOfRoom != "кухня" {
		out = out + ". "
		//fmt.Print(". ")
	} else {
		out = out + ", " + p.Target
		//fmt.Print(", ", p.Target)
	}
	out = out + "можно пройти - "
	//fmt.Print("можно пройти - ")
	for i, value := range p.CurrentRoom.AdjacentRooms {
		if i != len(p.CurrentRoom.AdjacentRooms)-1 {
			out = out + value.NameOfRoom + ", "
			//fmt.Print(value.NameOfRoom, ", ")
		} else {
			out = out + value.NameOfRoom
			//	fmt.Println(value.NameOfRoom)
		}
	}
	return out
}

func (p *Player) Take(thing string) string {
	for i, value := range p.CurrentRoom.Furnitures {
		for j, temp := range value.ThingsForTake {
			if thing == temp && p.Backpack {
				p.Inventory = append(p.Inventory, thing)
				p.CurrentRoom.Furnitures[i].ThingsForTake = append(p.CurrentRoom.Furnitures[i].ThingsForTake[:j], p.CurrentRoom.Furnitures[i].ThingsForTake[j+1:]...)
				//fmt.Println("предмет добавлен в инвентарь:", thing)
				return "предмет добавлен в инвентарь: " + thing
			} else if !p.Backpack {
				//fmt.Println("некуда класть")
				return "некуда класть"
			}
		}
	}
	//fmt.Println("нет такого")
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
							//fmt.Println("дверь открыта")
							return out
						}
					default:
						out = out + "не к чему применить"
						fmt.Println("не к чему применить")
						return out
					}
				}

			}
		}
	}
	out = out + "нет предмета в инвентаре - " + ThingFromInv
	//fmt.Println("нет предмета в инвентаре -", ThingFromInv)
	return out
}
func (p *Player) Go(GoTo string) string {
	var out string
	for i, value := range p.CurrentRoom.AdjacentRooms {
		if value.NameOfRoom == GoTo {
			if p.CurrentRoom.OpenDoor || p.CurrentRoom.AdjacentRooms[i].OpenDoor {
				p.CurrentRoom = *p.CurrentRoom.AdjacentRooms[i]
				out = out + p.CurrentRoom.MessageForGo + "можно пройти - "
				//fmt.Printf("%sможно пройти - ", p.CurrentRoom.MessageForGo)
				for j, temp := range p.CurrentRoom.AdjacentRooms {
					if j != len(p.CurrentRoom.AdjacentRooms)-1 {
						out = out + temp.NameOfRoom + ", "
						//	fmt.Print(temp.NameOfRoom, ", ")
					} else {
						out = out + temp.NameOfRoom
						//	fmt.Println(temp.NameOfRoom)
					}
				}
				return out
			} else {
				out = out + "дверь закрыта"
				//	fmt.Println("дверь закрыта")
				return out
			}
		}
	}
	out = out + "нет пути в " + GoTo
	//fmt.Println("нет пути в", GoTo)
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
	P.Inventory = nil
	P.Target = ""
	P.Backpack = false
	//VLad. Напоминалочка.
	RMyRoom = Room{
		Furnitures:    nil,
		AdjacentRooms: nil,
	}
	RKitchen.Furnitures = nil
	RKitchen.AdjacentRooms = nil
	RHallway.AdjacentRooms = nil
	RStreet.AdjacentRooms = nil

	P.Target = "надо собрать рюкзак и идти в универ. "
	P.Backpack = false
	RKitchen.NameOfRoom = "кухня"
	RKitchen.OpenDoor = true
	RKitchen.AdjacentRooms = append(RKitchen.AdjacentRooms, &RHallway)
	RKitchen.MessageForGo = "кухня, ничего интересного. "
	RKitchen.Message = "ты находишься на кухне, "
	var kitchenFurnit Furniture
	kitchenFurnit.Decor = "столе"
	kitchenFurnit.ThingsForTake = append(kitchenFurnit.ThingsForTake, "чай")
	RKitchen.Furnitures = append(RKitchen.Furnitures, kitchenFurnit)
	P.CurrentRoom = RKitchen
	RHallway.AdjacentRooms = append(RHallway.AdjacentRooms, &RKitchen, &RMyRoom, &RStreet)
	RHallway.OpenDoor = false
	RHallway.NameOfRoom = "коридор"
	RHallway.Message = "ты находишься в коридоре, "
	RHallway.MessageForGo = "ничего интересного. "
	RStreet.OpenDoor = false
	RStreet.NameOfRoom = "улица"
	RStreet.MessageForGo = "на улице весна. "
	RHouse = RHallway
	RHouse.NameOfRoom = "домой"
	RStreet.AdjacentRooms = append(RStreet.AdjacentRooms, &RHouse)
	RMyRoom.NameOfRoom = "комната"
	RMyRoom.AdjacentRooms = append(RMyRoom.AdjacentRooms, &RHallway)
	RMyRoom.OpenDoor = true
	RMyRoom.MessageForGo = "ты в своей комнате. "
	RMyRoom.Message = ""
	var MyRoomFurnit = []Furniture{
		{"столе", nil, []string{"ключи", "конспекты"}},
		{"стуле", []string{"рюкзак"}, nil},
	}

	RMyRoom.Furnitures = append(MyRoomFurnit)
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
	//for {
	//	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//	if err != nil {
	//		fmt.Println("Nothing")
	//	}
	//	str = strings.ReplaceAll(str, "\r\n", "")
	//
	//}
}

//urlDownload := "https://github.com/semyon-dev/stepik-go/blob/master/work_with_files/task_sep_1/task.data"     чтение файлов с гита
//resp, err := http.Get(urlDownload)																			чтение файлов с гита
//if err != nil {																								чтение файлов с гита
//	return																										чтение файлов с гита
//}																												чтение файлов с гита
//defer resp.Body.Close()																						чтение файлов с гита
//file, err := os.Create("test")																				чтение файлов с гита
//if err != nil {																								чтение файлов с гита
//	return																										чтение файлов с гита
//}																												чтение файлов с гита
//defer file.Close()
//_, err = io.Copy(file, resp.Body)

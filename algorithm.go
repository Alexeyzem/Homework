package main

import (
	"fmt"
	"sort"
)

type List struct {
	number int
	next   *List
}

type FirstList struct {
	len  int
	head *List
}

func initList() *FirstList {
	return &FirstList{}
}

func (f *FirstList) AddFront(number int) {
	list := &List{
		number: number,
	}
	if f.head == nil {
		f.head = list
	} else {
		list.next = f.head
		f.head = list
	}
	f.len++
	return
}

func (f *FirstList) AddBack(number int) {
	list := &List{
		number: number,
	}
	if f.head == nil {
		f.head = list
	} else {
		current := f.head
		for current.next != nil {
			current = current.next
		}
		current.next = list
	}
	f.len++
	return
}

func (f *FirstList) RemoveFront() error {
	if f.head == nil {
		return fmt.Errorf("RemoveFront: List is empty")
	}
	f.head = f.head.next
	f.len--
	return nil
}

func (f *FirstList) RemoveBack() error {
	if f.head == nil {
		return fmt.Errorf("RemoveBack: List is Empty")
	}
	var prev *List
	current := f.head
	for current.next != nil {
		prev = current
		current = current.next
	}
	if prev != nil {
		prev.next = nil
	} else {
		f.head = nil
	}
	f.len--
	return nil
}

func (f *FirstList) Front() (int, error) {
	if f.head == nil {
		return 0, fmt.Errorf("single List is empty")
	}
	return f.head.number, nil
}

func (f *FirstList) Size() int {
	return f.len
}

func (f *FirstList) Traverse() error {
	if f.head == nil {
		return fmt.Errorf("TranverseError: List is empty")
	}
	current := f.head
	for current != nil {
		fmt.Println(current.number)
		current = current.next
	}
	return nil
}

func BinarySearchForInt(item int, list []int) bool {
	sort.Ints(list)
	low := 0
	high := len(list) - 1
	for low <= high {
		var mid int
		mid = (low + high) / 2
		guess := list[mid]
		if guess == item {
			return true
		} else if guess > item {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false
}

func BinarySearchForString(item string, list []string) bool {
	sort.Strings(list)
	low := 0
	high := len(list) - 1
	for low <= high {
		var mid int
		mid = (low + high) / 2
		guess := list[mid]
		if guess == item {
			return true
		} else if guess > item {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false
}

func Sum(arr []int) int {
	if len(arr) == 0 {
		return 0
	} else {
		arr1 := append(arr[1:])
		return arr[0] + Sum(arr1)
	}
}

func Quantity(temle List) int {
	if temle.next == nil {
		return 1
	} else {
		return 1 + Quantity(*temle.next)
	}
}

func Max(temple List) int {
	if temple.next == nil {
		return temple.number
	}
	if temple.number > Max(*temple.next) {
		return temple.number
	} else {
		return Max(*temple.next)
	}
}

func FastSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	} else {
		mid := (len(arr) - 1) / 2
		pivot := arr[mid]
		var left []int
		var right []int
		for i, value := range arr {
			if i != mid {
				if value <= pivot {
					left = append(left, value)
				} else {
					right = append(right, value)
				}
			}
		}
		left = FastSort(left)
		right = FastSort(right)
		var output []int
		output = append(append(left, arr[mid]), right...)
		return output
	}
}

func Deikstra(graph map[string]map[string]int, costs map[string]int, parents map[string]string) int {
	var processed []string
	min := FindLowest(costs, processed)
	for min != "" {
		cost := costs[min]
		neighbors := graph[min]
		for key, value := range neighbors {
			newCost := cost + value
			if costs[key] > newCost {
				costs[key] = newCost
				parents[key] = min
			}
		}
		processed = append(processed, min)
		min = FindLowest(costs, processed)
	}
	return costs["fin"]
}

func FindLowest(m map[string]int, processed []string) string {
	lowCost := 922337203685477580
	var output string
	for key, value := range m {
		if value < lowCost && !BinarySearchForString(key, processed) {
			lowCost = value
			output = key
			processed = append(processed, key)
		}
	}
	return output
}

func GreedyAlgorithm(stations map[string][]string, statesNeded map[string]bool) []string {
	var finalStation []string
	for len(statesNeded) != 0 {
		var bestStation string
		var statesCovered []string
		for station, statesForStation := range stations {
			var covered []string
			for _, item := range statesForStation {
				if statesNeded[item] {
					covered = append(covered, item)
				}
			}
			if len(covered) > len(statesCovered) {
				bestStation = station
				statesCovered = covered
			}
		}
		for i := range statesCovered {
			delete(statesNeded, statesCovered[i])
		}
		finalStation = append(finalStation, bestStation)
		delete(stations, bestStation)
	}
	return finalStation
}

func main() {
	states := map[string]bool{
		"mt": true,
		"id": true,
		"nv": true,
		"ut": true,
		"wa": true,
		"or": true,
		"ca": true,
		"az": true,
	}
	stations := map[string][]string{
		"kone":   {"id", "nv", "ut"},
		"ktwo":   {"wa", "id", "mt"},
		"kthree": {"or", "nv", "ca"},
		"kfour":  {"nv", "ut"},
		"kfive":  {"ca", "az"},
	}
	fmt.Print(GreedyAlgorithm(stations, states))
}

package main

import (
	"testing"
)

const infinity int = 922337203685477580

type TestCaseForBinSeaForInt struct {
	item   int
	list   []int
	result bool
}

func TestBinarySearchForInt(t *testing.T) {
	casesForInt := []TestCaseForBinSeaForInt{
		{
			item:   5,
			list:   []int{1, 2, 3, 4, 5, 6, 10, 23, 45, 10},
			result: true,
		},
		{
			item:   0,
			list:   []int{1, 2, 34, 1, 4, 5, 2, 5, 9},
			result: false,
		},
		{
			item:   10,
			list:   []int{10, 2, 3, 1, 7, 8, 0},
			result: true,
		},
		{
			item:   11,
			list:   []int{1, 2, 2, 1, 3, 1, 11},
			result: true,
		},
		{
			item:   0,
			list:   []int{},
			result: false,
		},
	}
	for caseNum, item := range casesForInt {
		result := BinarySearchForInt(item.item, item.list)
		if result != item.result {
			t.Errorf("[%d] wrong result in BinarySearchForInt, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
}

type TestCaseForBinSeaForStr struct {
	item   string
	list   []string
	result bool
}

func TestBinarySearchForString(t *testing.T) {

	casesForStr := []TestCaseForBinSeaForStr{
		{
			item:   "alex",
			list:   []string{"vlad", "alexey", "alex", "alexander", "andrey"},
			result: true,
		},
		{
			item:   "0",
			list:   []string{"1", "number", "seven", "zero", "you"},
			result: false,
		},
		{
			item:   "I",
			list:   []string{"I", "you", "me", "phone"},
			result: true,
		},
		{
			item:   "mother",
			list:   []string{"father", "sister", "brother", "aunt", "mother"},
			result: true,
		},
		{
			item:   "",
			list:   []string{},
			result: false,
		},
	}
	for caseNum, item := range casesForStr {
		result := BinarySearchForString(item.item, item.list)
		if result != item.result {
			t.Errorf("[%d] wrong result in BinarySearchForString, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
}

type TestCaseForSum struct {
	array  []int
	result int
}

func TestSum(t *testing.T) {

	casesForSum := []TestCaseForSum{
		{
			array:  []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 9},
			result: 55,
		},
		{
			array:  []int{},
			result: 0,
		},
	}
	for caseNum, item := range casesForSum {
		result := Sum(item.array)
		if result != item.result {
			t.Errorf("[%d] wrong result in Sum, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
}

type TestCaseForQuantity struct {
	data   List
	result int
}

type TestCaseForMax struct {
	data   List
	result int
}

func TestMaxAndQuantity(t *testing.T) {
	singleListFirst := initList()
	singleListFirst.AddFront(1)
	singleListFirst.AddFront(2)
	singleListFirst.AddBack(15)
	singleListFirst.AddBack(3)
	singleListFirst.AddFront(10)
	singleListSecond := initList()
	singleListSecond.AddBack(1)
	singleListSecond.AddBack(1999)
	singleListSecond.AddBack(3)
	singleListSecond.AddBack(98)
	singleListSecond.AddBack(90)
	singleListSecond.AddBack(19)
	singleListSecond.AddBack(30)
	singleListSecond.AddBack(90)
	singleListSecond.AddBack(0)
	singleListThird := initList()
	singleListThird.AddBack(5)
	casesForQuantity := []TestCaseForQuantity{
		{
			data:   *singleListFirst.head,
			result: 5,
		},
		{
			data:   *singleListSecond.head,
			result: 9,
		},
		{
			data:   *singleListThird.head,
			result: 1,
		},
	}
	for caseNum, item := range casesForQuantity {
		result := Quantity(item.data)
		if result != item.result {
			t.Errorf("[%d] wrong result in Quantity, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
	CasesForMax := []TestCaseForMax{
		{
			data:   *singleListFirst.head,
			result: 15,
		},
		{
			data:   *singleListSecond.head,
			result: 1999,
		},
		{
			data:   *singleListThird.head,
			result: 5,
		},
	}
	for caseNum, item := range CasesForMax {
		result := Max(item.data)
		if result != item.result {
			t.Errorf("[%d] wrong result in Quantity, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
}

type TestCaseForFastSort struct {
	input  []int
	result []int
}

func TestFastSort(t *testing.T) {
	Cases := []TestCaseForFastSort{
		{
			input:  []int{1, 8, 3, 4, 1, 4, 7, 4, 3, 34565, 4322, 3454, 32},
			result: []int{1, 1, 3, 3, 4, 4, 4, 7, 8, 32, 3454, 4322, 34565},
		},
		{
			input:  []int{},
			result: []int{},
		},
		{
			input:  []int{5},
			result: []int{5},
		},
	}
	for caseNum, item := range Cases {
		result := FastSort(item.input)
		for i := range result {
			if result[i] != item.result[i] {
				t.Errorf("[%d] wrong result in Quantity, expected %#v, got %#v", caseNum, item.result, result)
			}
		}
	}
}

type TestCaseForDeikstra struct {
	graph   map[string]map[string]int
	parents map[string]string
	costs   map[string]int
	result  int
}

func TestDeikstra(t *testing.T) {
	graph1 := make(map[string]map[string]int)
	graph1["start"] = map[string]int{"a": 6, "b": 2}
	graph1["b"] = map[string]int{"a": 3, "fin": 5}
	graph1["a"] = map[string]int{"fin": 1}
	graph1["fin"] = map[string]int{}
	costs1 := make(map[string]int)
	costs1 = map[string]int{"a": 6, "b": 2, "fin": infinity}
	parents1 := make(map[string]string)
	parents1 = map[string]string{"a": "start", "b": "start", "fin": ""}
	graph2 := make(map[string]map[string]int)
	graph2["start"] = map[string]int{"a": 5, "b": 2}
	graph2["b"] = map[string]int{"a": 5, "c": 7}
	graph2["a"] = map[string]int{"c": 2, "d": 4}
	graph2["c"] = map[string]int{"fin": 1}
	graph2["d"] = map[string]int{"c": 6, "fin": 3}
	costs2 := make(map[string]int)
	costs2 = map[string]int{"a": 5, "b": 2, "c": infinity, "d": infinity, "fin": infinity}
	parents2 := make(map[string]string)
	parents2 = map[string]string{"a": "start", "b": "start", "c": "", "d": "", "fin": ""}
	graph3 := make(map[string]map[string]int)
	graph3["start"] = map[string]int{"a": 10}
	graph3["b"] = map[string]int{"fin": 30, "c": 1}
	graph3["a"] = map[string]int{"b": 20}
	graph3["c"] = map[string]int{"a": 1}
	costs3 := make(map[string]int)
	costs3 = map[string]int{"a": 10, "b": infinity, "c": infinity, "fin": infinity}
	parents3 := make(map[string]string)
	parents3 = map[string]string{"a": "start", "b": "", "c": "", "fin": ""}
	Cases := []TestCaseForDeikstra{
		{
			graph:   graph1,
			costs:   costs1,
			parents: parents1,
			result:  6,
		},
		{
			graph:   graph2,
			costs:   costs2,
			parents: parents2,
			result:  8,
		},
		{
			graph:   graph3,
			costs:   costs3,
			parents: parents3,
			result:  60,
		},
	}
	for caseNum, item := range Cases {
		result := Deikstra(item.graph, item.costs, item.parents)
		if result != item.result {
			t.Errorf("[%d] wrong result in Quantity, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
}

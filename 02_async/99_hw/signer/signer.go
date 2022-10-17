package main

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := sync.WaitGroup{}
	in := make(chan interface{})
	for _, joba := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		inWorker := in
		go func(jobaWork job) {
			defer wg.Done()
			jobaWork(inWorker, out)
			close(out)
		}(joba)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for value := range in {
		wg.Add(1)
		go func(value interface{}) {
			defer wg.Done()
			right := make(chan string)
			dataInt, ok := value.(int)
			if ok {
				data := strconv.Itoa(dataInt)
				go func(data string) {
					mu.Lock()
					temp := DataSignerMd5(data)
					mu.Unlock()
					right <- DataSignerCrc32(temp)
				}(data)
				out <- DataSignerCrc32(data) + "~" + <-right
			} else {
				log.Print("ошибка перевода в SingleHash")
			}
		}(value)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	internalWg := sync.WaitGroup{}
	externalWg := sync.WaitGroup{}
	for data := range in {
		externalWg.Add(1)
		go func(data interface{}) {
			defer externalWg.Done()
			intermediateData := make([]string, 6, 6)
			value, ok := (data).(string)
			if ok {
				for th := 0; th < 6; th++ {
					internalWg.Add(1)
					go func(th int) {
						intermediateData[th] = DataSignerCrc32(strconv.Itoa(th) + value)
						internalWg.Done()
					}(th)
				}
			} else {
				log.Print("ошибка перевода в MultiHash")
			}
			internalWg.Wait()
			output := strings.Join(intermediateData, "")
			out <- output
		}(data)
	}
	externalWg.Wait()
}
func CombineResults(in, out chan interface{}) {
	var result string
	var resultArray []string
	for data := range in {
		temp, ok := (data).(string)
		if ok {
			resultArray = append(resultArray, temp)
		} else {
			log.Print("ошибка перевода при CombineResults")
		}
	}
	sort.Strings(resultArray)
	result = strings.Join(resultArray, "_")
	out <- result
}
func main() {
}

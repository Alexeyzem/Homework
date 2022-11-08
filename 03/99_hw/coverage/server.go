package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var FileName = "dataset.xml"

type Row struct {
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}
type Root struct {
	List []Row `xml:"row"`
}

const (
	writeError    = "can not answer-%w"
	internalError = "internal server error"
)

func initMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["0"] = nil
	out["1"] = nil
	out["MyToken"] = nil
	return out
}
func SearchServer(w http.ResponseWriter, r *http.Request) {
	tokenMap := initMap()
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.WriteHeader(500)
		_, errW := w.Write([]byte(internalError))
		if err != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
	}
	token := r.Header.Get("AccessToken")
	if _, ok := tokenMap[token]; !ok {
		w.WriteHeader(401)
		_, errW := w.Write([]byte("wrong token, please log in"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	file, err := os.Open(FileName)
	if err != nil {
		log.Println(fmt.Errorf("can not open file with XML data-%w", err))
		w.WriteHeader(500)
		_, errW := w.Write([]byte(internalError + "\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(fmt.Errorf("can not close file-%w", err))
		}
	}()
	var dataXML []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dataXML = append(dataXML, scanner.Bytes()...)
	}
	data := new(Root)
	err = xml.Unmarshal(dataXML, &data)
	if err != nil {
		log.Print(fmt.Errorf("can not unmarshal xml file - %w", err))
		w.WriteHeader(500)
		_, errW := w.Write([]byte(internalError + "\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	switch orderField {
	case "Id":
		if orderBy == 1 || orderBy == -1 {
			sort.Slice(data.List, func(i, j int) (less bool) {
				return orderBy*data.List[i].Id < orderBy*data.List[j].Id
			})
		} else if orderBy != 0 {
			log.Println(fmt.Errorf("error - orderBy invalid"))
			w.WriteHeader(400)
			_, errW := w.Write([]byte("invalid order_by\n"))
			if errW != nil {
				log.Println(fmt.Errorf(writeError, errW))
			}
			return
		}
	case "Name":
		if orderBy == 1 {
			sort.Slice(data.List, func(i, j int) (less bool) {
				return data.List[i].FirstName+" "+data.List[i].LastName < data.List[j].FirstName+" "+data.List[j].LastName
			})
		} else if orderBy == -1 {
			sort.Slice(data.List, func(i, j int) (less bool) {
				return data.List[i].FirstName+" "+data.List[i].LastName > data.List[j].FirstName+" "+data.List[j].LastName
			})
		} else if orderBy != 0 {
			log.Println(fmt.Errorf("error - orderBy invalid"))
			w.WriteHeader(400)
			_, errW := w.Write([]byte("invalid order_by \n"))
			if errW != nil {
				log.Println(fmt.Errorf(writeError, errW))
			}
			return
		}
	case "Age":
		if orderBy == 1 || orderBy == -1 {
			sort.Slice(data.List, func(i, j int) (less bool) {
				return orderBy*data.List[i].Age < orderBy*data.List[j].Age
			})
		} else if orderBy != 0 {
			log.Println(fmt.Errorf("error - orderBy invalid"))
			w.WriteHeader(400)
			_, errW := w.Write([]byte("invalid order_by \n"))
			if errW != nil {
				log.Println(fmt.Errorf(writeError, errW))
			}
			return
		}
	default:
		log.Println(fmt.Errorf("error - %s", `OrderField invalid`))
		w.WriteHeader(400)
		_, errW := w.Write([]byte("OrderField invalid\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	var output []Row
	first, err := strconv.Atoi(offset)
	if err != nil {
		log.Println(fmt.Errorf("offset invalid - %w", err))
		w.WriteHeader(400)
		_, errW := w.Write([]byte("offset invalid\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	lim, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(fmt.Errorf("limit invalid - %w", err))
		w.WriteHeader(400)
		_, errW := w.Write([]byte("limit invalid\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	if query == "" {
		output = append(output, data.List[first-1:lim+first-1]...)
	} else {
		i := 1
		for _, value := range data.List {
			if strings.Contains(value.FirstName+" "+value.LastName, query) || strings.Contains(value.About, query) {
				if i >= first && len(output) < lim {
					output = append(output, value)
				}
				i++
			}
		}
	}
	var out []User
	var temp User
	for _, value := range output {
		temp.Age = value.Age
		temp.About = value.About
		temp.ID = value.Id
		temp.Gender = value.Gender
		temp.Name = value.FirstName + " " + value.LastName
		out = append(out, temp)
	}
	if len(out) == 0 {
		w.WriteHeader(412)
		_, errW := w.Write([]byte("The specified users are not in the database. Try again, MOTHER-FUCKA.\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	w.WriteHeader(200)
	outJSON, err := json.Marshal(out)
	if err != nil {
		log.Println(fmt.Errorf("can not marshal to JSOn - %w", err))
		w.WriteHeader(500)
		_, errW := w.Write([]byte(internalError + "\n"))
		if errW != nil {
			log.Println(fmt.Errorf(writeError, errW))
		}
		return
	}
	_, errW := w.Write(outJSON)
	if errW != nil {
		log.Println(fmt.Errorf(writeError, errW))
	}
}
func main() {
	http.HandleFunc("/", SearchServer)
	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}

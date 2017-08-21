package serviceClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Inputobj struct {
	Msg string
}

type client struct {
	consul *consul.Client
}

func NewClient() (client, error) {
	var Client client
	config := consul.DefaultConfig()
	config.Address = "localhost:8500"
	c, err := consul.NewClient(config)
	if err != nil {
		return Client, err
	}
	Client.consul = c
	return Client, nil
}

func (s *client) findService(service string) (string, error) {
	passingOnly := true
	addrs, _, err := s.consul.Health().Service(service, "", passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return "", fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return "", err
	}
	rand.Seed(time.Now().Unix())
	selectedService := addrs[rand.Intn(len(addrs))]
	address := selectedService.Service.Address + ":" + strconv.Itoa(selectedService.Service.Port)
	return address, nil
}

func (s *client) Send(service string, functionname string, msg Inputobj) {
	address, err := s.findService(service)
	if err != nil {
		return
	}
	address = "http://" + address + "/" + service + "/" + functionname
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg.Msg)
	res, _ := http.Post(address, "application/json; charset=utf-8", b)
	decoder := json.NewDecoder(res.Body)
	var inputjson Inputobj
	err = decoder.Decode(&inputjson.Msg)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(inputjson)
}
/*
func main() {
	Client, _ := NewClient()
        jsonObjToSend := map[string]int{"int1": 5, "int2": 8}
        mapB, err := json.Marshal(jsonObjToSend)
        if err!= nil {
             panic(err)
        }
        toSend := inputobj{msg:string(mapB)}
	Client.send("hi", "hello", toSend)
}
*/

package main

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

type inputobj struct {
	msg interface{}
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

func (s *client) send(service string, functionname string, msg inputobj) {
	address, err := s.findService(service)
	if err != nil {
		return
	}
	address = "http://" + address + "/" + service + "/" + functionname
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg.msg)
	res, _ := http.Post(address, "application/json; charset=utf-8", b)
	decoder := json.NewDecoder(res.Body)
	var inputjson inputobj
	err = decoder.Decode(&inputjson.msg)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Println(inputjson)
}

func main() {
	Client, _ := NewClient()
	Client.send("hi", "hello", inputobj{msg: map[string]map[string]int{"msg": {"int1": 5, "int2": 8}}})
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

type result func(error, resultobj)

type action func(inputobj, result)

type service struct {
	name   string
	port   int
	consul *consul.Client
}

type resultobj struct {
	result interface{}
}

type inputobj struct {
	msg interface{}
}

type requestParserStruct struct {
	rw       http.ResponseWriter
	req      *http.Request
	callback action
}

func (r *requestParserStruct) requestParser(rw http.ResponseWriter, req *http.Request) {
	r.rw = rw
	r.req = req
	decoder := json.NewDecoder(req.Body)
	var inputjson inputobj
	err := decoder.Decode(&inputjson.msg)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	r.callback(inputjson, func(err error, resultjson resultobj) {
		json.NewEncoder(rw).Encode(resultjson.result)
	})
}

func statusUpdate(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte{'h', 'e', 'a', 'l', 't', 'h', 'y'})
}

func getMyIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func registerService(name string, port int) (service, error) {
	fmt.Println("registerService:", name)
	var Service service
	Service.name = name
	Service.port = port
	config := consul.DefaultConfig()
	config.Address = "localhost:8500"
	c, err := consul.NewClient(config)
	if err != nil {
		return Service, err
	}
	Service.consul = c
	UUID := uuid.New().String()
	serviceID := name + "-" + UUID
	hostname := getMyIp()
	servicechecks := &consul.AgentServiceCheck{
		Script:   "curl --connect-timeout 10 " + "http" + "://" + hostname + ":" + strconv.Itoa(port) + "/" + serviceID + "/status",
		Interval: "10s",
		Timeout:  "8s",
		TTL:      "",
		HTTP:     "http" + "://" + hostname + ":" + strconv.Itoa(port) + "/" + serviceID + "/status",
		Status:   "passing",
		DeregisterCriticalServiceAfter: "10s",
	}
	reg := &consul.AgentServiceRegistration{
		ID:                serviceID,
		Name:              name,
		Port:              port,
		Tags:              []string{name},
		Address:           hostname,
		EnableTagOverride: true,
		Check:             servicechecks,
		Checks:            consul.AgentServiceChecks{},
	}
	err = Service.consul.Agent().ServiceRegister(reg)
	if err != nil {
		panic(err)
		return Service, err
	}
	http.HandleFunc("/"+serviceID+"/status", statusUpdate)
	return Service, nil
}

func (s *service) addFunction(name string, callback action) {
	fmt.Println("addFunction:", name)
	var request requestParserStruct
	request.callback = callback
	http.HandleFunc("/"+s.name+"/"+name, request.requestParser)
}

func (s *service) start() {
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

func main() {
	Service, _ := registerService("hi", 1234)
	fmt.Println(Service)
	Service.addFunction("hello", func(inputjson inputobj, done result) {
		fmt.Println(reflect.TypeOf(inputjson.msg))
		int1 := inputjson.msg.(map[string]interface{})["msg"].(map[string]interface{})["int1"]
		int2 := inputjson.msg.(map[string]interface{})["msg"].(map[string]interface{})["int2"]
		fmt.Println(reflect.TypeOf(int1), reflect.TypeOf(int2))
		sum := int1.(float64) + int2.(float64)
		done(nil, resultobj{result: map[string]float64{"sum": sum}})
	})
	Service.start()
}

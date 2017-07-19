package main

import (
        "encoding/json"
        "fmt"
        consul "github.com/hashicorp/consul/api"
        "net/http"
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
        fmt.Println(inputjson)
        r.callback(inputjson, func(err error, resultjson resultobj) {
                fmt.Println(resultjson)
        })

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
        reg := &consul.AgentServiceRegistration{
                ID:                name,
                Name:              name,
                Port:              port,
                Tags:              []string{name},
                Address:           "",
                EnableTagOverride: true,
        }
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
        Service.addFunction("hello", func(json inputobj, done result) {
                fmt.Println(json)
                done(nil, resultobj{result: "Hi"})
        })
        Service.start()
}

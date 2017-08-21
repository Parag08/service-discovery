package serviceRegistry

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"net"
	"net/http"
	"strconv"
)

type Result func(error, Resultobj)

type Action func(Inputobj, Result)

type service struct {
	name   string
	port   int
	consul *consul.Client
}
/*
type functionsInput struct {
        Int1  int   `json:"int1"`
        Int2  int   `json:"int2"`
}
*/
type Resultobj struct {
	Result string
}

type Inputobj struct {
	Msg string
}

type requestParserStruct struct {
	rw       http.ResponseWriter
	req      *http.Request
	callback Action
}

func (r *requestParserStruct) requestParser(rw http.ResponseWriter, req *http.Request) {
	r.rw = rw
	r.req = req
	decoder := json.NewDecoder(req.Body)
	var inputjson Inputobj
	err := decoder.Decode(&inputjson.Msg)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	r.callback(inputjson, func(err error, resultjson Resultobj) {
		json.NewEncoder(rw).Encode(resultjson.Result)
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

func RegisterService(name string, port int) (service, error) {
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

func (s *service) AddFunction(name string, callback Action) {
	fmt.Println("addFunction:", name)
	var request requestParserStruct
	request.callback = callback
	http.HandleFunc("/"+s.name+"/"+name, request.requestParser)
}

func (s *service) Start() {
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

/*
func main() {
	Service, _ := registerService("hi", 1234)
	fmt.Println(Service)
	Service.addFunction("hello", func(inputjson inputobj, done result) {
                fmt.Println(inputjson,inputjson.msg)
                var inputFuncObj functionsInput
                if err := json.Unmarshal([]byte(inputjson.msg),&inputFuncObj); err != nil {
                   panic(err)
                }
                fmt.Println("inputFuncObj",inputFuncObj)
		sum := inputFuncObj.Int1 + inputFuncObj.Int2
                result :=  map[string]int{"sum": sum}
                result2, err := json.Marshal(result)
                if err != nil {
                    panic(err)
                }
		done(nil, resultobj{result:string(result2) })
	})
	Service.start()
}*/

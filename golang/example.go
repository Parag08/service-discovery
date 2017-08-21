package main

import (
	"consul/serviceClient"
	"consul/serviceRegistry"
	"encoding/json"
	"fmt"
	"time"
)

type functionsInput struct {
	Int1 int `json:"int1"`
	Int2 int `json:"int2"`
}

func registerServiceMain(name string, port int) {
	Service, _ := serviceRegistry.RegisterService(name, port)
	fmt.Println(Service)
	Service.AddFunction("hello", func(inputjson serviceRegistry.Inputobj, done serviceRegistry.Result) {
		fmt.Println(inputjson, inputjson.Msg)
		var inputFuncObj functionsInput
		if err := json.Unmarshal([]byte(inputjson.Msg), &inputFuncObj); err != nil {
			panic(err)
		}
		fmt.Println("inputFuncObj", inputFuncObj)
		sum := inputFuncObj.Int1 + inputFuncObj.Int2
		result := map[string]int{"sum": sum}
		result2, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		done(nil, serviceRegistry.Resultobj{Result: string(result2)})
	})
	Service.Start()
}

func main() {
	go registerServiceMain("hi", 1234)
	time.Sleep(time.Second * 4)
	Client, _ := serviceClient.NewClient()
	jsonObjToSend := map[string]int{"int1": 5, "int2": 8}
	mapB, err := json.Marshal(jsonObjToSend)
	if err != nil {
		panic(err)
	}
	toSend := serviceClient.Inputobj{Msg: string(mapB)}
	Client.Send("hi", "hello", toSend)
	time.Sleep(time.Second * 10)
}

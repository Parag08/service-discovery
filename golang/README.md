# Service discovery with consul in golang

this is a very premitive implementtaion of consul based service discovery with go. The implementaion is incomplete

<b>goal</b>

all service can be registered as


1. Registering service   
`Service,err := registerService(<servicename> string, <port> int)`   
[after this step service should be register with consul and any node connected to that consul cluster can talk to the above mentioned service]  
    
2. Implementing or Adding a Function to Service  
`Service.addFunction(<functionname> string, functionToexcute func) ` 
  
3. Talking with other services  
`response,error :=Service.send(<servicename> string,<functionname> string,msgForTheFunction inputobj)`   
   
functions Implementations and structures   
   
inputobj   
```golang   
type inputobj struct {   
        msg string   
}   
```   
resultobj   
```golang   
type resultobj struct {  
        result string    
        err string  //this I havent implemented   
}  

    
type result func(error, resultobj)    
     
functionToexcute(inputjson inputobj, done result)   {
         /* do something on inputjson */   
         if err != nil {
            done(err,nil)
         }   else {
            done(nil,resultobj)
         }
}
```
if any one has any segestion regarding the implementation let me know   



<b>example</b>


```golang
//register a new service with a name hi and a function hello
Service, _ := serviceRegistry.RegisterService("hi", 1234)
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
    
//using the service with functionality

Client, _ := serviceClient.NewClient()
	jsonObjToSend := map[string]int{"int1": 5, "int2": 8}
	mapB, err := json.Marshal(jsonObjToSend)
	if err != nil {
		panic(err)
	}
	toSend := serviceClient.Inputobj{Msg: string(mapB)}
	Client.Send("hi", "hello", toSend)
```

# Service discovery with consul in golang

this is a very premitive implementtaion of consul based service discovery with go. The implementaion is incomplete

<b>goal</b>

all service can be registered as


1. Registering service   
Service,err := registerService(<servicename> string, <port> int)   
[after this step service should be register with consul and any node connected to that consul cluster can talk to the above mentioned service]  
    
2. Implementing or Adding a Function to Service  
Service.addFunction(<functionname> string, functionToexcute func)  
  
3. Talking with other services  
response,error :=Service.send(<servicename> string,<functionname> string,msgForTheFunction inputobj)   
   
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
        result interface{}    
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

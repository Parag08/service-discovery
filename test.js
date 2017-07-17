Service = require('./service-register.js')
service = new Service('testService')
service.addFunction('testaddfunction', function (req,res) {   //request or message,done
    console.log('request:recived:',req)
    res.send({'reply':'hi'})
})

service.send('testService','testaddfunction',{msg:{'aloha':'hi'}} ,function (err,resp) {
    console.log('msg:',resp,'err:',err)
    console.log('expected result:msg: { reply: hi } err: null')
})

service.send('testService','testfunction',{msg:{'aloha':'hi'}} ,function (err,resp) {
    console.log('msg:',resp,'err:',err)
    console.log('expected result:msg: null err: Service not responding')
})

service.addFunction('testaddfunction-2', function (req,res) {   //request or message,done
    console.log('request:recived:',req)
    res.send({'reply':'hi2'})
})


service.send('testService','testaddfunction-2',{msg:{'aloha':'hi'}} ,function (err,resp) {
    console.log('msg:',resp,'err:',err)
    console.log('expected result:msg: { reply: hi2 } err: null')
})

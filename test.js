Service = require('./service-register.js')
service = new Service('testService')
service.addFunction('testaddfunction', function (req,res) {   //request or message,done
    console.log('request:recived:',req)
    res.send({'reply':'hi'})
})

service.send('testService','testaddfunction',{msg:{'aloha':'hi'}} ,function (err,resp) {
    console.log('msg:',resp,'err:',err)
})

service.send('testService','testfunction',{msg:{'aloha':'hi'}} ,function (err,resp) {
    console.log('msg:',resp,'err:',err)
})

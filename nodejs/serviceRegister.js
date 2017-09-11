const consul = require('consul')();
const uuid = require('uuid');
const ip = require("ip");

function Service(serviceName,port) {
    this.serviceName = serviceName;
    this.port = port||29999
}

Service.prototype.getServiceName = function() {
    return this.serviceName
}

Service.prototype.start = function (callback) {
    app.listen(this.port, callback)
}

Service.prototype.registerService =  function() {
    serviceID = this.serviceName + '-'+ uuid.v4()
    let details = {
	name:this.serviceName,
        id:serviceID,
        address:ip.address(),
        port:this.port,
        check: {
           http: "http" + "://" + ip.address() + ":" + this.port + "/" + serviceID + "/status",
           interval: '10s',
           timeout: '5s',
           deregistercriticalserviceafter :'30s',
        }
    }
    console.log('registering:',details)
    consul.agent.service.register(details, function(err) {
        if (err) throw err;
    });
}



if (typeof require != 'undefined' && require.main==module) {
    service = new Service('testService')
    console.log(service.getServiceName())
    service.registerService()
}

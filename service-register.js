const express = require('express');
var bodyParser = require('body-parser')
const consul = require('consul')();
const uuid = require('uuid');
var request = require('request');
const app = express();
const os = require('os');
const timeout = 10
app.use(bodyParser.json())
const PID = process.pid;
const PORT = Math.floor(process.argv[2]) || 29999;
const HOST = os.hostname();
var PORT = 29999;

known_nodejs_services_instances = {}

var watcher = consul.watch({
    method: consul.health.service,
    options: {
        service: 'nodejs-services',
        passing: true
    }
});

watcher.on('change', data => {
    data.forEach(entry => {
        ID = entry.Service.ID
        path = `http://${entry.Service.Address}:${entry.Service.Port}/`
        service = ID.substr(0, ID.indexOf("%"));
        if (known_nodejs_services_instances[service] != undefined) {
            known_nodejs_services_instances[service].push(path)
        } else {
            console.log('new service added to cluster:', service.substr(7))
            known_nodejs_services_instances[service] = []
            known_nodejs_services_instances[service].push(path)
        }
    });
});

watcher.on('error', err => {
    console.error('watch error', err);
});

function createConsulID(serviceName, functionName) {
    return `nodejs-${serviceName}%${functionName}-${HOST}-${PORT}-${uuid.v4()}`
}

function Service(serviceName,port) {
    this.serviceName = serviceName;
    PORT = port||29999
    app.listen(PORT, function() {})
}

function findServiceFunction(servicename, functionName, req, callback) {
    var i = 0
    var ID = 'nodejs-' + servicename
    var interval = setInterval(function() {
        i = i + 1
        if (known_nodejs_services_instances[ID] != undefined) {
            path = known_nodejs_services_instances[ID][Math.floor(Math.random() * known_nodejs_services_instances[ID].length)]
            path = path + servicename + '/' + functionName + '/'
            console.log('path:', path, 'req:', req)
            request({
                    url: path,
                    json: true,
                    method: "POST",
                    headers: {
                        "content-type": "application/json",
                    },
                    body: req

                },
                function(error, response, body) {
                    if (!error && response.statusCode == 200) {
                        callback(null,body)
                    } else if (!error) {
                        callback(response.body,null)
                    }  else {
                        callback('function not responding',null)
                        console.error(error)
                        console.error('timeout error Service not responding','service:',servicename,'function:',functionName,'msg:',req)
                    }
                }
            );
            clearInterval(interval);
        }
        if (i >= timeout) {
            console.error('timeout error Service not found','service:',servicename,'function:',functionName,'msg:',req)
            callback('Service not found',null)
            clearInterval(interval);
        }
    }, 1000);
}

module.exports = Service

Service.prototype.getServiceName = function() {
    return this.serviceName
}

Service.prototype.send = function(servicename, functionName, req, callback) {
    findServiceFunction(servicename, functionName, req, callback);
}

Service.prototype.addFunction = function(functionName, callback) {
    CONSUL_ID = createConsulID(this.serviceName, functionName)
    console.log('addingfucntion:', functionName, 'service:', this.serviceName)
    app.post('/' + this.serviceName + '/' + functionName + '/', function(req, res) {
        callback(req.body, function (err,reply) {
            if (!err) {
                res.send(reply)
            } else {
                res.status(400).send(err);
            }
        })
    })
    let details = {
        name: 'nodejs-services',
        address: HOST,
        check: {
            ttl: '10s',
            deregister_critical_service_after: '1m'
        },
        port: PORT,
        id: CONSUL_ID
    };
    consul.agent.service.register(details, err => {
        if (err) {
            throw new Error(err);
        }
        console.log('registered with Consul');

        setInterval(() => {
            consul.agent.check.pass({
                id: `service:${CONSUL_ID}`
            }, err => {
                if (err) throw new Error(err);
            });
        }, 1000);
        process.on('SIGINT', () => {
            console.log('SIGINT. De-Registering...');
            let details = {
                id: CONSUL_ID
            };
            consul.agent.service.deregister(details, (err) => {
                console.log('de-registered.', err);
                process.exit();
            });
        });
    });
}

import consul
import uuid
from consul.base import Check

class client:
    def __init__(self,host='127.0.0.1', port=8500, token=None, scheme='http', consistency='default', dc=None, verify=True, cert=None):
        self.host = host
        self.port = port
        self.token = token
        self.scheme = scheme 
        self.consistency = consistency
        self.dc = dc
        self.verify = verify
        self.cert = cert 
        self.conuslClient = consul.Consul()
    def register(self,serviceName,servicePath,check ='http'):
        self.UUID = serviceName+str(uuid.uuid4())
        self.check =  Check.http(servicePath,'5s',timeout='10s', deregister='20s', header=None)
        self.conuslClient.agent.service.register(serviceName,service_id=self.UUID,address=servicePath)
        self.conuslClient.agent.check.register(serviceName,check=self.check,service_id=self.UUID)
    def checkStatus(self): 
        print('hi')
        print(self.conuslClient.agent.services(),self.conuslClient.agent.checks())

if __name__ == '__main__':
     clientObj = client()
     clientObj.register('test','http://test:8767')
     clientObj.checkStatus()
     

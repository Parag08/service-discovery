import socket   
import unittest
import transport.http as http
import discovery.consul as consul

class ServiceTestCase(unittest.TestCase):
    def setUp(self):
        self.service = Service('testService',port=7889,host='testhost')
    def test_service_init(self):
        self.assertEqual(self.service.serviceName,'testService','service name not set properly')
        self.assertEqual(self.service.port,7889,'service port not set properly')
        self.assertEqual(self.service.host,'testhost','service port not set properly')
    def test_service_addFunction(self):
        def add(a,b):
            return a + b
        self.service.addFucntion('add',add)
        self.assertEqual(self.service.path,['http://testhost:7889/add'],'registering function not completed')
        self.assertEqual(self.service.funcDictionary.keys(),['add'],'registering function not completed')
    def tearDown(self):
        print('teardown')


class Service:
    def __init__ (self,serviceName,port=1234,host=None):
        self.serviceName = serviceName
        self.port = port
        self.host = host
        if host == None:
            self.host = [l for l in ([ip for ip in socket.gethostbyname_ex(socket.gethostname())[2] if not ip.startswith("127.")][:1], [[(s.connect(('8.8.8.8', 53)), s.getsockname()[0], s.close()) for s in [socket.socket(socket.AF_INET, socket.SOCK_DGRAM)]][0][1]]) if l][0][0]
        self.path = []
        self.funcDictionary = {}
        self.client = consul.client()
        self.client.register(self.serviceName,'http://'+self.host+':'+str(self.port))
    def addFucntion(self,functionName,function):
        self.path.append('http://'+self.host+':'+str(self.port)+'/'+functionName)
        self.funcDictionary[functionName] = function
#class client:




if __name__ == '__main__':
    unittest.main()

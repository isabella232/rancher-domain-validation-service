#Domian Validation

rancher-domain-validaiton-service
========

A microservice that does micro things.

## Building

`make`


## Running

`./bin/rancher-domain-validaiton-service`





## License
Copyright (c) 2014-2016 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


##Test Environment Setup
- RANCHER SERVER : http://54.255.182.226:8080/
- DNS SERVER: 54.255.182.226
- HOST SERVER: 54.169.69.238

###DNS server setup

- DNS Server side: Add a DNS Record in bind6 server
	1. `sudo apt-get install bind9`
	2. Add the zone in 
		- `vim /etc/bind/named.conf.local `
		- e.g.`zone "zp.com"  { type master; file "/etc/bind/db.zp.com"; };`
	3. set up the zone database
		- `cp db.local db.zp.com`
		- `vim db.zp.com`
		
```
;
; BIND data file for local loopback interface
;
$TTL	604800
@	IN	SOA	localhost. root.localhost. (
			      2		; Serial
			 604800		; Refresh
			  86400		; Retry
			2419200		; Expire
			 604800 )	; Negative Cache TTL
;
@	IN	NS	localhost.
@	IN	A	127.0.0.1
@	IN	AAAA	::1
test       IN      A        54.169.69.238
www IN      A               54.169.69.238
@            IN      A      54.169.69.238
_hna-challenge    IN   TXT  "4501163876f8a13521f9233cb3d0464f36e61bbc8965d"

```


	
- Client side:

	1. Set the DNS to test DNS serve
	2. Add the DNS server address in the /etc/reslv.conf
	3. Test `ping` host

	
```

	sudo vi /etc/resolvconf/resolv.conf.d/head
	search nyc3.example.com  # your private domain
nameserver 10.128.10.11  # ns1 private IP address
nameserver 10.128.20.12  # ns2 private IP address

sudo resolvconf -u
```

###Web server setup

- Client side: 
  1. Start up a web service(Apache)
  2. create a hna.txt under `${webroot}/.well-known/hna.txt`



## API Document
The "X-API-Account-Id" will extract from token validation server

###Sample Error Message 

`{"type":"error","code":"400","status":"BadRequest","message":"Error delete the record: \u003cnil\u003e"}`    

    
The "X-API-Account-Id" will praseing from token validation server
    
###Create: POST /v1-domains/domains 
    - input`Http.Header("Cookie","token=xxxxx")` `{"domainName": "food.com","projectid": "1d1","X-API-Account-Id": "1a1"}`
    - output `{"type":"sucess","code":"200","status":"Succeed","message":"inserting the record succeed"}`
    

###LIST: GET /v1-domains/domains/

- List all the domain name in the current env: GET /v1-domains/domains/
	
		 - input 
		 	`Http.Header("X-API-ProjectID-Id","1a9")`
		 	`Http.Header("Cookie","token=xxxxx")`
	    - output 
	```
	[{"accountId":"1a8","projectId":"1a9","domainName":"food.com","state":"active","hashvalue":"29644da38d8b870a5cf519dc77746d4bd99a470cee44b","containerID":"1d2"},{"accountId":"1a8","projectId":"1a9","domainName":"food.com","state":"Pending","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d3"},{"accountId":"1a8","projectId":"1a9","domainName":"alpha.com","state":"active","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d4"},{"accountId":"1a8","projectId":"1a9","domainName":"beta.com","state":"active","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d5"}]
	```
		
		
- List all the domain name in the system when token belongs to `admin`: GET /v1-domains/domains/
	
		 - input 
		 	`Http.Header("X-API-ProjectID-Id","1a9")``Http.Header("Cookie","token=xxxxx")`
	    - output 
			```
			[{"accountId":"1a8","projectId":"1a9","domainName":"food.com","state":"active","hashvalue":"29644da38d8b870a5cf519dc77746d4bd99a470cee44b","containerID":"1d2"},{"accountId":"1a8","projectId":"1a9","domainName":"food.com","state":"Pending","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d3"},{"accountId":"1a8","projectId":"1a9","domainName":"alpha.com","state":"active","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d4"},{"accountId":"1a8","projectId":"1a9","domainName":"beta.com","state":"active","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d5"}]
			```
		
- List all the information for the container: GET /v1-domains/domains/{containerid}
	
	- input 
		 	`Http.Header("X-API-ProjectID-Id","1a9")``Http.Header("Cookie","token=xxxxx")`
	- output 
	 
		 
			`		[{"accountId":"1a8","projectId":"1a9","domainName":"food.com","state":"Pending","hashvalue":"4501163876f8a13521f9233cb3d0464f36e61bbc8965d","containerID":"1d3"}]
			`

###Delete: DELETE /v1-domains/domains/containerid
- delete the domain name under current envid. Admin token can delete any container.

	- input`Http.Header("Cookie","token=xxxxx")``Http.Header("X-API-ProjectID-Id","ia3")`
	- output `{"type":"sucess","code":"200","status":"Succeed","message":"Delete the record succeed"}`

###Validate: POST /v1-domains/domain/{id}?action=validate

User can validate the domain by HTTP webroot or DNS TXT record

- input
	`Http.Header("PL=rancher; token=YQmiC6QzqcWP8jQ5W9e69Gkymm2UhTqGpgRkqqkJ; CSRF=12CCE5D3BC")`
- output `{"type":"sucess","code":"200","status":"Succeed","message":"Validate the domain by webroot . Update the record succeed"}`
`{"type":"error","code":"404","status":"NotFound","message":"Cannot validate domain name"}`

###Domain Filter:POST /v1-domains/filter/projects/1a9/loadbalancerservice
Domian filter will filter out invalid domain name when user create Load balancer

- input

`Http.Header("Cookie","PL=rancher; token=YQmiC6QzqcWP8jQ5W9e69Gkymm2UhTqGpgRkqqkJ")`
`Http.Header("x-api-csrf","12CCE5D3BC")`

```
{  
   "assignServiceIpAddress":false,
   "scale":1,
   "startOnCreate":true,
   "type":"loadBalancerService",
   "name":"lbbob",
   "description":null,
   "stackId":"1st22",
   "launchConfig":{  
      "instanceTriggeredStop":"stop",
      "kind":"container",
      "networkMode":"managed",
      "privileged":false,
      "publishAllPorts":false,
      "readOnly":false,
      "startOnCreate":true,
      "stdinOpen":false,
      "tty":false,
      "vcpu":1,
      "imageUuid":"docker:rancher/lb-service-haproxy:v0.4.6",
      "type":"launchConfig",
      "restartPolicy":{  
         "name":"always"
      },
      "ports":[  
         "23:23/tcp"
      ],
      "expose":[  

      ],
      "labels":{  

      },
      "blkioWeight":null,
      "cgroupParent":null,
      "count":null,
      "cpuCount":null,
      "cpuPercent":null,
      "cpuPeriod":null,
      "cpuQuota":null,
      "cpuSet":null,
      "cpuSetMems":null,
      "cpuShares":null,
      "createIndex":null,
      "created":null,
      "deploymentUnitUuid":null,
      "description":null,
      "diskQuota":null,
      "domainName":null,
      "externalId":null,
      "firstRunning":null,
      "healthInterval":null,
      "healthRetries":null,
      "healthState":null,
      "healthTimeout":null,
      "hostname":null,
      "ioMaximumBandwidth":null,
      "ioMaximumIOps":null,
      "ip":null,
      "ip6":null,
      "ipcMode":null,
      "isolation":null,
      "kernelMemory":null,
      "memory":null,
      "memoryMb":null,
      "memoryReservation":null,
      "memorySwap":null,
      "memorySwappiness":null,
      "milliCpuReservation":null,
      "oomScoreAdj":null,
      "pidMode":null,
      "pidsLimit":null,
      "removed":null,
      "requestedIpAddress":null,
      "shmSize":null,
      "startCount":null,
      "stopSignal":null,
      "user":null,
      "userdata":null,
      "usernsMode":null,
      "uts":null,
      "uuid":null,
      "volumeDriver":null,
      "workingDir":null,
      "networkLaunchConfig":null
   },
   "lbConfig":{  
      "type":"lbConfig",
      "config":null,
      "certificateIds":[  

      ],
      "stickinessPolicy":null,
      "portRules":[  
         {  
            "protocol":"http",
            "type":"portRule",
            "priority":1,
            "hostname":"food.com",
            "sourcePort":23,
            "serviceId":"1s27",
            "targetPort":23
         }
      ]
   },
   "created":null,
   "externalId":null,
   "fqdn":null,
   "healthState":null,
   "kind":null,
   "removed":null,
   "selectorLink":null,
   "uuid":null,
   "vip":null
}
```
- output

Error: `{"type":"error","code":"403","status":"Forbidden","message":"Domain abs.com is not valid"}`
Success:

```
{  
   "type":"collection",
   "resourceType":"loadBalancerService",
   "links":{  
      "self":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice"
   },
   "createTypes":{  

   },
   "actions":{  

   },
   "data":[  

   ],
   "sortLinks":{  
      "accountId":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=accountId",
      "createIndex":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=createIndex",
      "created":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=created",
      "description":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=description",
      "externalId":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=externalId",
      "healthState":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=healthState",
      "id":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=id",
      "kind":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=kind",
      "name":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=name",
      "removeTime":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=removeTime",
      "removed":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=removed",
      "selectorContainer":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=selectorContainer",
      "selectorLink":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=selectorLink",
      "stackId":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=stackId",
      "state":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=state",
      "system":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=system",
      "uuid":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=uuid",
      "vip":"http:\/\/54.255.182.226:8080\/v2-beta\/projects\/1a9\/loadbalancerservice?sort=vip"
   },
   "pagination":{  
      "first":null,
      "previous":null,
      "next":null,
      "limit":100,
      "total":null,
      "partial":false
   },
   "sort":null,
   "filters":{  
      "accountId":null,
      "createIndex":null,
      "created":null,
      "description":null,
      "externalId":null,
      "healthState":null,
      "id":null,
      "kind":null,
      "name":null,
      "removeTime":null,
      "removed":null,
      "selectorContainer":null,
      "selectorLink":null,
      "stackId":null,
      "state":null,
      "system":null,
      "uuid":null,
      "vip":null
   },
   "createDefaults":{  

   }
}
```





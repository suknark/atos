# Proxy to view and manage data in memcached, aerospike and elasticsearch

Docs for commands aerospike  
http://www.aerospike.com/docs/reference/info/  
Docs for commands memcached  
https://code.google.com/p/memcached/  
Docs for commands elasticsearch  
https://www.elastic.co/guide/en/elasticsearch/guide/master/_talking_to_elasticsearch.html

**Simple usage:**
```
root@duma:~$ ./atos 
> h
Simple usage:
 type "aerospike" to use aerospike storage
 type "memcached" to use memcached-storage
> memcached
goched> stats items
STAT items:1:number 10
STAT items:1:age 79546
```
or  
```
root@duma:~$ ./atos memcached
goched> stats items
STAT items:1:number 10
STAT items:1:age 79546  
```
or
```
root@duma:~$ ./atos memcached stats items
STAT items:1:number 10
STAT items:1:age 79546
root@duma:~$
```
You can move from one service to other:  
```
root@duma:~$ ./atos memcached
goched> aerospike
gospike> q   
> exit
```
Elasticsearch usage:  
```
root@duma:~$ ./atos                                                                                                                                                                              
> elasticsearch
gostic> get /
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Content-Length: 343

{
  "status" : 200,
  "name" : "George Tarleton",
  "cluster_name" : "elasticsearch",
  "version" : {
    "number" : "1.5.2",
    "build_hash" : "62ff9868b4c8a0c45860bebb259e21980778ab1c",
    "build_timestamp" : "2015-04-27T09:21:06Z",
    "build_snapshot" : false,
    "lucene_version" : "4.10.4"
  },
  "tagline" : "You Know, for Search"
}

gostic> 
```
Last commit add readline from github.com/fiorix/go-readline - tab and history

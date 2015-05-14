# Proxy to view and manage data in memcached and aerospike

Docs for commands aerospike  
http://www.aerospike.com/docs/reference/info/  
Docs for commands memcached  
https://code.google.com/p/memcached/

*Simple usage:*
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
You can move from one service to other:  
```
root@duma:~$ ./atos memcached
goched> aerospike
gospike> q   
> exit
```

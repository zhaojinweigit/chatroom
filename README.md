# chatroom
golang websocket chatroom

build and run:
	./build.sh 
	./server 
	
visit website:
	visit http://127.0.0.1:10328/?name=zjw

you can also visit http://127.0.0.1:10328/  and input your name in the page
if you wanto to change ip,change ip in js/index.js and server.go
if you install redis, use var isHistory bool = true  to turn on the history feature
and you can put your headphoto to ./image/ and  write it in func init()


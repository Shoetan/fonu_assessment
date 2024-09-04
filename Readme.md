## WEB SOCKET SERVER
One of the drawbacks of HTTP is that it does not support real-time communication. Instead of having to poll every time web sockets eliminate that by having a Persistent connection with the server to keep real-time communication up and running 

## ARCHITECTURE
Here I have two distant servers called server A and server B respectively. Starting this are just regular HTTP servers but are upgraded to a web socket server
```
var upgrader =  websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,}
```
Both servers run on defined port servers A: 8080 and B: 8081 using TCP protocol.

## GETTING STARTED

* Clone repository to your machine
``git clone git@github.com:Shoetan/fonu_assessment.git``
* The project is divided into servers and utils. Servers hold the implementation for the two distinct servers and utils hold the implementation of most of the reusable functions.
* To start both servers cd into the servers folder by ``cd servers``
* On two different terminal windows cd unto server a in one and server b in another. ``cd serverA`` and ``cd serverB``
* You have the setup like  this
![terminal setup](./Resources/start.png)
* To start both servers run ``go run main.go start`` on both terminals. I am using Cobra to build this as a CLI application. NB: On starting the servers it waits for 10 seconds before trying to establish a connection to the other server. This just helps to make sure both servers are running else it terminates.
* The servers will come up and you will have an interface that looks like below
![servers up](./Resources/connetion.png)
* A client has been connected to server A and server B  
![client](./Resources/client.png)
* To send a message from a connected client to a server use the command line and enter ``1 <your message>``. In the background your message is stripped from your command




## OBSERVATIONS AND CONTRAINTS
* The concept of connecting servers ( in this context both servers acting as clients to each other is beyond 

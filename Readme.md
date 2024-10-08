## WEB SOCKET SERVER
One of the drawbacks of HTTP is that it does not support real-time communication. Instead of having to poll every time web sockets eliminate that by having a Persistent connection with the server to keep real-time communication up and running 

## ARCHITECTURE
Here I have two distant servers called server A and server B respectively. Starting this are just regular HTTP servers but are upgraded to a web socket server

Clients can connect to the TCP servers via the command this way the TCP server can be interacted with.
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

* To start both servers run ``go run main.go start`` on both terminals. I am using Cobra to build this as a CLI application. NB: On starting the servers it waits for 10 seconds before trying to establish a connection to the other server. This helps to make sure both servers are running, or else it terminates.
  
* The servers will come up and you will have an interface that looks like below
![servers up](./Resources/connetion.png)

* To connect a client to either server A or server B  navigate to the clients folder ``cd clients`` and run ``go run main.go`` to connect to an already specified server
![client](./Resources/client.png)

* To send a broadcast message from a connected client to a server use the command line and enter ``1 <your message>``. In the background, your message is stripped from your command
![strip](./Resources/strip.png)

* In action this looks like writing the below

![send message](./Resources/broad.png)

* The server will receive your message and display it

![receive message](./Resources/LOOP2.png)

* The server also checks if the message is valid or not
  
![receive message](./Resources/invalid.png) 

* The servers and runtime are running concurrently in different go routines hence it does not time out until you send a close signal with ``ctrl + c``

## TESTING

Most important codes are found in the utils folder and tests can be run
``go test -run TestConnectToWebSocketServer``
``go test -run TestHandleConnections``

## OBSERVATIONS
* The concept of connecting servers ( in this context both servers act as clients to each other) while this is possible this approach is not recommended as per best practices and will lead to very complex or even spaghetti code spent a lot time trying to get it to work. WebSocket servers usually wait for clients to connect to them. They don't usually start the connection themselves. Trying to connect them together will only confuse their roles.

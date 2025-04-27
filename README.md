# ED-tcp-go-chat
Smol tcp server/client chat Im making for educational purposes. 

---

# Starting client and a server

Simply run ```go run server.go``` in one terminal and ```go run client.go``` in another terminal.

Congratulations! You are connected and can start sending messages. Start a third terminal and client, so you can now chat with yourself and meet the demands of your schizophrenia!

---

# Available commands in client

-   ```/help``` for a list of commands in client
-   ```/leave``` leave your current room. Redirects you to General, which is created by default
-   ```/create <name>``` create a room with a given name
-   ```/join <name>``` join a room with a given name
-   ```/list``` get a list of all rooms

---

# Server start options

- ```-log``` Defult is ```false```, meaning server logs will go directly into terminal. Setting it ```true``` makes server output log into ```log.log``` file.
- ```-ip``` Default is ```127.0.0.1```, meaning server will start on ```localhost```.
- ```-port``` Default is ```8080```, meaning it will start on this port.

# ED-tcp-go-chat
Smol tcp server/client chat Im making for educational purposes. 

---

# Starting client and a server

Simply run ```go run server.go``` in one terminal and ```go run client.go``` or ```go run client_no_bubble.go``` in another terminal.

Congratulations! You are connected and can start sending messages. Start a third terminal and client, so you can now chat with yourself and meet the demands of your schizophrenia!

---

# Available commands in client

>[!WARNING]
>NO BUBBLE CLIENT DOESNT WORK PROPEPRLY AS OF 29.04.2025, SINCE READING AND WRITING HAS BEEN REWORKED TO JSON FORMAT.

-   ```/help``` for a list of commands in client **NEEDS REWORK**
-   ```/leave``` leave your current room. Redirects you to General, which is created by default **NO INTERFACE IN TUI VERSION**
-   ```/create <name>``` create a room with a given name
-   ```/join <name>``` join a room with a given name **NO INTERFACE IN TUI VERSION**
-   ```/list``` get a list of all rooms **NO INTERFACE IN TUI VERSION**

---

# Server start options

- ```-log``` Defult is ```false```, meaning server logs will go directly into terminal. Setting it ```true``` makes server output log into ```log.log``` file.
- ```-ip``` Default is ```127.0.0.1```, meaning server will start on ```localhost```.
- ```-port``` Default is ```8080```, meaning it will start on this port.

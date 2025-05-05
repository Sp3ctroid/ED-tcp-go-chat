# ED-tcp-go-chat
Smol tcp server/client chat Im making for educational purposes. 

---

# Starting client and a server

Simply run ```go run server.go``` in one terminal and ```go run client.go``` or ```go run client_no_bubble.go``` in another terminal.

Congratulations! You are connected and can start sending messages. Start a third terminal and client, so you can now chat with yourself and meet the demands of your schizophrenia!

---

# Available commands in NO TUI client

>[!WARNING]
>NO BUBBLE CLIENT DOESNT WORK PROPEPRLY AS OF 29.04.2025, SINCE READING AND WRITING HAS BEEN REWORKED TO JSON FORMAT.

-   ```/help``` for a list of commands in client
-   ```/leave``` leave your current room. Redirects you to General, which is created by default
-   ```/create <name>``` create a room with a given name
-   ```/join <name>``` join a room with a given name
-   ```/list``` get a list of all rooms

---

# Available commands in TUI client

-   ```ctrl+n``` - create new room
-   ```ctrl+j``` - join room by its name
-   ```ctrl+a``` - open a list of rooms. By pressing ```ENTER``` you can join selected room
- ```ctrl+u``` - change your username

You can go back to your chat room by pressing ```esc```

---

# Server start options

- ```-log``` Defult is ```false```, meaning server logs will go directly into terminal. Setting it ```true``` makes server output log into ```log.log``` file.
- ```-ip``` Default is ```127.0.0.1```, meaning server will start on ```localhost```.
- ```-port``` Default is ```8080```, meaning it will start on this port.

---

# Client start options

- ```-ip``` Default is ```127.0.0.1```, meaning client will try to connect to a server, that is launched on ```localhost```
- ```-port``` Default is ```8080```.

---

# Server storage

To store its users and rooms, server uses corresponding interfaces ```UserStore``` and ```RoomStore```. This ensures, that developer can easily swap between storages by implementing these interfaces. Introduced this so that I can start breaking server down into files, so that code is more readable and modifiable.

Curently there's only map storage, which can also be considered a mock storage for testing.

To make a DB storage, you simply need to implement interfaces mentioned above.

---

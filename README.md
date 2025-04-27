# ED-tcp-go-chat
Smol tcp server/client chat Im making for educational purposes. 

---

# Starting client and a server

Simply run ```go run server.go``` in one terminal and ```go run client.go``` in another terminal.

Congratulations! You are connected and can start sending messages. Start a third terminal and client, so you can now chat with yourself and meet the demands of your schizophrenia!

---

# Available commands

-   ```/help``` for a list of commands in client
-   ```/leave``` leave your current room. Redirects you to General, which is created by default
-   ```/create <name>``` create a room with a given name
-   ```/join <name>``` join a room with a given name
-   ```/list``` get a list of all rooms

---

# Logging

Server automatically logs every action it performs in ```log.log``` file, that is created once the server starts. If a file exists it will continue writing in it.

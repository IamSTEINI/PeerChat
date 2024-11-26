PeerChat
================

A simple GO p2p chat application for peer-to-peer communication. This project is open-source and free to use for anyone, but please give credits to me ;)

I may want to do this but encrypted (maybe with PQC techniques). Feel free to fork it.


As a strong advocate for privacy and security, I want to find solutions that prioritize user privacy and security while maintaining decentralization and robust encryption.

**Privacy and Security Concerns**
--------------------------------

Mandatory storage and accessibility of chat logs can compromise user privacy and create vulnerabilities for potential data breaches. These measures may not effectively prevent illegal activities, as criminals can find ways to circumvent them. Instead, they infringe on individuals' right to privacy and may do more harm than good.

**Decentralized and Encrypted Solutions**
-----------------------------------------

Decentralization is crucial for ensuring free speech and protecting against data breaches.

<br>
<br>
**Centralized, unencrypted storage is more vulnerable to hacking than decentralized, encrypted access.**

___

**Disclaimer:** This project is not 100% perfect and might have some bugs or areas for improvement. It's a simple implementation of a p2p chat and is intended for educational or personal use. If you find any issues or have suggestions for improvement, feel free to contribute to the project.

**How does it work?**
---------------------

1. **Node Selection**: The user selects an option to either host a node or join an existing node.
2. **Node Hosting**: If the user chooses to host a node, the application listens for incoming connections on a random available port. The node's IP and port are displayed for sharing with other peers.
3. **Node Joining**: If the user chooses to join a node, they are prompted to enter the IP and port of the node they want to join. The application attempts to establish a connection with the specified node.
4. **Connection Establishment**: When a connection is established, the application sends a request to the connected node to share its peers. This allows the application to discover other nodes in the network.
5. **Peer Sharing**: Periodically, each node shares its list of connected peers with other nodes in the network. This enables nodes to discover and connect to new peers.
6. **Real-time Chat**: Users can send messages to other peers in the network. Messages are broadcast to all connected peers in real-time.
7. **Connection Management**: The application manages connections, updating the list of active peers and handling disconnections.

**Features:**
------------

-   Peer-to-peer communication without a central server
-   Ability to host a node and join other nodes
-   Real-time chat functionality
-   Periodic sharing of peers for network discovery

**How to use:**
--------------

1. Clone the repository and run the application with `go run peer.go`
2. Select an option to either host a node or join an existing node.
3. If joining, enter the IP:PORT of the node you want to join.
4. Start chatting with other peers in the network.

**Contributing:**
--------------

If you'd like to contribute to this project, please fork the repository, make your changes, and submit a pull request. Your contributions are welcome and will be appreciated.

**License:**
---------

This project is licensed under the MIT License.

### THANKS

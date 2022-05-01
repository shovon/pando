# tree

This is a document that describes the distributed udag management server.

Summary: nodes will not be auto-arranging themselves into a tree topology; instead it's a server that will hold a representation of how the clients will be related to each other.

Properties of the server's internal representation of a tree will be the following:

- each client will be given IDs of 3 neighbouring nodes
- the graph will be acyclic
- the graph will be undirected

By the definition of the last two-points, the graph will therefore be a tree.

If two non-neighbouring nodes need to communicate with eachother, either the application utilizing the spanning tree will need to design the application to support the act of relaying messages from node-to-node, or the application will need to utilize a third-party server, unrelated to the spanningtree server

## Protocol

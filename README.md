# seamstress
Fabric tools and packages

The purpose of this repo is to build a well documented set of simple tools
for connecting to Fabric instances on IBM Blockchain Platform and locally.

The initial structure of this repo is a single main package CLI with minimal
external dependencies. This main package may have multiple subcomamnds to
achieve a variety of useful functions. This may change as more specialized
tooling becomes needed.

# Usage (work in progress - README driven development)

Set up IBM Blockchain Platform (with CA, Organizations, Peers, Orderer, Channels, Smart contracts, etc).

* [Build Network](https://cloud.ibm.com/docs/blockchain?topic=blockchain-ibp-console-build-network)
* [Deploy smart contract](https://cloud.ibm.com/docs/blockchain?topic=blockchain-ibp-console-smart-contracts)

Export connection profile and MSP from IBM Blockchain Platform. Console > Organizations > YourOrg > Create connection profile > Select all peers > Download connection profile.

Connect to channel using seamstress:

```
$ seamstress connect path-to-connection-profile
connection established (proof: ...)
```

Invoke smart contract using seamstress:
```
seamstress invoke action-name -d '{"data":"values"}'
```

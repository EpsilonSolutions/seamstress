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

Configure seamstress, example:
```
$ seamstress configure
enter filepath to profile:
./local/Org1msp_profile.json
enter channel ID:
channel1
enter organization name:
OrgMSP
enter smart contract ID / chaincode ID:
simple-contract
```

Invoke smart contract using seamstress, example creating a key value pair:
```
seamstress invoke create '{"key":"V31OSDW57K","value":"IRRGjYWnQ9"}'
```

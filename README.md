# seamstress
Fabric tools and packages

The purpose of this repo is to build a well documented set of simple tools
for connecting to Fabric instances on IBM Blockchain Platform and locally.

This repo contains:

* A command line interface `cmd/seamstress-cli` which can be used to
  conveniently invoke chaincode from your command line.

* A worker microservice (in the top level directory) which reads messages from
  a NATS message bus and passes them to a preconfigured smart contract on a
  preconfigured channel.

* A library package `fabric` which can be used to create a range of simple
  clients tools for invoking a fabric smart contract.

This repo does *NOT* contain:

* Something that sets up a fabric blockchain (TODO: add link)

* Something that sets up a channel on a fabric blockchain (TODO: add link)

* Something that installs or updates a smart contract on a fabric blockchain (TODO: add link)

The omissions here are deliberate, because the goal of this repo is to simplify
application development. Administration of the blockchain can be left to other
tools.

# Usage for cmd/seamstress-cli

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

# A word on the worker

The worker service contained in this repo does not transform or validate the
message contents. It serves as an adaptor to allow any platform that can write
to NATS to invoke fabric smart contracts.

This adaptor pattern means that application developers can write their APIs in
any language or platform that they want to without having to do any direct
integration with platform specific fabric SDKs.

The use of NATS also allows an application stack to asynchronously perform high
latency invocations of smart contracts (~10s) without blocking resources on
their API services (which can therefore be written more rapidly with less of a
focus on performance / resource utilization).

A secondary goal therefore is for this worker to be highly performant and side
scalable. Of course there is also an absolute limit on the fabric blockchain
side so this service architecture should be able to deal with backpressure
too. Putting all of these challenges together and abstracting them into a
simple function independent worker service frees up other parts of the stack to
focus on business logic and feature requirements.


# Vassago
Efficient and Authenticated Provenance Query on Multiple Blockchains

This project is the demo of the paper: ["Vassago: Efficient and Authenticated Provenance Query on Multiple Blockchains" SRDS'21](https://ieeexplore.ieee.org/document/9603540) 

This demo is built atop Hyperledger Fabric. We use Channel in Hyperledger Fabric to simulate multiple blockchains and Chaincode as Smart Contract to define provenance dependency, write transactions, and perform provenance queries.

Simply out, to run this demo, you need to start a local Hyperledger Fabric project, then deploy Vassago on the project. If you succeed, you can test Vassago's provenance query latency. 

## Dependency
To deploy Vassago, you need to deploy Hyperledger Fabric v2.2. Create Channels with the following commands. Move the chaincodes to the correct position to be successfully installed.

How to deploy Hyperledger Fabric, please visit [here](https://hyperledger-fabric.readthedocs.io/en/release-2.2/whatsnew.html).

## How to use
### Start Fabric network

```shell
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

./network.sh up
```

### Set Environment variables for Org1
```shell
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

### Establish Channel
```shell
./network.sh createChannel -c dbchain
./network.sh createChannel -c tb1chain
./network.sh createChannel -c tb2chain
./network.sh createChannel -c tb3chain
./network.sh createChannel -c tb4chain
./network.sh createChannel -c tb5chain
./network.sh createChannel -c tb6chain
```

### Install chaincode
```shell
./network.sh deployCC -c dbchain -ccn dbchaincode -ccp ../vassago_chaincode/DBChaincode2 -ccl go
./network.sh deployCC -c tb1chain -ccn tb1chaincode -ccp ../vassago_chaincode/TB1Chaincode -ccl go
./network.sh deployCC -c tb2chain -ccn tb2chaincode -ccp ../vassago_chaincode/TB2Chaincode -ccl go
./network.sh deployCC -c tb3chain -ccn tb3chaincode -ccp ../vassago_chaincode/TB3Chaincode -ccl go
./network.sh deployCC -c tb4chain -ccn tb4chaincode -ccp ../vassago_chaincode/TB4Chaincode -ccl go
./network.sh deployCC -c tb5chain -ccn tb5chaincode -ccp ../vassago_chaincode/TB5Chaincode -ccl go
./network.sh deployCC -c tb6chain -ccn tb6chaincode -ccp ../vassago_chaincode/TB6Chaincode -ccl go
```

### Initialize chaincode
```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb1chain -n tb1chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb2chain -n tb2chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb3chain -n tb3chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb4chain -n tb4chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb5chain -n tb5chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C tb6chain -n tb6chaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C dbchain -n dbchaincode --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:10051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'
```



### Query

### Intra-blockchain query
On Blockchain TB6,  perform serial query the provenance of product with serial number B1613

```shell
peer chaincode query -C tb6chain -n tb6chaincode -c '{"function":"IntraChainTransfer","Args":["B1613"]}'
```

### Inter-blockchain query
Perform cross-chain query for Banana with serial number B1613, find the dependent chain, query in parallel according to the dependency, and test the latency
```shell
peer chaincode query -C tb6chain -n tb6chaincode -c '{"function":"CrossChainTransfer","Args":["Banana"]}'
./scratch.sh B1613 tb4chain tb5chain tb6chain
```

#### Parallel query latency test
For the product with serial number M1614, perform provenance query on blockchain TB1~TB6, and test the parallel latency.
```shell
./scratch.sh M1614 tb1chain tb2chain tb3chain tb4chain tb5chain tb6chain
```

#### Serial query latency test
For the product with serial number M1614, perform provenance query on blockchain TB1~TB6, and test the serial latency.
```shell
./serial.sh M1614 tb1chain tb2chain tb3chain tb4chain tb5chain tb6chain
```


### Down the system

```shell
./network.sh down
```

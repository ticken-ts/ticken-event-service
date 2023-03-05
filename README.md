# Ticken - Event Service

## Architectural design

TODO 
## Running for development
To run in dev mode you need to set the "ENV" environment variable to "dev"

Running in dev mode will mock most external services, you only need to run the following services:
- MongoDB
- Ganache (for the blockchain)

You still need to create the .env and the config.json files using the examples.

## Running locally

This is project is built in way that it can be run locally. 
To achieve this run locally three services:

* Mongo DB instance
* This server (ticken-ticket-service)
* Hyperledger Fabric peer with two chaincodes:
  * ticken-ticket-chaincode
  * ticken-event-chaincode

Before starting clone the following repos in the same folder:
* [ticken-dev](https://github.com/tpp-facu-javi/ticken-dev): contains
all docker images that we are going to use and the scripts to run them.

* [ticken-chaincodes](https://github.com/tpp-facu-javi/ticken-chaincodes): contains 
ticken-event chaincode and ticken-ticket chaincode

All scripts are going to be inside the folder `dev-services` inside `ticken-dev`

### Running the MongoDB instance

```
sh ./start-mongo.sh
```

This is going to start a docker container with a mongo db image.
The image name is `ticken-mongo`

### Running the Hyperledger Fabric Peer

```
sh ./start-pvtbc.sh
```

This is going to start all the images needed to run an hyperledger fabric peer and it
will deploy all necessary chaincodes.

### Running ticken-event-service

Once you run successfully the hyperledger fabric peer and the MongoDB instance, 
you can start this service

**Running without Docker**

- Start private blockchain with start-pvtbc.sh from ticken-dev
```bash
$ bash ../ticken-dev/dev-services/start-pvtbc.sh
```
- Create the `config.json` file based on `config.example.json` and paste the full path to the certificates and private key from an Organization 1 user in the `pvtbc` section:
```json
{
    "msp_id": "Org1MSP",
    "peer_endpoint": "localhost:7051",
    "gateway_peer": "peer0.org1.example.com",
    "certificate_path": "<<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
    "private_key_path": "<<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk",
    "tls_certificate_path": "<<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/ca.crt"
}
```
- Launch a mongodb instance using the script from `ticken-dev`
```bash
$ bash ../ticken-dev/dev-services/start-mongo.sh
```
- Create `.env` file and paste the mongodb url
```bash
ENV="dev"
CONFIG_FILE_PATH="."
CONFIG_FILE_NAME="config"
DB_CONN_STRING="mongodb://admin:admin@localhost:27017/?authSource=admin" # <---- paste here
BUS_CONN_STRING="." # doesn't matter for dev environment
```

- Copy the certificates and private key from Organization 2's User and set peerEndpoint in `./sync/user_service_client.go`
```golang
// <<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem
const Certificate = "-----BEGIN CERTIFICATE-----\nMIICKjCCAdCgAwIBAgIQNpXmL8..."

// <<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/keystore/priv_sk
const PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCq..."

// <<PATH_TO_PROJECT>>/ticken-dev/test-pvtbc/test-network/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/tls/ca.crt
const TLSCertificated = "-----BEGIN CERTIFICATE-----\nMIICWDCCAf6gAwIBAgIRAL..."
```

```golang
func (usc *UserServiceClient) GetUserMembership(userID string) *UserMembership {
	return &UserMembership{
		MspID:          "Org2MSP",
		PeerEndpoint:   "localhost:9051", // <--- set this
		GatewayPeer:    "peer0.org2.example.com",
		Certificate:    Certificate,
		PrivateKey:     PrivateKey,
		TLSCertificate: TLSCertificated,
	}
}
```

- Build
```bash
$ go build
```

- Run
```bash
$ ./ticken-event-service
```

**Running Docker**

## Running tests

**Running specific package**

Use the following commnad to run the test in specific package
```
go test ./<paht_to_package>
```

**Running all tests**

Use the following commnad to run all tests in the project
```
go test ./...
```

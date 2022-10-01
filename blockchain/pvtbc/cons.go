package pvtbc

const globalPath = "/home/javier/Facultad/TPP/fabric-samples"
const cryptoPath = globalPath + "/test-network/organizations/peerOrganizations/org1.example.com"

const (
	mspID    = "Org1MSP"
	certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath  = cryptoPath + "/users/User1@org1.example.com/msp/keystore/priv_sk"
)

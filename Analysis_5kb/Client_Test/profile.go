package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"manufacturer": {
		CertPath:     "../Automobile-Network/organizations/peerOrganizations/manufacturer.auto.com/users/User1@manufacturer.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Automobile-Network/organizations/peerOrganizations/manufacturer.auto.com/users/User1@manufacturer.auto.com/msp/keystore/",
		TLSCertPath:  "../Automobile-Network/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.manufacturer.auto.com",
		MSPID:        "ManufacturerMSP",
	},

	"dealer": {
		CertPath:     "../Automobile-Network/organizations/peerOrganizations/dealer.auto.com/users/User1@dealer.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Automobile-network/organizations/peerOrganizations/dealer.auto.com/users/User1@dealer.auto.com/msp/keystore/",
		TLSCertPath:  "../Automobile-network/organizations/peerOrganizations/dealer.auto.com/peers/peer0.dealer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.dealer.auto.com",
		MSPID:        "DealerMSP",
	},

	"mvd": {
		CertPath:     "../Automobile-Network/organizations/peerOrganizations/mvd.auto.com/users/User1@mvd.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Automobile-Network/organizations/peerOrganizations/mvd.auto.com/users/User1@mvd.auto.com/msp/keystore/",
		TLSCertPath:  "../Automobile-Network/organizations/peerOrganizations/mvd.auto.com/peers/peer0.mvd.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.mvd.auto.com",
		MSPID:        "MvdMSP",
	},

	"manufacturer2": {
		CertPath:     "../Automobile-network/organizations/peerOrganizations/manufacturer.auto.com/users/User2@manufacturer.auto.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Automobile-network/organizations/peerOrganizations/manufacturer.auto.com/users/User2@manufacturer.auto.com/msp/keystore/",
		TLSCertPath:  "../Automobile-network/organizations/peerOrganizations/manufacturer.auto.com/peers/peer0.manufacturer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.manufacturer.auto.com",
		MSPID:        "ManufacturerMSP",
	},

	"minifab-manufacturer": {
		CertPath:     "../MinifabNetwork/vars/keyfiles/peerOrganizations/manufacturer.auto.com/users/Admin@manufacturer.auto.com/msp/signcerts/Admin@manufacturer.auto.com-cert.pem",
		KeyDirectory: "../MinifabNetwork/vars/keyfiles/peerOrganizations/manufacturer.auto.com/users/Admin@manufacturer.auto.com/msp/keystore/",
		TLSCertPath:  "../MinifabNetwork/vars/keyfiles/peerOrganizations/manufacturer.auto.com/peers/peer1.manufacturer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7003",
		GatewayPeer:  "peer1.manufacturer.auto.com",
		MSPID:        "manufacturer-auto-com",
	},

	"minifab-dealer": {
		CertPath:     "../MinifabNetwork/vars/keyfiles/peerOrganizations/dealer.auto.com/users/Admin@dealer.auto.com/msp/signcerts/Admin@dealer.auto.com-cert.pem",
		KeyDirectory: "../MinifabNetwork/vars/keyfiles/peerOrganizations/dealer.auto.com/users/Admin@dealer.auto.com/msp/keystore/",
		TLSCertPath:  "../MinifabNetwork/vars/keyfiles/peerOrganizations/dealer.auto.com/peers/peer1.dealer.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7004",
		GatewayPeer:  "peer1.dealer.auto.com",
		MSPID:        "dealer-auto-com",
	},

	"minifab-mvd": {
		CertPath:     "../MinifabNetwork/vars/keyfiles/peerOrganizations/mvd.auto.com/users/Admin@mvd.auto.com/msp/signcerts/Admin@mvd.auto.com-cert.pem",
		KeyDirectory: "../MinifabNetwork/vars/keyfiles/peerOrganizations/mvd.auto.com/users/Admin@mvd.auto.com/msp/keystore/",
		TLSCertPath:  "../MinifabNetwork/vars/keyfiles/peerOrganizations/mvd.auto.com/peers/peer1.mvd.auto.com/tls/ca.crt",
		PeerEndpoint: "localhost:7005",
		GatewayPeer:  "peer1.mvd.auto.com",
		MSPID:        "mvd-auto-com",
	},
}

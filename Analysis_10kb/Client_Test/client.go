package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var txnSuccess int

var txnFail int

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func submitTxnFn(organization string, channelName string, chaincodeName string, contractName string, txnType string, privateData map[string][]byte, txnName string, args ...string) bool {

	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	// defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(60*time.Second),
		client.WithSubmitTimeout(60*time.Second),
		client.WithCommitStatusTimeout(2*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	// defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	fmt.Printf("\n-->Submiting Transaction: %s with ID %s,\n", txnName, args[0])

	switch txnType {
	case "invoke":
		_, err := contract.SubmitTransaction(txnName, args...)

		if err != nil {
			txnFail += 1
			return false
		}else {
			txnSuccess +=1
		}
		return true

	// case "query":
	// 	evaluateResult, err := contract.EvaluateTransaction(txnName, args...)
	// 	if err != nil {
	// 		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	// 	}

	// 	// return fmt.Sprintf("*** Result:%s\n", result)
	// 	var result string
	// 	if isByteSliceEmpty(evaluateResult) {
	// 		result = string(evaluateResult)
	// 	} else {
	// 		result = formatJSON(evaluateResult)
	// 	}

	// 	// return fmt.Sprintf("*** Result:%s\n", result)
	// 	return result

	// case "private":
	// 	result, err := contract.Submit(
	// 		txnName,
	// 		client.WithArguments(args...),
	// 		client.WithTransient(privateData),
	// 	)

	// 	if err != nil {
	// 		panic(fmt.Errorf("failed to submit transaction: %w", err))
	// 	}

	// 	return fmt.Sprintf("*** Transaction committed successfully\n result: %s \n", result)

	// case "test":
	// 	startTime := time.Now()


	// 	for i := 0; i < totalTxns; i++ {
	// 		result, err := contract.EvaluateTransaction(txnName, args...)
	// 		if err != nil {
	// 			panic(fmt.Errorf("failed to submit transaction: %w", err))
	// 		}
	// 		fmt.Printf("Transaction %d submitted with result: %s\n", i+1, result)
	// 	}

	// 	endTime := time.Now()  // Track end time for TPS calculation
	// 	duration := endTime.Sub(startTime).Seconds()
	// 	fmt.Printf("Total Transactions: %d, TPS: %f\n", totalTxns, float64(totalTxns)/duration)


	}
	return false
}

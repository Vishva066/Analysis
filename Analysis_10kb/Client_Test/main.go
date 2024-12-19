package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type TxnCounters struct {
	success int
	failure int
	mu      sync.Mutex
}

type GatewayConfig struct {
	gateway  *client.Gateway
	network  *client.Network
	contract *client.Contract
}

func initializeGateway(organization, channelName, chaincodeName, contractName string) (*GatewayConfig, error) {
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(30*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %v", err)
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	return &GatewayConfig{
		gateway:  gw,
		network:  network,
		contract: contract,
	}, nil
}

func worker(
	jobs <-chan string,
	wg *sync.WaitGroup,
	counters *TxnCounters,
	gwConfig *GatewayConfig,
	txnName string,
) {
	defer wg.Done()

	for carKey := range jobs {
		fmt.Printf("\n-->Submitting Transaction: %s with ID %s\n", txnName, carKey)

		_, err := gwConfig.contract.SubmitTransaction(txnName, carKey, "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", "Nexon", "White", "KBA3", "22/07/2023")

		counters.mu.Lock()
		if err != nil {
			counters.failure++
			fmt.Printf("Transaction failed for %s: %v\n", carKey, err)
		} else {
			counters.success++
		}
		counters.mu.Unlock()
	}
}

func main() {
	var counters TxnCounters
	startTime := time.Now()
	totalTxns := 500
	numWorkers := 100
	jobs := make(chan string, totalTxns)

	// Initialize gateway once
	gwConfig, err := initializeGateway(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"CarContract",
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize gateway: %v", err))
	}
	// Ensure gateway is closed after we're done
	defer gwConfig.gateway.Close()

	var wg sync.WaitGroup

	// Start worker goroutines with shared gateway config
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(jobs, &wg, &counters, gwConfig, "CreateCar")
	}

	// Send transactions as jobs
	for i := 0; i < totalTxns; i++ {
		carKey := fmt.Sprintf("Car2-%d", i)
		jobs <- carKey
	}
	close(jobs)

	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	// fmt.Printf("\nResults:\n")
	// fmt.Printf("Successes: %d, Failures: %d\n", counters.success, counters.failure)
	// fmt.Printf("Total Transactions: %d, Duration: %.2f seconds, TPS: %.2f\n",
	//     totalTxns, duration, float64(totalTxns)/duration)
	file, err := os.Create("results.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}

	defer file.Close()

	fmt.Fprintf(file, "Performance Results:\n")
	fmt.Fprintf(file, "Successes: %d, Failures: %d\n", counters.success, counters.failure)
	fmt.Fprintf(file, "Total Transactions: %d, Duration: %.2f seconds, TPS: %.2f\n",
		totalTxns, duration, float64(totalTxns)/duration)

	if err != nil {
		panic(fmt.Sprintf("Failed to write to file: %v", err))
	}

	fmt.Println("Results have been written to results.txt")
}
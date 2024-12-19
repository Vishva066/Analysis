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
	totalTxns := 4000000
	numWorkers := 1000
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

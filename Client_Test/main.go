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
	id int,
	jobs <-chan string,
	wg *sync.WaitGroup,
	counters *TxnCounters,
	gwConfig *GatewayConfig,
	txnName string,
) {
	defer wg.Done()

	for carKey := range jobs {
		fmt.Printf("\n-->Submitting Transaction: %s with ID %s\n", txnName, carKey)

		_, err := gwConfig.contract.SubmitTransaction(txnName, carKey, "Tata Motors Limited is an Indian multinational automotive company, headquartered in Mumbai and part of the Tata Group. Prices indicated are Ex-showroom prices. Prices are subject to change without prior information at the discretion of Tata Motors. Contact your nearest Tata Motors dealer for exact prices.", "Nexon", "White", "The Apennine Colossus is a stone statue, approximately 11 metres (36 feet) tall, in the estate of Villa Demidoff (originally Villa di Pratolino) in Vaglia in Tuscany, Italy. A personification of the Apennine Mountains, the colossal figure was created by Giambologna, a Flemish-born Italian sculptor, in the late 1580s. The statue has the appearance of an elderly man crouched at the shore of a lake, squeezing the head of a sea monster through whose open mouth water originally emanated into the pond in front of the statue. The colossus is depicted naked, with stalactites in the thick beard and long hair to show the metamorphosis of man and mountain, blending his body with the surrounding nature. It is made of stone and plaster and the interior houses a series of chambers and caves on three levels. Initially, the back of the statue was protected by a structure resembling a cave, which was demolished around 1690 by the sculptor Giovanni Battista Foggini, who built a statue of a dragon to adorn the back of the colossus. The Italian sculptor Rinaldo Barbetti renovated the statue in 1876.", "22/07/2023")

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
		go worker(w, jobs, &wg, &counters, gwConfig, "CreateCar")
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

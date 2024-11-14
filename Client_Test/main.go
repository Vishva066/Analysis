package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {

// 	// use this functions to evaluate and submit txns
// 	// try calling these functions
// 	startTime := time.Now()
// 	for i := 0; i <= 50; i++ {
// 		carKey := fmt.Sprintf("Car1-%d", i)

// 		result := submitTxnFn(
// 			"manufacturer",
// 			"autochannel",
// 			"KBA-Automobile",
// 			"CarContract",
// 			"invoke",
// 			make(map[string][]byte),
// 			"CreateCar",
// 			1000,
// 			carKey,
// 			"Tata",
// 			"Harrier",
// 			"Black",
// 			"fac01",
// 			"25/10/2023",
// 		)

// 		fmt.Println(result)
// 	}

// 		endTime := time.Now()  // Track end time for TPS calculation
// 		duration := endTime.Sub(startTime).Seconds()
// 		fmt.Printf("Total Transactions: 50, TPS: %f\n", float64(50)/duration)

// 	// submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "test", make(map[string][]byte), "ReadCar",3000,"Car1-9")
// }

import (
    "fmt"
    "sync"
    "time"

    "github.com/hyperledger/fabric-gateway/pkg/client"   //Fabric Client package
)

type TxnCounters struct {
	success int
	failure int
	mu      sync.Mutex //Mutex to lock certain resources while one is updating
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup, counters *TxnCounters, organization string, channelName string, chaincodeName string ,contractName string, txnType string, txnName string) {
    defer wg.Done() //To execute this function after the current function block is excuted
    // This function tell the waiting list its work is done

    //Client code for fabric
    orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

    clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)

    Id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

    gw, err := client.Connect(
		Id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(30*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}

    network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// fmt.Printf("\n-->Submiting Transaction: %s with ID %s,\n", txnName, args[0])

    //Creating random key values for Car

    for carKey := range jobs {
        // result := submitTxnFn(
        //     "manufacturer",
        //     "autochannel",
        //     "KBA-Automobile",
        //     "CarContract",
        //     "invoke",
        //     make(map[string][]byte),
        //     "CreateCar",
        //     carKey,
        //     "Tata",
        //     "Harrier",
        //     "Black",
        //     "fac01",
        //     "25/10/2023",
        // )
        fmt.Printf("\n-->Submiting Transaction: %s with ID %s,\n", txnName, carKey)

        _, err := contract.SubmitTransaction(txnName, carKey, "Tata", "Harrier", "Black", "fac01", "25/10/2024")

        //Only one user per time to add the txn failures or success

		counters.mu.Lock()
		if err != nil {
			counters.failure++
		} else {
			// counters.failure++
            counters.success++
		}
        // Unlocking the resources

		counters.mu.Unlock()
    }
}


func main() {
	var counters TxnCounters //Defining the struct
    startTime := time.Now() //For TPS Calculation
    totalTxns := 1000 //Total no of transactions
    numWorkers := 100   // Total no of workers running parallely
    jobs := make(chan string, totalTxns)  //  channel to hold transaction jobs (queue)
    //This can hold the totalTxns no of strings
    var wg sync.WaitGroup //Set up wait group

    // Start 5 worker goroutines
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, &wg, &counters,"manufacturer", "autochannel","KBA-Automobile", "CarContract", "invoke", "CreateCar" )
    }

    // Send 50 transactions as jobs
    for i := 0; i < totalTxns; i++ {
        carKey := fmt.Sprintf("Car1-%d", i)
        jobs <- carKey  // Add job to the jobs channel
    }
    close(jobs)  // Close the channel to signal no more jobs

    wg.Wait()  // Wait for all workers to finish

    endTime := time.Now()
    duration := endTime.Sub(startTime).Seconds()
	fmt.Printf("Successes: %d, Failures: %d\n", counters.success, counters.failure)
    fmt.Printf("Total Transactions: %d, TPS: %f\n", totalTxns, float64(totalTxns)/duration)
}


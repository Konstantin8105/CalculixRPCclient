package clientCalculix

import (
	"fmt"
	"sync"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// CalculateForDat - calculation
func (c *ClientCalculix) CalculateForDat(inpBody []string) (datBody []string, err error) {

	type block struct {
		inp string
		dat string
	}
	var blockData []block
	blockChannel := make(chan block)
	errChannel := make(chan error)
	quitBlock := make(chan bool)
	quitErr := make(chan bool)

	var wg sync.WaitGroup

	go func() {
		for b := range blockChannel {
			blockData = append(blockData, b)
			fmt.Println("DAT client calculated task : ", len(blockData))
		}
		quitBlock <- true
	}()

	go func() {
		for e := range errChannel {
			err = fmt.Errorf("Error: %v\n%v", e, err)
		}
		quitErr <- true
	}()

	for _, inp := range inpBody {
		// Increment the WaitGroup counter.
		wg.Add(1)

		// Launch a goroutine to fetch the inp file.
		go func(inpFile string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
		BACK:
			client, err := c.getServer()
			defer func() {
				err2 := client.Close()
				if err2 != nil {
					errChannel <- fmt.Errorf("Errors:%v\n%v", err2, err)
				}
			}()
			if err != nil {
				if err.Error() == serverCalculix.ErrorServerBusy {
					goto BACK
				}
				errChannel <- err
			}

			var dat serverCalculix.DatBody
			err = client.Call("Calculix.ExecuteForDat", inpFile, &dat)
			if err != nil {
				if err.Error() == serverCalculix.ErrorServerBusy {
					goto BACK
				}
				errChannel <- err
			}
			blockChannel <- block{inp: inpFile, dat: dat.A}
		}(inp)
	}
	// Wait for all inp body
	wg.Wait()

	// Close all opened channels
	close(errChannel)
	close(blockChannel)

	<-quitBlock
	<-quitErr

	//repair sequene of result dat
	for _, inp := range inpBody {
		for _, block := range blockData {
			if inp == block.inp {
				datBody = append(datBody, block.dat)
				goto NewInp
			}
		}
	NewInp:
	}

	return datBody, err
}

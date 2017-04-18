package clientCalculix

import (
	"fmt"
	"sync"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// CalculateForDat - calculation
func (c *ClientCalculix) CalculateForDat(inpBody []string) (datBody []string, err error) {

	inpMap := make(map[int]string)
	datMap := make(map[int]string)

	type block struct {
		index int
		value string
	}
	blockChannel := make(chan block)
	errChannel := make(chan error)
	quitBlock := make(chan bool)
	quitErr := make(chan bool)

	go func() {
		for b := range blockChannel {
			datMap[b.index] = b.value
			fmt.Printf("DAT client calculated task : %4v of %4v\n", len(datMap), len(inpMap))
		}
		quitBlock <- true
	}()

	go func() {
		for e := range errChannel {
			err = fmt.Errorf("Error: %v\n%v", e, err)
		}
		quitErr <- true
	}()

	errGlobal := err

	var wg sync.WaitGroup

	for index, inp := range inpBody {
		// add to inp map
		inpMap[index] = inp

		// Increment the WaitGroup counter.
		wg.Add(1)

		// Launch a goroutine to fetch the inp file.
		go func(index int, inpFile string) {
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
				return
			}
			if errGlobal != nil {
				return
			}

			var dat serverCalculix.DatBody
			err = client.Call("Calculix.ExecuteForDat", inpFile, &dat)
			if err != nil {
				if err.Error() == serverCalculix.ErrorServerBusy {
					goto BACK
				}
				errChannel <- err
				return
			}
			blockChannel <- block{index: index, value: dat.A}
		}(index, inp)
	}
	// Wait for all inp body
	wg.Wait()

	// Close all opened channels
	close(errChannel)
	close(blockChannel)

	<-quitBlock
	<-quitErr
	close(quitBlock)
	close(quitErr)

	//repair sequene of result dat
	size := len(inpBody)
	for index := 0; index < size; index++ {
		for k, v := range datMap {
			if index == k {
				datBody = append(datBody, v)
				goto NewInp
			}
		}
	NewInp:
	}

	return datBody, err
}

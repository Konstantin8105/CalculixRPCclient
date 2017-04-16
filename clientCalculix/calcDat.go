package clientCalculix

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// CalculateForDat - calculation
func (c *ClientCalculix) CalculateForDat(inpBody []string) (datBody []string, err error) {
	for _, inp := range inpBody {
	BACK:
		client, err := c.getServer()
		defer func() {
			err2 := client.Close()
			if err2 != nil {
				err = fmt.Errorf("Errors:%v\n%v", err2, err)
			}
		}()
		if err.Error() == serverCalculix.ErrorServerBusy {
			goto BACK
		}
		if err != nil {
			return datBody, err
		}

		var dat serverCalculix.DatBody
		err = client.Call("Calculix.ExecuteForDat", inp, &dat)
		if err != nil {
			return datBody, err
		}
		datBody = append(datBody, dat.A)
	}
	return datBody, nil
}

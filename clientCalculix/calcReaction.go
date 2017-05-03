package clientCalculix

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Konstantin8105/Convert-INP-to-STD-format/calculixResult"
)

// CalculateForReactionLoad - calculation reaction
func CalculateForReactionLoad() {
	file := "example.dat"

	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}
	// open file
	inFile, err := os.Open(file)
	if err != nil {
		return
	}
	defer func() {
		errFile := inFile.Close()
		if errFile != nil {
			if err != nil {
				err = fmt.Errorf("%v ; %v", err, errFile)
			} else {
				err = errFile
			}
		}
	}()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	//
	//
	//
	//
	//
	//
	//
	// forces (fx,fy,fz) for set FIX and time  0.4000000E-01
	// forces (fx,fy,fz) for set LOAD and time  0.2000000E-01
	// 204  3.485854E+00  1.025290E+01  3.092803E+01

	v, err := calculixResult.SupportForcesSummary(lines)
	//v, err := calculixResult.SupportForces(lines)
	for _, vv := range v {
		if vv.NodeName == "FIX" {
			fmt.Println("vv = ", vv)
		}
	}
	fmt.Println("ERROR = ", err)
}

package clientCalculix

import (
	"fmt"
	"strconv"
	"strings"
)

// CalculateForBuckle - calculation
func (c *ClientCalculix) CalculateForBuckle(inpBody []string) (factors []float64, err error) {
	dats, err := c.CalculateForDat(inpBody)
	if err != nil {
		return factors, err
	}
	for _, dat := range dats {
		factor, err := getBucklingFactor(dat)
		if err != nil {
			return factors, err
		}
		factors = append(factors, factor)
	}
	return factors, nil
}

func getBucklingFactor(dat string) (factor float64, err error) {
	lines := strings.Split(dat, "\n")

	bucklingHeader := "B U C K L I N G   F A C T O R   O U T P U T"
	var found bool
	var numberLine int
	amountBuckling := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !found {
			// empty line
			if len(line) == 0 {
				continue
			}
			if len(line) != len(bucklingHeader) {
				continue
			}
			if line == bucklingHeader {
				found = true
			}
		} else {
			numberLine++
			if numberLine >= 5+amountBuckling {
				break
			}
			if numberLine >= 5 {
				m, f, err := parseBucklingFactor(line)
				if err != nil {
					return -1., err
				}
				if m != numberLine-4 {
					return -1., fmt.Errorf("Wrong MODE NO: %v (%v) in line: %v", m, numberLine-4, line)
				}
				return f, nil
			}
		}
	}
	return -1., fmt.Errorf("Cannot found in dat")
}

// Example:
//      4   0.4067088E+03
func parseBucklingFactor(line string) (mode int, factor float64, err error) {
	s := strings.Split(line, "   ")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}

	var index int

	for index = 0; index < len(s); index++ {
		if len(s[index]) == 0 {
			continue
		}
		i, err := strconv.ParseInt(s[index], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
		}
		mode = int(i)
		break
	}

	for index++; index < len(s); index++ {
		if len(s[index]) == 0 {
			continue
		}
		factor, err = strconv.ParseFloat(s[index], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
		}
		break
	}
	return mode, factor, nil
}

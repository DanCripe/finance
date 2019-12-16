package main

import (
	"strconv"
	"strings"
)

var stockData = `-3.06
30.23
7.49
9.97
1.33
37.2
22.68
33.10
28.34
20.89
-9.03
-11.85
-21.97
28.36
10.74
4.83
15.61
5.48
-36.55
25.94
14.2
2.10
15.89
32.15
13.52
1.38
11.77
21.61
-4.23`

var bondData = `6.24
15.00
9.36
14.21
-8.04
23.58
1.43
9.94
14.92
-8.25
16.66
5.57
15.12
0.38
4.49
2.87
1.96
10.21
20.10
-11.12
8.46
16.04
2.97
-9.10
10.75
1.28
0.69
2.80
-0.02`

var inflationData = `6.1
3.1
2.9
2.7
2.7
2.5
3.3
1.7
1.6
2.7
3.4
1.6
2.4
1.9
3.3
3.4
2.5
4.1
0.1
2.7
1.5
3.0
1.7
1.5
0.8
0.7
2.1
2.1
1.9`

var stocks []float64
var bonds []float64
var inflation []float64

func init() {
	valList := strings.Split(stockData, "\n")
	for _, val := range valList {
		x, _ := strconv.ParseFloat(val, 64)
		stocks = append(stocks, x/100.0)
	}

	valList = strings.Split(bondData, "\n")
	for _, val := range valList {
		x, _ := strconv.ParseFloat(val, 64)
		bonds = append(bonds, x/100.0)
	}
	valList = strings.Split(inflationData, "\n")
	for _, val := range valList {
		x, _ := strconv.ParseFloat(val, 64)
		inflation = append(inflation, x/100.0)
	}
}

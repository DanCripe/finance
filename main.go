package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	// "strings"

	"github.com/aybabtme/uniplot/histogram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)
var balance = flag.Float64("balance", 50.0, "Stock Percent")
var adjustment = flag.Float64("adjustment", 0.0, "Stock Percent Adjustment Per Year")

type value struct {
	idx int
	val float64
}

func main() {
	flag.Parse()

	*balance /= 100.0
	*adjustment /= 100.0

	fmt.Printf("-------------------------------------------\n")
	fmt.Printf("Balance: %6.2f%%\n", *balance*100.0)
	fmt.Printf("Adjust:  %6.2f%%\n", *adjustment*100.0)

	highest := float64(0.0)
	lowest := float64(1000000000.0)
	lidx := -1
	hidx := -1
	lyear := 110
	var avgReturns []float64
	var avgInflations []float64
	var results []float64
	var successCount int
	var attempts = 10000
	for i := 0; i < attempts; i++ {
		val, avgReturn, avgInflation, year, success := try(i, *balance, *adjustment, false)
		if val < lowest || year < lyear {
			lowest = val
			lidx = i
			lyear = year
		}
		if val > highest {
			highest = val
			hidx = i
		}

		if success {
			successCount++
		}

		avgReturns = append(avgReturns, avgReturn*100.0)
		avgInflations = append(avgInflations, avgInflation*100.0)
		results = append(results, val)
	}

	values := make([]float64, len(results))
	copy(values, results)

	sort.Float64Slice(results).Sort()

	fmt.Printf("Success Rate: %7.3f%%\n", float64(successCount)*100.0/float64(attempts))
	fmt.Printf("Low:  %6d %13.2f %d\n", lidx, lowest, lyear)

	fmt.Printf("High: %6d %13.2f\n", hidx, highest)
	fmt.Printf("Med:         %13.2f\n", results[len(results)/2])
	for idx, value := range values {
		if value == results[len(results)/2] {
			fmt.Printf("Median index is %d\n", idx)
			v, _, _, _, _ := try(idx, *balance, *adjustment, true)
			fmt.Printf("Median value is %13.2f\n", v)
		}
	}

	// low index
	// try(lidx, *balance, *adjustment, true)

	val := average(avgReturns)
	p.Printf("Average Return:    %11.2f\n", val)
	val = average(avgInflations)
	p.Printf("Average Inflation: %11.2f\n", val)

	hist := histogram.Hist(100, results)
	histogram.Fprintf(os.Stdout, hist, histogram.Linear(5), func(v float64) string {
		return p.Sprintf("%16.2f", v)
	})
	// p.Printf("Average Returns Histogram\n")
	// hist := histogram.Hist(30, avgReturns)
	// histogram.Fprint(os.Stdout, hist, histogram.Linear(5))

	// p.Printf("Average Inflation\n")
	// hist = histogram.Hist(30, avgInflations)
	// histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
}

func try(seed int, mix, adjust float64, display bool) (float64, float64, float64, int, bool) {
	r := rand.New(rand.NewSource(int64(seed)))

	cfg, err := readConfig()

	if err != nil {
		p.Printf("Error reading configuration: %v\n", err)
		os.Exit(1)
	}

	if len(stocks) != len(bonds) || len(stocks) != len(inflation) {
		p.Printf("Lengths are not correct: %d %d %d\n", len(stocks), len(bonds), len(inflation))
		os.Exit(1)
	}

	data := convert(cfg)

	var returns []float64
	var inflations []float64

	final := false
	ss := float64(0.0)
	sss := float64(0.0)
	if display {
		p.Printf("%4.4s %7.7s %13s %13s %13s %13s %13s %13s %4s   %5s   %4s\n", "YEAR", "AGES", "ORDINARY", "CAP GAINS", "TAX FREE", "EXPENSES", "SOC SEC", "SPOUSE SS", "PCT", "RTN", "INF")
	}

	var year int
	var wr float64

	for year = data.StartYear; year <= data.EndYear && !final; year++ {
		idx1 := r.Uint32() % uint32(len(stocks))
		idx2 := r.Uint32() % uint32(len(stocks))
		idx3 := r.Uint32() % uint32(len(stocks))

		avgReturn := calcAvg(stocks[idx1], bonds[idx2], mix)
		inflation := inflation[idx3]

		remaining := data.Expenses

		if data.Age >= data.Income.SocialSecurityAge {
			remaining = remaining - (data.Income.SocialSecurity * (1.0 - data.TaxRate/2.0))
			ss = data.Income.SocialSecurity
		}

		if data.SpouseAge >= data.Income.SpouseSocialSecurityAge {
			remaining = remaining - (data.Income.SpouseSocialSecurity * (1.0 - data.TaxRate/2.0))
			sss = data.Income.SpouseSocialSecurity
		}

		if ss == 0.0 {
			// always use up to Standard Deduction as ordinary income
			if data.Assets.OrdinaryIncome < data.StdDeduction {
				remaining -= data.Assets.OrdinaryIncome
				data.Assets.OrdinaryIncome = 0.0
			} else {
				remaining -= data.StdDeduction
				data.Assets.OrdinaryIncome -= data.StdDeduction
			}
		}

		if data.Assets.CapitalGains == 0.0 && data.Assets.TaxFree == 0.0 {
			remaining = remaining * (1.0 + data.TaxRate)
		}

		wr = ((remaining + data.StdDeduction) * 100.0) / (data.Assets.OrdinaryIncome + data.Assets.TaxFree + data.Assets.CapitalGains)

		if display {
			p.Printf("%s %3d %3d %s %s %s %13.2f %13.2f %13.2f %4.2f %7.3f %6.3f\n", fmt.Sprintf("%d", year), data.Age, data.SpouseAge, colorize(data.Assets.OrdinaryIncome, green, red), colorize(data.Assets.CapitalGains, green, red), colorize(data.Assets.TaxFree, green, red), data.Expenses, ss, sss, wr, avgReturn*100.0, inflation*100.0)
		}

		if year+1 == data.EndYear {
			avgReturn := average(returns)
			avgInflation := average(inflations)

			return data.Assets.OrdinaryIncome + data.Assets.CapitalGains + data.Assets.TaxFree, avgReturn, avgInflation, year + 1, !final
		}

		if remaining <= data.Assets.CapitalGains {
			data.Assets.CapitalGains -= remaining
			remaining = 0.0
		} else {
			remaining -= data.Assets.CapitalGains
			data.Assets.CapitalGains = 0.0
		}

		if remaining <= data.Assets.TaxFree {
			data.Assets.TaxFree -= remaining
			remaining = 0.0
		} else {
			remaining -= data.Assets.TaxFree
			data.Assets.TaxFree = 0.0
		}

		if remaining <= data.Assets.OrdinaryIncome {
			data.Assets.OrdinaryIncome -= remaining
			remaining = 0.0
		} else {
			data.Assets.OrdinaryIncome = 0.0
			final = true
		}

		returns = append(returns, avgReturn)
		inflations = append(inflations, inflation)
		data.StdDeduction *= inflation + 1.0

		avgReturn += 1.0
		data.Assets.OrdinaryIncome *= avgReturn
		data.Assets.CapitalGains *= avgReturn
		data.Assets.TaxFree *= avgReturn

		data.Income.SocialSecurity *= inflation + 1.0
		data.Income.SpouseSocialSecurity *= inflation + 1.0

		data.Expenses *= inflation + 1.0

		if mix >= adjust {
			mix -= adjust
		} else {
			mix = 0.0
		}

		data.Age++
		data.SpouseAge++
	}

	avgReturn := average(returns)
	avgInflation := average(inflations)

	if display {
		p.Printf("%11.2f %11.2f\n", avgReturn*100.0, avgInflation*100.0)
	}

	return data.Assets.OrdinaryIncome + data.Assets.CapitalGains + data.Assets.TaxFree, avgReturn, avgInflation, year, !final
}

func calcAvg(a, b float64, mix float64) float64 {
	// A1 = A0 * (1 + returnA)
	// B1 = B0 * (1 + returnB)
	// ((A1 + B1) - (A0 + B0)) / (A0 + B0)
	// ((A0 * (1 + returnA) + B0 * (1 + returnB)) - (A0 + B0)) / (A0 + B0)
	// (A0 + A0 * returnA + B0 + B0 * returnB - A0 - B0) / (A0 + B0)
	// (A0 * returnA + B0 * returnB) / (A0 + B0)
	return (a * mix) + (b * (1.00 - mix))
}

func average(items []float64) float64 {
	return sum(items) / float64(len(items))
}

func sum(items []float64) float64 {
	val := float64(0.0)
	for _, item := range items {
		val += item
	}
	return val
}

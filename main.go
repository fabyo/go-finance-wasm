package main

import (
	"math"
	"syscall/js"
)

func main() {
	// registra a função calcFinance no objeto global JS
	js.Global().Set("calcFinance", js.FuncOf(calcFinance))

	// mantém o programa vivo (senão o Go sai e o WASM morre)
	select {}
}

// calcFinance(propertyValue, downPayment, annualRate, months)
func calcFinance(this js.Value, args []js.Value) any {
	if len(args) < 4 {
		return map[string]any{
			"error": "expected 4 arguments: propertyValue, downPayment, annualRate, months",
		}
	}

	propertyValue := args[0].Float()
	downPayment := args[1].Float()
	annualRate := args[2].Float()
	months := args[3].Int()

	if propertyValue <= 0 {
		return map[string]any{"error": "propertyValue must be > 0"}
	}
	if months <= 0 {
		return map[string]any{"error": "months must be > 0"}
	}
	if downPayment < 0 || downPayment >= propertyValue {
		return map[string]any{"error": "downPayment must be >= 0 and < propertyValue"}
	}

	principal := propertyValue - downPayment
	monthlyRate := annualRate / 12.0 / 100.0

	var monthlyPayment float64
	if monthlyRate == 0 {
		monthlyPayment = principal / float64(months)
	} else {
		// PMT = P * i / (1 - (1+i)^(-n))
		pow := math.Pow(1+monthlyRate, float64(-months))
		monthlyPayment = principal * monthlyRate / (1 - pow)
	}

	totalPaid := monthlyPayment * float64(months)
	totalInterest := totalPaid - principal

	return map[string]any{
		"propertyValue":   propertyValue,
		"downPayment":     downPayment,
		"principal":       principal,
		"annualRate":      annualRate,
		"monthlyRate":     monthlyRate,
		"months":          months,
		"monthlyPayment":  monthlyPayment,
		"totalPaid":       totalPaid,
		"totalInterest":   totalInterest,
	}
}

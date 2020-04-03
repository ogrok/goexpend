package state

import (
	"github.com/adaminoue/goexpend/src/models"
	"math"
	"strconv"
	"time"
)

const maxReportWidth = 100
var bufferCoefficient = 1.05

func ShowFullReport() error {
	items, err := GetAllActiveItems()

	if err != nil {
		return err
	}

	var reportItems []models.ReportViewItem

	for _, i := range items {
		reportItems = append(reportItems, i.ToReport())
	}

	config, err := GetConfig()

	if err != nil {
		return err
	}

	var report = models.Report{
		Items:          reportItems,
		Income:         config.Income,
		Year:           config.CurrentYear,
		Month:          config.CurrentMonth,
	}

	widths := report.CalculateColWidths()
	report.CalculateTotals()

	if widths.TotalWidth() > maxReportWidth {
		// TODO for table: handle max width reduction, wrapping, etc.
	}

	// TODO detail table generation should go here

	bufferAmount := int(math.Round(float64(report.TotalRemaining) * bufferCoefficient))
	bufferToView := strconv.FormatFloat(bufferCoefficient, 'f', -2, 64)

	println("\n" + time.Month(report.Month).String() + " " + strconv.Itoa(report.Year) + " Report\n")
	println(strconv.Itoa(report.TotalRealized) + " / " + strconv.Itoa(report.TotalAccrued) + " realized expenses")
	println(strconv.Itoa(report.TotalRemaining) + " remaining\n")
	println("goexpend recommends keeping at least " + strconv.Itoa(bufferAmount) + " in your account.")
	println("This is " + bufferToView + " times more than the remaining balance.\n")

	return nil
}

package models

import "strconv"

// denotes number of extra spaces after each column
const extraSpace = 1

type Report struct {
	Items          []ReportViewItem
	Income         int
	Year           int
	Month          int

	TotalAccrued   int
	TotalRealized  int
	TotalRemaining int
}

type ReportViewItem struct {
	Name        string
	Category    string
	Description string
	Accrued     string
	Realized    string
	Mutable     rune
	SideNote	string
}

type ReportMaxWidths struct {
	NameWidth        int
	CategoryWidth    int
	DescriptionWidth int
	AccruedWidth     int
	RealizedWidth    int
	MutableWidth     int
	SideNoteWidth    int
}

func (r *Report) CalculateTotals() {
	accrued, realized, remaining := 0, 0, 0

	for _, i := range r.Items {
		a, _ := strconv.Atoi(i.Accrued)
		accrued += a

		b, _ := strconv.Atoi(i.Realized)
		realized += b

		remaining += a - b
	}

	r.TotalAccrued, r.TotalRealized, r.TotalRemaining = accrued, realized, remaining
}

func (r *Report) CalculateColWidths() ReportMaxWidths {
	nameMax, catMax, descMax, AccMax, RealMax, SideMax := 0, 0, 0, 0, 0, 0

	for _, i := range r.Items {
		if len(i.Name) > nameMax {
			nameMax = len(i.Name)
		}
		if len(i.Category) > catMax {
			catMax = len(i.Category)
		}
		if len(i.Description) > descMax {
			descMax = len(i.Description)
		}
		if len(i.Accrued) > AccMax {
			AccMax = len(i.Accrued)
		}
		if len(i.Realized) > RealMax {
			RealMax = len(i.Realized)
		}
		if len(i.SideNote) > SideMax {
			SideMax = len(i.SideNote)
		}
	}

	return ReportMaxWidths{
		NameWidth:        nameMax + extraSpace,
		CategoryWidth:    catMax + extraSpace,
		DescriptionWidth: descMax,              // no extra space bc of wrapping
 		AccruedWidth:     AccMax + extraSpace,
		RealizedWidth:    RealMax + extraSpace,
		MutableWidth:     1,                    // always single-character width here
		SideNoteWidth:    SideMax + extraSpace,
	}
}

func (w *ReportMaxWidths) TotalWidth() int {
	return w.AccruedWidth + w.CategoryWidth + w.DescriptionWidth +
		w.MutableWidth + w.NameWidth + w.RealizedWidth + w.SideNoteWidth
}

func (a *ActiveItem) ToReport() ReportViewItem {
	var mutableRune rune

	if a.Mutable {
		mutableRune = 'M'
	} else {
		mutableRune = ' '
	}

	return ReportViewItem{
		Name:        a.Name,
		Category:    a.Category,
		Description: a.Description,
		Accrued:     strconv.Itoa(a.Accrued),
		Realized:    strconv.Itoa(a.Realized),
		Mutable:     mutableRune,
		SideNote:    generateSideNote(a),
	}
}

// currently only shows overbudget and extra-accrual comments
func generateSideNote(a *ActiveItem) string {
	overspend := a.Realized - a.Accrued

	if overspend > 0 {
		return strconv.Itoa(overspend) + " over budget"
	}

	extraAccrual := a.Accrued - a.Amount

	if extraAccrual > 0 {
		return strconv.Itoa(extraAccrual) + " extra accrued"
	}

	return ""
}
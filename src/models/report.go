package models

import "strconv"

// denotes number of extra spaces after each column
const extraSpace = 1

type Report struct {
	Items  []ReportViewItem
	Income int
	Year   int
	Month  int
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
		DescriptionWidth: descMax, // no extra space bc of wrapping
 		AccruedWidth:     AccMax + extraSpace,
		RealizedWidth:    RealMax + extraSpace,
		MutableWidth:     1,
		SideNoteWidth:    SideMax + extraSpace,
	}
}

func (a *ActiveItem) ToReport() ReportViewItem {
	var mutableRune rune

	if a.Mutable {
		mutableRune = 'Y'
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

// TODO the logic that does side notes from item context
func generateSideNote(a *ActiveItem) string {
	return ""
}
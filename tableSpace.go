package main

import "github.com/icza/gowut/gwu"

// NewHSpace create a horezontal space that can be added where you like. For example in a table cell.
func NewHSpace(width int) gwu.Comp {
	l := gwu.NewLabel("")
	l.Style().SetDisplay(gwu.DisplayBlock).SetWidthPx(width)
	return l
}

// NewVSpace create a virtical space that can be added where you like. For example in a table cell.
func NewVSpace(height int) gwu.Comp {
	l := gwu.NewLabel("")
	l.Style().SetDisplay(gwu.DisplayBlock).SetHeightPx(height)
	return l
}

// NewSpace create a space that can be added where desired. For example into a table cell.
func NewSpace(width, height int) gwu.Comp {
	l := gwu.NewLabel("")
	l.Style().SetDisplay(gwu.DisplayBlock).SetSizePx(width, height)
	return l
}

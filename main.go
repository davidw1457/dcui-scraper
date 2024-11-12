package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("DCUI Scraper")

	updateButton := widget.NewButton("Update DCUI Database", updateDatabase)
	filterText := canvas.NewText("Filters", color.White)
	titleFilterButton := widget.NewButton("Title", titleFilter)
	dateFilterButton := widget.NewButton("Date Range", dateFilter)

	leftPane := container.New(layout.NewVBoxLayout(), updateButton, filterText, titleFilterButton, dateFilterButton)
	// TODO: Add a widget.NewList to hold filter contents 
	centerPane := container.New(layout.NewVBoxLayout(), canvas.NewText("Filter Options:", color.White))
	// TODO: Add a widget.NewTable to hold filter output
	rightPane := container.New(layout.NewVBoxLayout(), canvas.NewText("Filter Output:", color.White))

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), leftPane, widget.NewSeparator(), centerPane, widget.NewSeparator(), rightPane, layout.NewSpacer()))

	myWindow.ShowAndRun()
}

func updateDatabase() {
	fmt.Println("updateButton pressed")
}

func titleFilter() {
	fmt.Println("titleFilterButton pressed")
}

func dateFilter() {
	fmt.Println("dateFilter pressed")
}

func buildContent() {
	
}
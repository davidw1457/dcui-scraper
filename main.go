package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/davidw1457/dcui-scraper/database"
)

const userRWX = 0o700

var mainLog *log.Logger //nolint:gochecknoglobals

func main() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		initError(err.Error())
		os.Exit(1)
	}

	sep := string(os.PathSeparator)

	logPath := userHome + sep + ".dcui" + sep + "logs"

	_, err = os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(logPath, userRWX)
			if err != nil {
				initError(err.Error())
				os.Exit(1)
			}
		} else {
			initError(err.Error())
			os.Exit(1)
		}
	}

	logFile, err := os.OpenFile(logPath+sep+"dcui-scraper.log",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, userRWX)
	if err != nil {
		initError(err.Error())
		os.Exit(1)
	}

	mainLog = log.New(logFile, "main: ", log.LstdFlags)

	mainLog.Println("opening backend database")

	dbase, err := database.New()
	if err != nil {
		mainLog.Fatalln("unable to open database")
	}

	defer dbase.Close()

	myApp := app.New()
	myWindow := myApp.NewWindow("DCUI Scraper")

	updateButton := widget.NewButton("Update DCUI Database", func() {
		err := dbase.RefreshDatabase()
		if err != nil {
			// TODO: Update UI if there is an error
			mainLog.Println(err)
		}
	})
	filterText := canvas.NewText("Filters", color.White)
	titleFilterButton := widget.NewButton("Title", titleFilter)
	dateFilterButton := widget.NewButton("Date Range", dateFilter)

	leftPane := container.New(layout.NewVBoxLayout(), updateButton, filterText, titleFilterButton, dateFilterButton)
	// TODO: Add a widget.NewList to hold filter contents
	centerPane := container.New(layout.NewVBoxLayout(), canvas.NewText("Filter Options:", color.White))
	// TODO: Add a widget.NewTable to hold filter output
	rightPane := container.New(layout.NewVBoxLayout(), canvas.NewText("Filter Output:", color.White))

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), leftPane, widget.NewSeparator(), centerPane,
		widget.NewSeparator(), rightPane, layout.NewSpacer()))

	myWindow.ShowAndRun()
}

func titleFilter() {
	fmt.Println("titleFilterButton pressed")
}

func dateFilter() {
	fmt.Println("dateFilter pressed")
}

func buildContent() {
	// TODO: Move UI building here, maybe?
}

func initError(err string) {
	myApp := app.New()
	myWindow := myApp.NewWindow("ERROR")
	myWindow.SetContent(canvas.NewText(err, color.White))
	myWindow.ShowAndRun()
}

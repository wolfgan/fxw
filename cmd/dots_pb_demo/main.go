package main

import (
	"fmt"
	xwidget "fxw"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Calendar")

	progress := binding.NewFloat()
	pbv := xwidget.NewDotsProgressBarWithData(progress)
	pbv.SetMinSize(fyne.NewSize(200, 200))
	data := binding.NewString()

	pbi := xwidget.NewDotsProgressBarInfinity()
	pbi.SetMinSize(fyne.NewSize(200, 200))

	go func() {
		direction := 0.05
		for {
			val, _ := progress.Get()
			switch {
			case val >= 1.0:
				direction = -0.05
			case val <= 0.0:
				direction = 0.05
			}

			progress.Set(val + direction)
			data.Set(fmt.Sprintf("Progress: %.1f", (val+direction)*100))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	w.SetContent(
		container.NewHBox(
			container.NewVBox(pbv,
				container.NewCenter(
					widget.NewLabelWithData(data),
				),
			),
			container.NewVBox(pbi,
				container.NewGridWithColumns(2,
					widget.NewButton("-", func() {
						pbi.SetSpeed(pbi.Speed() + 0.1)
					}),
					widget.NewButton("+", func() {
						pbi.SetSpeed(pbi.Speed() - 0.1)
					}),
				),
			),
		),
	)
	w.ShowAndRun()
}

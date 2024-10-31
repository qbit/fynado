package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

var (
	ffffea = color.RGBA{R: 255, G: 255, B: 234, A: 255}
	sixes  = color.RGBA{R: 204, G: 204, B: 204, A: 255}
	blue   = color.RGBA{R: 74, G: 144, B: 226, A: 255}
)

func main() {
	debug := flag.Bool("d", false, "enable debugging (work / break time become 10 seconds)")
	makeIcon := flag.Bool("i", false, "draw an icon and exit.")
	flag.Parse()
	enabled := true
	rounds := 0
	workTime := 25 * 60
	breakTime := 5 * 60

	if *debug {
		workTime = 10
		breakTime = 10
	}

	icon := CountIcon{
		Enabled: enabled,
	}

	var desk desktop.App

	disable := func() {
		enabled = false
		icon.Enabled = enabled
		desk.SetSystemTrayIcon(icon.Draw(0.0))
		systray.SetTooltip("Disabled")
	}

	enable := func() {
		enabled = true
		icon.Enabled = enabled
	}

	if !*makeIcon {
		a := app.New()
		driver := fyne.CurrentApp().Driver().(desktop.Driver)
		w := driver.CreateSplashWindow()

		if castDesk, ok := a.(desktop.App); ok {
			desk = castDesk
			desk.SetSystemTrayIcon(icon.Draw(1.0))
			desk.SetSystemTrayMenu(
				fyne.NewMenu("fynado",
					fyne.NewMenuItem("Enable", func() {
						enable()
					}),
					fyne.NewMenuItem("Disable", func() {
						disable()
					}),
				),
			)
		} else {
			log.Fatalln("can't initialize deskto.App")
		}

		a.Lifecycle().SetOnExitedForeground(func() {
			time.Sleep(time.Second * 3)
			w.RequestFocus()
		})

		rect := canvas.NewRectangle(color.Black)
		infoRect := canvas.NewRectangle(ffffea)

		reminderText := canvas.NewText("Time for a break!", color.Black)
		reminderText.TextSize = 80
		reminderText.TextStyle = fyne.TextStyle{Bold: true}
		reminderText.Alignment = fyne.TextAlignCenter

		timerText := canvas.NewText("", blue)
		timerText.TextSize = 50
		timerText.TextStyle = fyne.TextStyle{Bold: true}
		timerText.Alignment = fyne.TextAlignCenter

		roundsText := canvas.NewText("", color.Black)
		roundsText.TextSize = 50
		roundsText.Alignment = fyne.TextAlignCenter

		w.Resize(fyne.NewSize(1000, 600))
		w.SetContent(
			container.NewStack(
				rect,
				container.NewPadded(infoRect,
					container.NewGridWithRows(4,
						container.NewBorder(nil, reminderText, nil, nil),
						container.NewBorder(timerText, nil, nil, nil),
						container.NewBorder(roundsText, nil, nil, nil),
						container.NewBorder(nil, widget.NewButton("Extended Break", func() {
							rounds = 0
							w.Hide()
							disable()
						}), nil, nil),
					)),
			),
		)

		ctrlQ := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
		ctrlW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierControl}
		w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
			a.Quit()
		})
		w.Canvas().AddShortcut(ctrlW, func(shortcut fyne.Shortcut) {
			w.Hide()
		})

		go func() {
			duration := workTime
			w.Hide()

			log.Println("work time starting")

			for {
				if !enabled {
					time.Sleep(time.Second)
					continue
				}
				duration--
				timerText.Text = fmt.Sprintf("%02d:%02d", duration/60, duration&60)
				timerText.Refresh()

				w.Content().Refresh()
				systray.SetTooltip(timerText.Text)

				remainingPct := float64((workTime - duration)) / float64(workTime) * 1

				desk.SetSystemTrayIcon(icon.Draw(remainingPct))

				if duration == 0 {
					rounds++

					if rounds > 3 {
						roundsText.Text = fmt.Sprintf("%d rounds! Take an extended break!", rounds)
					} else {
						roundsText.Text = fmt.Sprintf("Rounds: %d", rounds)
					}
					roundsText.Refresh()

					log.Printf("break %d starting", rounds)
					showBreak(breakTime, w, timerText)
					duration = workTime
				}
				time.Sleep(time.Second)
			}
		}()

		a.Run()
	} else {
		data := icon.Draw(.73).Content()
		os.WriteFile("icon.png", data, 0644)
		os.Exit(0)
	}
}

func showBreak(breakTime int, w fyne.Window, t *canvas.Text) {

	w.Show()
	duration := breakTime
	for {
		t.Text = fmt.Sprintf("%02d:%02d", duration/60, duration%60)
		t.Refresh()

		systray.SetTooltip(t.Text)

		if duration == 0 {
			w.Hide()
			return
		}
		duration--
		time.Sleep(time.Second)
	}
}

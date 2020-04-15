package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/google/gousb"
	"github.com/ramin0/live/ifuse/ifuse"
	"github.com/skratchdot/open-golang/open"
)

const iDeviceVendorProduct = "05ac:12a8"

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("iFuse")

	serialNumbers, err := getSerialNumbers()
	if err != nil {
		log.Fatalf("Failed to get serial numbers: %v", err)
	}
	if len(serialNumbers) == 0 {
		log.Fatal("No iDevices found")
	}

	w.SetContent(widget.NewVBox(
		widget.NewLabel("Select iDevice:"),
		widget.NewSelect(serialNumbers, func(serialNumber string) {
			onSerialNumberSelected(a, w, serialNumber)
		}),
		widget.NewButton("Quit", a.Quit),
	))

	w.ShowAndRun()
}

func onSerialNumberSelected(a fyne.App, w fyne.Window, serialNumber string) {
	log.Printf("Serial Number: %s", serialNumber)

	apps, err := ifuse.ListApps(serialNumber)
	if err != nil {
		log.Fatalf("Failed to list apps for %s: %v", serialNumber, err)
	}

	var appNames []string
	appsMap := map[string]ifuse.App{}
	for _, a := range apps {
		appNames = append(appNames, a.Name)
		appsMap[a.Name] = a
	}
	sort.Strings(appNames)

	w.SetContent(widget.NewVBox(
		widget.NewLabel("Select application:"),
		widget.NewSelect(appNames, func(appName string) {
			app := appsMap[appName]
			onAppSelected(a, w, serialNumber, app)
		}),
		widget.NewButton("Quit", a.Quit),
	))
}

func onAppSelected(a fyne.App, w fyne.Window, serialNumber string, app ifuse.App) {
	mountPath, err := ifuse.MountAppDocuments(serialNumber, app)
	if err != nil {
		log.Fatalf("Failed to mount app documents (%s): %v", app.ID, err)
	}

	if err := open.Start(mountPath); err != nil {
		log.Fatalf("Failed to open mount path %q: %v", mountPath, err)
	}

	a.Quit()
}

func getSerialNumbers() ([]string, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return fmt.Sprintf("%s:%s", desc.Vendor, desc.Product) == iDeviceVendorProduct
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, d := range devices {
			d.Close()
		}
	}()

	formatSerialNumber := func(serialNumber string) string {
		serialNumber = strings.Trim(serialNumber, "\x00")
		return fmt.Sprintf("%s-%s", serialNumber[:8], serialNumber[8:])
	}

	var serialNumbers []string
	for _, d := range devices {
		serialNumber, err := d.SerialNumber()
		if err != nil {
			return nil, err
		}
		serialNumbers = append(serialNumbers, formatSerialNumber(serialNumber))
	}
	return serialNumbers, nil
}

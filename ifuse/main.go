package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/gousb"
	"github.com/ramin0/live/ifuse/ifuse"
)

const iDeviceVendorProduct = "05ac:12a8"

func main() {
	serialNumbers, err := getSerialNumbers()
	if err != nil {
		log.Fatalf("Failed to get serial numbers: %v", err)
	}

	if len(serialNumbers) == 0 {
		log.Fatal("No iDevices found")
	}

	serialNumber := serialNumbers[0]
	apps, err := ifuse.ListApps(serialNumber)
	if err != nil {
		log.Fatalf("Failed to list apps for %s: %v", serialNumber, err)
	}

	app := apps[0]
	mountPath, err := ifuse.MountAppDocuments(serialNumber, app)
	if err != nil {
		log.Fatalf("Failed to mount app documents (%s): %v", app.ID, err)
	}
	fmt.Println(mountPath)
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

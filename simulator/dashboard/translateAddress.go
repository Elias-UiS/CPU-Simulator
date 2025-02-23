package dashboard

// import (
// 	"CPU-Simulator/simulator/pkg/cpu"
// 	"CPU-Simulator/simulator/pkg/temp"
// 	"fmt"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/translator"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func setupCalculatorTab() fyne.CanvasObject {
	// Create input fields and button for the top row.

	input1 := widget.NewEntry()
	input1.SetPlaceHolder("VPN/PFN")
	input1.Resize(fyne.NewSize(150, 30))
	input1.OnChanged = func(text string) {
		// Filter out non-digit characters
		filtered := ""
		for _, r := range text {
			if r >= '0' && r <= '9' {
				filtered += string(r)
			}
		}
		if filtered != text {
			input1.SetText(filtered)
		}
	}

	input2 := widget.NewEntry()
	input2.SetPlaceHolder("Offset")
	input2.Resize(fyne.NewSize(150, 30))
	input2.OnChanged = func(text string) {
		// Filter out non-digit characters
		filtered := ""
		for _, r := range text {
			if r >= '0' && r <= '9' {
				filtered += string(r)
			}
		}
		if filtered != text {
			input2.SetText(filtered)
		}
	}

	// Create result labels
	intResult := widget.NewEntry()
	intResult.SetPlaceHolder("Don't edit")
	bitsResult := widget.NewEntry()
	bitsResult.SetPlaceHolder("Don't edit")

	hexResult := widget.NewEntry()
	hexResult.SetPlaceHolder("Don't edit")

	intParts := widget.NewLabel("Int Parts: [vpn] [offset]") // New field for the intermediate step

	intPartsVpnResult := widget.NewEntry()
	intPartsVpnResult.SetPlaceHolder("Don't edit")
	intPartsOffsetResult := widget.NewEntry()
	intPartsOffsetResult.SetPlaceHolder("Don't edit")

	calcButton := widget.NewButton("Calculate", func() {
		input1Int, err1 := strconv.Atoi(input1.Text)
		input2Int, err2 := strconv.Atoi(input2.Text)
		if err1 != nil || err2 != nil {
			logger.Log.Println("Invalid input: please enter valid numbers")
			return
		}
		address := translator.TranslateVPNandOffsetToAddress(input1Int, input2Int)

		// Update the label with the calculated address

		intResult.SetText(fmt.Sprintf("%d", address))
		binaryAddress := fmt.Sprintf("%032b", address)                   // Get 32-bit binary representation
		formattedBinary := binaryAddress[:16] + " " + binaryAddress[16:] // Insert space in the middle
		bitsResult.SetText(formattedBinary)
		hexResult.SetText(fmt.Sprintf("0x%X", address)) // Hex format
		vpnBaseAddress := input1Int << 16               // Left shift VPN by 16 bits
		intPartsVpnResult.SetText(fmt.Sprintf("%d", vpnBaseAddress))
		intPartsOffsetResult.SetText(fmt.Sprintf("%d", input2Int))
	})
	calcButton.Resize(fyne.NewSize(150, 30))

	// Top row container: Two inputs and button on the right.
	topRow := container.NewGridWithColumns(3,
		input1,
		input2,
		calcButton,
	)

	intPartsResultContainer := container.NewGridWithColumns(2,
		intPartsVpnResult,
		intPartsOffsetResult,
	)

	// Rows for displaying values

	intRow := widget.NewLabel("Int: ")
	bitsRow := widget.NewLabel("Bits: ")
	hexRow := widget.NewLabel("Hexadecimal: ")
	bumpRow := widget.NewLabel("")

	resultContainerText := container.NewGridWithRows(4,
		intRow,
		intParts,
		bitsRow,
		hexRow,
	)

	resultContainerBox := container.NewGridWithRows(4,
		intResult,
		intPartsResultContainer,
		bitsResult,
		hexResult,
	)

	resultContainer := container.NewGridWithColumns(2,
		resultContainerText,
		resultContainerBox,
	)

	// Create the first vertical container.
	topContainer := container.NewGridWithRows(2,
		topRow,
		resultContainer,
	)

	vpnOffsetResult := widget.NewLabel("VPN: 0 , Offset: 0")

	inputAddress := widget.NewEntry()
	inputAddress.SetPlaceHolder("Address(int)")
	inputAddress.Resize(fyne.NewSize(300, 35))
	inputAddress.OnChanged = func(text string) {
		// Filter out non-digit characters
		filtered := ""
		for _, r := range text {
			if r >= '0' && r <= '9' {
				filtered += string(r)
			}
		}
		if filtered != text {
			inputAddress.SetText(filtered)
		}
	}

	calcButton2 := widget.NewButton("Calculate", func() {
		inputAddressInt, err := strconv.Atoi(inputAddress.Text)
		if err != nil {
			logger.Log.Println("Invalid input: please enter valid numbers")
			vpnOffsetResult.SetText(fmt.Sprintf("VPN: Error , Offset: Error"))
			return
		}
		vpn, offset := translator.TranslateAddressToVPNandOffset(inputAddressInt)

		// Update the label with the calculated address
		vpnOffsetResult.SetText(fmt.Sprintf("VPN: %d , Offset: %d", vpn, offset))

	})

	// Top row container: Two inputs and button on the right.
	addressContainer := container.NewGridWithColumns(2,
		inputAddress,
		calcButton2,
	)

	bottomContainer := container.NewVBox(
		bumpRow,
		addressContainer,
		vpnOffsetResult,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
		bumpRow,
	)

	// Combine both vertical containers in a horizontal container.
	calculatorTab := container.NewVBox(
		topContainer,
		bottomContainer,
	)
	calculatorTab.Resize(fyne.NewSize(600, 30))

	return calculatorTab
}

package app

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

func infoSpinner(s ...any) *pterm.SpinnerPrinter {
	sp := pterm.DefaultSpinner.WithDelay(1 * time.Second)
	spinner, _ := sp.Start(s...)

	return spinner
}

func errorSpinner(spinner *pterm.SpinnerPrinter, err error) error {
	spinner.Fail(err)

	return err
}

func printWarn(fmt string, s ...any) {
	style := pterm.Bold.ToStyle().Add(*pterm.FgYellow.ToStyle())
	pterm.Warning.WithMessageStyle(&style).Printfln(fmt, s...)
}

func printInfo(fmt string, s ...any) {
	pterm.Info.WithMessageStyle(pterm.NewStyle(pterm.FgBlue, pterm.Bold)).Printfln(fmt, s...)
}

func printNote(fmt string, s ...any) {
	prefix := pterm.Info.Prefix

	pterm.Info.Prefix = pterm.Prefix{Text: "NOTE", Style: pterm.NewStyle(pterm.BgGray, pterm.Bold)}
	pterm.Info.WithMessageStyle(pterm.NewStyle(pterm.FgGray, pterm.Bold)).Printfln(fmt, s...)
	pterm.Info.Prefix = prefix
}

func newline() {
	fmt.Println()
}

func updateSpinner(spinner *pterm.SpinnerPrinter, fmt string, s ...any) {
	style := pterm.FgDefault.Sprintf(fmt, s...)
	spinner.UpdateText(style)
}

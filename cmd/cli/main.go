package main

import (
	"fmt"
	"os"
	"time"

	"github.com/f00b455/golang-template/pkg/core"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
)

var (
	name string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hello-cli",
	Short: "A colorful Hello-World CLI application",
	Long:  `A colorful Hello-World CLI application built with Go and Cobra.`,
	Run:   runHelloCommand,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&name, "name", "World", "Name to greet")
}

func runHelloCommand(cmd *cobra.Command, args []string) {
	// Welcome message
	magenta := color.New(color.FgMagenta).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s\n\n", magenta("🎉 Welcome to Hello CLI!"))

	// Run spinner
	if err := runSpinner(); err != nil {
		fmt.Printf("%s\n", red(fmt.Sprintf("❌ Error: %v", err)))
		os.Exit(1)
	}

	// Run progress bar
	if err := runProgressBar(); err != nil {
		fmt.Printf("%s\n", red(fmt.Sprintf("❌ Error: %v", err)))
		os.Exit(1)
	}

	// Display greeting
	fmt.Println()
	displayGreeting(name)

	fmt.Printf("%s\n", green("✅ All done! Have a great day!"))
}

func runSpinner() error {
	// Use shorter delay for testing
	spinnerDelay := 5 * time.Second
	if os.Getenv("GO_ENV") == "test" {
		spinnerDelay = 100 * time.Millisecond
	}

	cfg := yacspin.Config{
		Frequency: 100 * time.Millisecond,
		CharSet:   yacspin.CharSets[14], // Dots spinner
		Message:   " Preparing something awesome...",
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return err
	}

	if err := spinner.Start(); err != nil {
		return err
	}

	time.Sleep(spinnerDelay)

	spinner.Message(" Ready!")
	if err := spinner.Stop(); err != nil {
		return err
	}

	fmt.Println("Spinner completed: Ready!")
	return nil
}

func runProgressBar() error {
	// Use shorter delay for testing
	progressDelay := 30 * time.Millisecond
	if os.Getenv("GO_ENV") == "test" {
		progressDelay = 1 * time.Millisecond
	}

	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription("Progress"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "|",
			BarEnd:        "|",
		}))

	for i := 0; i <= 100; i++ {
		_ = bar.Add(1) // Ignore error from progress bar
		time.Sleep(progressDelay)
	}

	fmt.Println("\nProgress completed: 100%")
	return nil
}

func displayGreeting(name string) {
	config := core.FooConfig{
		Prefix: "✨",
		Suffix: "✨",
	}

	greetingMessage := core.FooGreet(config, name)

	// Create colorful text (simplified gradient effect)
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

	// Simple rainbow effect by cycling through colors
	coloredMessage := ""
	colors := []func(...interface{}) string{cyan, yellow, magenta}
	for i, char := range greetingMessage {
		colorFunc := colors[i%len(colors)]
		coloredMessage += colorFunc(string(char))
	}

	// Create boxed output (simplified)
	boxWidth := len(greetingMessage) + 4

	fmt.Printf("    ╔%s╗\n", fmt.Sprintf("%*s", boxWidth-2, ""))
	fmt.Printf("    ║%s║\n", fmt.Sprintf("%*s", boxWidth-2, "Hello CLI"))
	fmt.Printf("    ╠%s╣\n", fmt.Sprintf("%*s", boxWidth-2, ""))
	fmt.Printf("    ║ %s ║\n", coloredMessage)
	fmt.Printf("    ║%s║\n", fmt.Sprintf("%*s", boxWidth-2, ""))
	fmt.Printf("    ╚%s╝\n", fmt.Sprintf("%*s", boxWidth-2, ""))
}

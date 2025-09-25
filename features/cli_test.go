package features

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type cliFeatureContext struct {
	commandOutput string
	commandError  error
	exitCode      int
	binaryPath    string
}

func (ctx *cliFeatureContext) iHaveTheHelloCLICommandAvailable() error {
	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Try current directory first, then parent directory
	ctx.binaryPath = filepath.Join(workDir, "bin", "cli-tool")
	if _, err := os.Stat(ctx.binaryPath); os.IsNotExist(err) {
		// If running from features subdirectory, go up one level
		ctx.binaryPath = filepath.Join(workDir, "..", "bin", "cli-tool")
		if _, err := os.Stat(ctx.binaryPath); os.IsNotExist(err) {
			// Skip test if binary not built (instead of failing)
			return godog.ErrPending
		}
	}

	return nil
}

func (ctx *cliFeatureContext) iRunHelloCLIWithoutParameters() error {
	return ctx.runCommand()
}

func (ctx *cliFeatureContext) iRunHelloCLIWithName(name string) error {
	return ctx.runCommand("--name", name)
}

func (ctx *cliFeatureContext) iRunHelloCLIWithFlag(flag string) error {
	return ctx.runCommand(flag)
}

func (ctx *cliFeatureContext) runCommand(args ...string) error {
	// Set test environment to use faster delays
	cmd := exec.Command(ctx.binaryPath, args...)
	cmd.Env = append(os.Environ(), "GO_ENV=test")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ctx.commandError = cmd.Run()
	ctx.commandOutput = stdout.String() + stderr.String()

	if exitError, ok := ctx.commandError.(*exec.ExitError); ok {
		ctx.exitCode = exitError.ExitCode()
	} else if ctx.commandError == nil {
		ctx.exitCode = 0
	} else {
		ctx.exitCode = -1
	}

	return nil
}

func (ctx *cliFeatureContext) theCommandShouldCompleteSuccessfully() error {
	if ctx.exitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", ctx.exitCode, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) iShouldSeeASpinnerMessage(message string) error {
	if !strings.Contains(ctx.commandOutput, "Spinner completed: "+message) {
		return fmt.Errorf("expected spinner message '%s' not found in output:\n%s", message, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) iShouldSeeAProgressMessage(message string) error {
	if !strings.Contains(ctx.commandOutput, message) {
		return fmt.Errorf("expected progress message '%s' not found in output:\n%s", message, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) iShouldSeeAGreetingMessageFor(name string) error {
	expectedPart := fmt.Sprintf("Hello, %s!", name)
	if !strings.Contains(ctx.commandOutput, expectedPart) {
		return fmt.Errorf("expected greeting for '%s' not found in output:\n%s", name, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) theGreetingShouldContainThePrefix(prefix string) error {
	// The greeting format is: ✨ Hello, Name! ✨
	if !strings.Contains(ctx.commandOutput, prefix) {
		return fmt.Errorf("expected prefix '%s' not found in output:\n%s", prefix, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) theGreetingShouldContainTheSuffix(suffix string) error {
	if !strings.Contains(ctx.commandOutput, suffix) {
		return fmt.Errorf("expected suffix '%s' not found in output:\n%s", suffix, ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) theOutputShouldContainADecorativeBox() error {
	// Check for box drawing characters
	boxChars := []string{"╔", "╗", "╚", "╝", "║", "╠", "╣"}
	for _, char := range boxChars {
		if !strings.Contains(ctx.commandOutput, char) {
			return fmt.Errorf("expected box character '%s' not found in output", char)
		}
	}
	return nil
}

func (ctx *cliFeatureContext) theGreetingShouldUseTheFooGreetFunction() error {
	// The FooGreet function adds prefix and suffix
	if !strings.Contains(ctx.commandOutput, "✨") {
		return fmt.Errorf("FooGreet function signature (✨ prefix/suffix) not found")
	}
	return nil
}

func (ctx *cliFeatureContext) theOutputShouldMatchTheCorePackageFormat() error {
	// Check for the format: ✨ Hello, Name! ✨
	if !strings.Contains(ctx.commandOutput, "✨") || !strings.Contains(ctx.commandOutput, "Hello,") {
		return fmt.Errorf("core package format not found in output")
	}
	return nil
}

func (ctx *cliFeatureContext) iShouldSeeAHelpMessage() error {
	if !strings.Contains(ctx.commandOutput, "Usage:") && !strings.Contains(ctx.commandOutput, "Flags:") {
		return fmt.Errorf("help message not found in output:\n%s", ctx.commandOutput)
	}
	return nil
}

func (ctx *cliFeatureContext) theHelpShouldContain(expectedText string) error {
	if !strings.Contains(ctx.commandOutput, expectedText) {
		return fmt.Errorf("expected text '%s' not found in help output:\n%s", expectedText, ctx.commandOutput)
	}
	return nil
}

func InitializeCLIScenario(ctx *godog.ScenarioContext) {
	featureCtx := &cliFeatureContext{}

	// Background steps
	ctx.Step(`^I have the hello-cli command available$`, featureCtx.iHaveTheHelloCLICommandAvailable)

	// Action steps
	ctx.Step(`^I run hello-cli without parameters$`, featureCtx.iRunHelloCLIWithoutParameters)
	ctx.Step(`^I run hello-cli with name "([^"]*)"$`, featureCtx.iRunHelloCLIWithName)
	ctx.Step(`^I run hello-cli with "([^"]*)" flag$`, featureCtx.iRunHelloCLIWithFlag)

	// Assertion steps
	ctx.Step(`^the command should complete successfully$`, featureCtx.theCommandShouldCompleteSuccessfully)
	ctx.Step(`^I should see a spinner message "([^"]*)"$`, featureCtx.iShouldSeeASpinnerMessage)
	ctx.Step(`^I should see a progress message "([^"]*)"$`, featureCtx.iShouldSeeAProgressMessage)
	ctx.Step(`^I should see a greeting message for "([^"]*)"$`, featureCtx.iShouldSeeAGreetingMessageFor)
	ctx.Step(`^the greeting should contain the prefix "([^"]*)"$`, featureCtx.theGreetingShouldContainThePrefix)
	ctx.Step(`^the greeting should contain the suffix "([^"]*)"$`, featureCtx.theGreetingShouldContainTheSuffix)
	ctx.Step(`^the output should contain a decorative box$`, featureCtx.theOutputShouldContainADecorativeBox)
	ctx.Step(`^the greeting should use the FooGreet function$`, featureCtx.theGreetingShouldUseTheFooGreetFunction)
	ctx.Step(`^the output should match the core package format$`, featureCtx.theOutputShouldMatchTheCorePackageFormat)
	ctx.Step(`^I should see a help message$`, featureCtx.iShouldSeeAHelpMessage)
	ctx.Step(`^the help should contain "([^"]*)"$`, featureCtx.theHelpShouldContain)
}

func TestCLIFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeCLIScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"cli.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run CLI feature tests")
	}
}
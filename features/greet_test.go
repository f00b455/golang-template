package features

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/f00b455/golang-template/pkg/shared"
)

type greetFeatureContext struct {
	name   string
	result string
}

func (ctx *greetFeatureContext) iHaveTheName(name string) error {
	ctx.name = name
	return nil
}

func (ctx *greetFeatureContext) iCallTheGreetFunction() error {
	ctx.result = shared.Greet(ctx.name)
	return nil
}

func (ctx *greetFeatureContext) iShouldReceive(expectedMessage string) error {
	if ctx.result != expectedMessage {
		return fmt.Errorf("expected '%s', but got '%s'", expectedMessage, ctx.result)
	}
	return nil
}

func (ctx *greetFeatureContext) iAmUsingTheSharedGreetFunction() error {
	// This is just a background step - no action needed
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	featureCtx := &greetFeatureContext{}

	ctx.Step(`^I am using the shared greet function$`, featureCtx.iAmUsingTheSharedGreetFunction)
	ctx.Step(`^I have the name "([^"]*)"$`, featureCtx.iHaveTheName)
	ctx.Step(`^I call the greet function$`, featureCtx.iCallTheGreetFunction)
	ctx.Step(`^I should receive "([^"]*)"$`, featureCtx.iShouldReceive)
}

func TestSharedFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"greet.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run shared feature tests")
	}
}

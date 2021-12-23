//go:build provider
// +build provider

package main

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"testing"
)
type Settings struct {
	Host            string
	ProviderName    string
	BrokerBaseURL   string
	BrokerUsername  string // Basic authentication
	BrokerPassword  string // Basic authentication
	ConsumerName    string
	ConsumerVersion string // a git sha, semantic version number
	ProviderVersion string
}

func TestGuestLedgerBookProvider(t *testing.T) {
	pact := dsl.Pact{
		Consumer: "GuestLadgerClient",
		Provider: "GuestLadgerApi",
		DisableToolValidityCheck: true,
	}
	port,_ := utils.GetFreePort()
	go StartServer(port)
	_,err := pact.VerifyProvider(t,types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://localhost:%d", port),
		PactURLs: []string{"https://eneskzlcn.pactflow.io/pacts/provider/GuestLadgerApi/consumer/GuestLadgerClient/version/2e88856dd7a768eaa2cfbeb603cd37caf9427648"},
		PublishVerificationResults: true,
		ProviderVersion:            "1.0.0",
		BrokerToken:                "L0IzB6WxiCRX7sEdAQoWlQ",
		Tags: []string{"main"},
		StateHandlers: map[string]types.StateHandler{
			"get guest ladger by email successfully": func() error {
				return nil
			},
			"get all guest ladger book successfully": func() error {
				return nil
			},
			"new guest ladger posted successfuly": func() error {
				return nil
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

package common

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/config"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/helpers"

	cliConfig "github.com/oasisprotocol/cli/config"
)

const DefaultMarker = " (*)"

var (
	selectedNetwork  string
	selectedParaTime string
	selectedAccount  string

	noParaTime bool
)

var (
	// SelectorFlags contains the common selector flags for network/paratime/wallet.
	SelectorFlags *flag.FlagSet
	// SelectorNPFlags contains the common selector flags for network/paratime.
	SelectorNPFlags *flag.FlagSet
	// SelectorNFlags contains the common selector flags for network.
	SelectorNFlags *flag.FlagSet
	// SelectorNAFlags contains the common selector flags for network/account.
	SelectorNAFlags *flag.FlagSet
)

// NPASelection contains the network/paratime/account selection.
type NPASelection struct {
	NetworkName string
	Network     *config.Network

	ParaTimeName string
	ParaTime     *config.ParaTime

	AccountName string
	Account     *cliConfig.Account
}

// GetNPASelection returns the user-selected network/paratime/account combination.
func GetNPASelection(cfg *cliConfig.Config) *NPASelection {
	var s NPASelection
	s.NetworkName = cfg.Networks.Default
	if selectedNetwork != "" {
		s.NetworkName = selectedNetwork
	}
	if s.NetworkName == "" {
		cobra.CheckErr(fmt.Errorf("no networks configured"))
	}
	s.Network = cfg.Networks.All[s.NetworkName]
	if s.Network == nil {
		cobra.CheckErr(fmt.Errorf("network '%s' does not exist", s.NetworkName))
	}

	if !noParaTime {
		s.ParaTimeName = s.Network.ParaTimes.Default
		if selectedParaTime != "" {
			s.ParaTimeName = selectedParaTime
		}
		if s.ParaTimeName != "" {
			s.ParaTime = s.Network.ParaTimes.All[s.ParaTimeName]
			if s.ParaTime == nil {
				cobra.CheckErr(fmt.Errorf("paratime '%s' does not exist", s.ParaTimeName))
			}
		}
	}

	s.AccountName = cfg.Wallet.Default
	if selectedAccount != "" {
		s.AccountName = selectedAccount
	}
	if s.AccountName != "" {
		if testName := helpers.ParseTestAccountAddress(s.AccountName); testName != "" {
			testAcc, err := LoadTestAccountConfig(testName)
			cobra.CheckErr(err)
			s.Account = testAcc
		} else {
			s.Account = cfg.Wallet.All[s.AccountName]
			if s.Account == nil {
				cobra.CheckErr(fmt.Errorf("account '%s' does not exist in the wallet", s.AccountName))
			}
		}
	}

	return &s
}

func init() {
	SelectorFlags = flag.NewFlagSet("", flag.ContinueOnError)
	SelectorFlags.StringVar(&selectedNetwork, "network", "", "explicitly set network to use")
	SelectorFlags.StringVar(&selectedParaTime, "paratime", "", "explicitly set paratime to use")
	SelectorFlags.BoolVar(&noParaTime, "no-paratime", false, "explicitly set that no paratime should be used")
	SelectorFlags.StringVar(&selectedAccount, "account", "", "explicitly set account to use")

	SelectorNPFlags = flag.NewFlagSet("", flag.ContinueOnError)
	SelectorNPFlags.StringVar(&selectedNetwork, "network", "", "explicitly set network to use")
	SelectorNPFlags.StringVar(&selectedParaTime, "paratime", "", "explicitly set paratime to use")
	SelectorNPFlags.BoolVar(&noParaTime, "no-paratime", false, "explicitly set that no paratime should be used")

	SelectorNAFlags = flag.NewFlagSet("", flag.ContinueOnError)
	SelectorNAFlags.StringVar(&selectedNetwork, "network", "", "explicitly set network to use")
	SelectorNAFlags.StringVar(&selectedAccount, "account", "", "explicitly set account to use")

	SelectorNFlags = flag.NewFlagSet("", flag.ContinueOnError)
	SelectorNFlags.StringVar(&selectedNetwork, "network", "", "explicitly set network to use")

	// Backward compatibility.
	SelectorFlags.StringVar(&selectedAccount, "wallet", "", "explicitly set account to use. OBSOLETE, USE --account INSTEAD!")
	err := SelectorFlags.MarkHidden("wallet")
	cobra.CheckErr(err)
}

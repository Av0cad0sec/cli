package account

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/oasisprotocol/oasis-core/go/common/entity"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
	"github.com/oasisprotocol/oasis-sdk/client-sdk/go/connection"

	"github.com/oasisprotocol/cli/cmd/common"
	cliConfig "github.com/oasisprotocol/cli/config"
)

var (
	entityCmd = &cobra.Command{
		Use:   "entity",
		Short: "Entity management in the network's registry",
	}

	entityRegisterCmd = &cobra.Command{
		Use:   "register <entity.json>",
		Short: "Register a new entity or update an existing one",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cliConfig.Global()
			npa := common.GetNPASelection(cfg)
			txCfg := common.GetTransactionConfig()
			filename := args[0]

			if npa.Account == nil {
				cobra.CheckErr("no accounts configured in your wallet")
			}

			// When not in offline mode, connect to the given network endpoint.
			ctx := context.Background()
			var conn connection.Connection
			if !txCfg.Offline {
				var err error
				conn, err = connection.Connect(ctx, npa.Network)
				cobra.CheckErr(err)
			}

			// Load entity descriptor.
			rawDescriptor, err := os.ReadFile(filename)
			cobra.CheckErr(err)

			// Parse entity descriptor.
			var descriptor entity.Entity
			if err = json.Unmarshal(rawDescriptor, &descriptor); err != nil {
				cobra.CheckErr(fmt.Errorf("malformed entity descriptor: %w", err))
			}

			// Load the account and ensure it corresponds to the entity.
			acc := common.LoadAccount(cfg, npa.AccountName)
			signer := acc.ConsensusSigner()
			if signer == nil {
				cobra.CheckErr(fmt.Errorf("account '%s' does not support signing consensus transactions", npa.AccountName))
			}
			if !signer.Public().Equal(descriptor.ID) {
				cobra.CheckErr(fmt.Errorf("entity ID '%s' does not correspond to selected account '%s' (%s)",
					descriptor.ID, npa.AccountName, signer.Public()))
			}

			// Sign entity descriptor.
			fmt.Println("Signing the entity descriptor...")
			fmt.Println("(In case you are using a hardware-based signer you may need to confirm on device.)")
			sigDescriptor, err := entity.SignEntity(signer, registry.RegisterEntitySignatureContext, &descriptor)
			cobra.CheckErr(err)

			// Prepare transaction.
			tx := registry.NewRegisterEntityTx(0, nil, sigDescriptor)

			sigTx, err := common.SignConsensusTransaction(ctx, npa, acc, conn, tx)
			cobra.CheckErr(err)

			common.BroadcastOrExportTransaction(ctx, npa.ParaTime, conn, sigTx, nil, nil)
		},
	}

	entityDeregisterCmd = &cobra.Command{
		Use:   "deregister",
		Short: "Deregister an existing entity",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cliConfig.Global()
			npa := common.GetNPASelection(cfg)
			txCfg := common.GetTransactionConfig()

			if npa.Account == nil {
				cobra.CheckErr("no accounts configured in your wallet")
			}

			// When not in offline mode, connect to the given network endpoint.
			ctx := context.Background()
			var conn connection.Connection
			if !txCfg.Offline {
				var err error
				conn, err = connection.Connect(ctx, npa.Network)
				cobra.CheckErr(err)
			}

			// Prepare transaction.
			tx := registry.NewDeregisterEntityTx(0, nil)

			acc := common.LoadAccount(cfg, npa.AccountName)
			sigTx, err := common.SignConsensusTransaction(ctx, npa, acc, conn, tx)
			cobra.CheckErr(err)

			common.BroadcastOrExportTransaction(ctx, npa.ParaTime, conn, sigTx, nil, nil)
		},
	}
)

func init() {
	entityRegisterCmd.Flags().AddFlagSet(common.SelectorNAFlags)
	entityRegisterCmd.Flags().AddFlagSet(common.TransactionFlags)

	entityDeregisterCmd.Flags().AddFlagSet(common.SelectorNAFlags)
	entityDeregisterCmd.Flags().AddFlagSet(common.TransactionFlags)

	entityCmd.AddCommand(entityRegisterCmd)
	entityCmd.AddCommand(entityDeregisterCmd)
}

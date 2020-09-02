package handler

import (
	"encoding/hex"

	ethkeymanager "github.com/bloxapp/eth-key-manager"
	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/public-key/flag"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"
)

// Account generates a new wallet account at specific index and prints the account.
func (h *PublicKey) Generate(cmd *cobra.Command, args []string) error {
	err := types.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	// Get index flag.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the index flag value")
	}
	// Validate
	if indexFlagValue < 0 {
		return errors.New("provided index is negative")
	}

	// Get seed flag.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode seed")
	}

	store := in_memory.NewInMemStore()
	options := &ethkeymanager.KeyVaultOptions{}
	options.SetStorage(store)

	_, err = ethkeymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault.")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	account, err := wallet.CreateValidatorAccount(seedBytes, &indexFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to create validator account")
	}

	publicKey := map[string]string{
		"validationPubKey": hex.EncodeToString(account.ValidatorPublicKey().Marshal()),
		"withdrawalPubKey": hex.EncodeToString(account.WithdrawalPublicKey().Marshal()),
	}

	err = h.printer.JSON(publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to print public-key JSON")
	}
	return nil
}

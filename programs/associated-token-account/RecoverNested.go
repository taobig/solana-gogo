// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package associatedtokenaccount

import (
	"errors"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/text/format"
	"github.com/gagliardetto/treeout"
)

type RecoverNested struct {
	Wallet     solana.PublicKey `bin:"-" borsh_skip:"true"`
	OwnerMint  solana.PublicKey `bin:"-" borsh_skip:"true"`
	NestedMint solana.PublicKey `bin:"-" borsh_skip:"true"`

	// [0] = [WRITE] NestedAssociatedTokenAccount
	// ··········· Nested associated token account, must be owned by [3]
	//
	// [1] = [] NestedMint
	// ··········· Token mint for the nested associated token account
	//
	// [2] = [WRITE] WalletAssociatedTokenAccount
	// ··········· Wallet's associated token account
	//
	// [3] = [] OwnerAssociatedTokenAccount
	// ··········· Owner associated token account address, must be owned by [5]
	//
	// [4] = [] OwnerMint
	// ··········· Token mint for the owner associated token account
	//
	// [5] = [WRITE, SIGNER] Wallet
	// ··········· Wallet address for the owner associated token account
	//
	// [6] = [] TokenProgram
	// ··········· SPL Token program ID
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewRecoverNestedInstructionBuilder creates a new `RecoverNested` instruction builder.
func NewRecoverNestedInstructionBuilder() *RecoverNested {
	nd := &RecoverNested{}
	return nd
}

func (inst *RecoverNested) SetWallet(wallet solana.PublicKey) *RecoverNested {
	inst.Wallet = wallet
	return inst
}

func (inst *RecoverNested) SetOwnerMint(ownerMint solana.PublicKey) *RecoverNested {
	inst.OwnerMint = ownerMint
	return inst
}

func (inst *RecoverNested) SetNestedMint(nestedMint solana.PublicKey) *RecoverNested {
	inst.NestedMint = nestedMint
	return inst
}

func (inst RecoverNested) Build() *Instruction {

	// Find associated token addresses
	ownerAssociatedTokenAddress, _, _ := solana.FindAssociatedTokenAddress(
		inst.Wallet,
		inst.OwnerMint,
	)
	dstAssociatedTokenAddress, _, _ := solana.FindAssociatedTokenAddress(
		inst.Wallet,
		inst.NestedMint,
	)
	nestedAssociatedTokenAddress, _, _ := solana.FindAssociatedTokenAddress(
		ownerAssociatedTokenAddress, // ATA is wrongly used as a wallet_address
		inst.NestedMint,
	)

	keys := []*solana.AccountMeta{
		{
			PublicKey:  nestedAssociatedTokenAddress,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  inst.NestedMint,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  dstAssociatedTokenAddress,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  ownerAssociatedTokenAddress,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inst.OwnerMint,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inst.Wallet,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			PublicKey:  solana.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	inst.AccountMetaSlice = keys

	return &Instruction{BaseVariant: bin.BaseVariant{
		Impl:   inst,
		TypeID: bin.TypeIDFromUint8(Instruction_RecoverNested),
	}}
}

// ValidateAndBuild validates the instruction accounts.
// If there is a validation error, return the error.
// Otherwise, build and return the instruction.
func (inst RecoverNested) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RecoverNested) Validate() error {
	if inst.Wallet.IsZero() {
		return errors.New("Wallet not set")
	}
	if inst.OwnerMint.IsZero() {
		return errors.New("OwnerMint not set")
	}
	if inst.NestedMint.IsZero() {
		return errors.New("NestedMint not set")
	}
	ownerAssociatedTokenAddress, _, err := solana.FindAssociatedTokenAddress(
		inst.Wallet,
		inst.OwnerMint,
	)
	if err != nil {
		return fmt.Errorf("error while FindAssociatedTokenAddress: %w", err)
	}
	_, _, err = solana.FindAssociatedTokenAddress(
		inst.Wallet,
		inst.NestedMint,
	)
	if err != nil {
		return fmt.Errorf("error while FindAssociatedTokenAddress: %w", err)
	}
	_, _, err = solana.FindAssociatedTokenAddress(
		ownerAssociatedTokenAddress,
		inst.NestedMint,
	)
	if err != nil {
		return fmt.Errorf("error while FindAssociatedTokenAddress: %w", err)
	}
	return nil
}

func (inst *RecoverNested) EncodeToTree(parent treeout.Branches) {
	parent.Child(format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch treeout.Branches) {
			programBranch.Child(format.Instruction("RecoverNested")).
				//
				ParentFunc(func(instructionBranch treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=7]").ParentFunc(func(accountsBranch treeout.Branches) {
						accountsBranch.Child(format.Meta("nestedAssociatedTokenAccount", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(format.Meta("             nestedTokenMint", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(format.Meta("   dstAssociatedTokenAccount", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(format.Meta(" ownerAssociatedTokenAccount", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(format.Meta("              ownerTokenMint", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(format.Meta("                      wallet", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(format.Meta("                tokenProgram", inst.AccountMetaSlice.Get(6)))
					})
				})
		})
}

func (inst RecoverNested) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	return encoder.WriteBytes([]byte{}, false)
}

func (inst *RecoverNested) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	return nil
}

func NewRecoverNestedInstruction(
	Wallet solana.PublicKey,
	OwnerMint solana.PublicKey,
	NestedMint solana.PublicKey,
) *RecoverNested {
	return NewRecoverNestedInstructionBuilder().
		SetWallet(Wallet).
		SetOwnerMint(OwnerMint).
		SetNestedMint(NestedMint)
}

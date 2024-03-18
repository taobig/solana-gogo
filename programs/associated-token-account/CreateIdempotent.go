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
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/text/format"
	"github.com/gagliardetto/treeout"
)

type CreateIdempotent struct {
	Create

	associatedToken solana.PublicKey
	programID       solana.PublicKey
}

// NewCreateIdempotentInstructionBuilder creates a new `CreateIdempotent` instruction builder.
func NewCreateIdempotentInstructionBuilder() *CreateIdempotent {
	nd := &CreateIdempotent{}
	return nd
}

func (inst *CreateIdempotent) SetPayer(payer solana.PublicKey) *CreateIdempotent {
	inst.Payer = payer
	return inst
}

func (inst *CreateIdempotent) SetWallet(wallet solana.PublicKey) *CreateIdempotent {
	inst.Wallet = wallet
	return inst
}

func (inst *CreateIdempotent) SetMint(mint solana.PublicKey) *CreateIdempotent {
	inst.Mint = mint
	return inst
}

func (inst *CreateIdempotent) SetAssociatedToken(associatedToken solana.PublicKey) *CreateIdempotent {
	inst.associatedToken = associatedToken
	return inst
}

func (inst *CreateIdempotent) SetProgramID(tokenProgramID solana.PublicKey) *CreateIdempotent {
	inst.programID = tokenProgramID
	return inst
}

func (inst CreateIdempotent) Build() *Instruction {

	//// Find the associatedTokenAddress;
	//associatedTokenAddress, _, _ := solana.FindAssociatedTokenAddress(
	//	inst.Wallet,
	//	inst.Mint,
	//)

	keys := []*solana.AccountMeta{
		{
			PublicKey:  inst.Payer,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			//PublicKey:  associatedTokenAddress,
			PublicKey:  inst.associatedToken,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  inst.Wallet,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inst.Mint,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  solana.SystemProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			//PublicKey:  solana.TokenProgramID,
			PublicKey:  inst.programID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  solana.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	inst.AccountMetaSlice = keys

	return &Instruction{BaseVariant: bin.BaseVariant{
		Impl:   inst,
		TypeID: bin.TypeIDFromUint8(Instruction_CreateIdempotent),
	}}
}

func (inst *CreateIdempotent) EncodeToTree(parent treeout.Branches) {
	parent.Child(format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch treeout.Branches) {
			programBranch.Child(format.Instruction("CreateIdempotent")).
				//
				ParentFunc(func(instructionBranch treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=7]").ParentFunc(func(accountsBranch treeout.Branches) {
						accountsBranch.Child(format.Meta("                 payer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(format.Meta("associatedTokenAddress", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(format.Meta("                wallet", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(format.Meta("             tokenMint", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(format.Meta("         systemProgram", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(format.Meta("          tokenProgram", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(format.Meta("            sysVarRent", inst.AccountMetaSlice.Get(6)))
					})
				})
		})
}

func NewCreateIdempotentInstruction(
	Payer solana.PublicKey,
	associatedToken solana.PublicKey,
	Wallet solana.PublicKey,
	Mint solana.PublicKey,
	programID solana.PublicKey,
) *CreateIdempotent {
	return NewCreateIdempotentInstructionBuilder().
		SetPayer(Payer).
		SetAssociatedToken(associatedToken).
		SetWallet(Wallet).
		SetMint(Mint).
		SetProgramID(programID)
}

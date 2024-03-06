package ed25519_sig_verify

import (
	"encoding/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/text/format"
	"github.com/gagliardetto/treeout"
	"math"
)

var ProgramID = solana.Ed25519ProgramID

const ProgramName = "Ed25519SigVerify"

type Instruction struct {
	Signer    solana.PublicKey
	Signature solana.Signature
	Message   []byte
}

func (inst *Instruction) EncodeToTree(parent treeout.Branches) {
	parent.Child(format.Program(ProgramName, ProgramID)).
		ParentFunc(func(programBranch treeout.Branches) {
			programBranch.Child(format.Instruction("Ed25519SigVerify")).
				ParentFunc(func(instructionBranch treeout.Branches) {
					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch treeout.Branches) {
						paramsBranch.Child(format.Param("Signer", inst.Signer))
						paramsBranch.Child(format.Param("Signature", inst.Signature))
						paramsBranch.Child(format.Param("Message", inst.Message))
					})
				})
		})
}

func (inst *Instruction) ProgramID() solana.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() []*solana.AccountMeta {
	return nil
}

func (inst *Instruction) Data() ([]byte, error) {
	data := make([]byte, 0, 16+32+64+len(inst.Message))
	data = append(data, 1) // num of signatures
	data = append(data, 0) // padding

	signerOffset := uint16(16)
	sigOffset := signerOffset + 32
	msgOffset := sigOffset + 64
	data = binary.LittleEndian.AppendUint16(data, sigOffset)
	data = binary.LittleEndian.AppendUint16(data, math.MaxUint16)
	data = binary.LittleEndian.AppendUint16(data, signerOffset)
	data = binary.LittleEndian.AppendUint16(data, math.MaxUint16)
	data = binary.LittleEndian.AppendUint16(data, msgOffset)
	data = binary.LittleEndian.AppendUint16(data, uint16(len(inst.Message)))
	data = binary.LittleEndian.AppendUint16(data, math.MaxUint16)

	data = append(data, inst.Signer[:]...)
	data = append(data, inst.Signature[:]...)
	data = append(data, inst.Message...)

	return data, nil
}

func NewEd25519SigVerifyInstruction(signer solana.PublicKey, sig solana.Signature, msg []byte) *Instruction {
	return &Instruction{
		Signer:    signer,
		Signature: sig,
		Message:   msg,
	}
}

func NewEd25519SigVerifyInstructionWithWallet(signer *solana.Wallet, msg []byte) (*Instruction, error) {
	sig, err := signer.PrivateKey.Sign(msg)
	if err != nil {
		return nil, err
	}
	return &Instruction{
		Signer:    signer.PublicKey(),
		Signature: sig,
		Message:   msg,
	}, nil
}

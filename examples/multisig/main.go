// Command multisig is a small, runnable explanation of how ClueMultiSig's
// _validatorSignUpdate uses updateHash.
//
// It mirrors the Solidity code off-chain:
//
//	approvalDigest(code, version, updateHash)
//	  = keccak256(0x1901 || domainSeparator || structHash)
//	structHash = keccak256(abi.encode(APPROVAL_TYPEHASH, code, version, updateHash))
//
// Run with:
//
//	go run ./examples/multisig
package main

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// keccak256("Approval(uint32 code,uint32 version,bytes32 updateHash)")
	approvalTypehash = crypto.Keccak256Hash([]byte(
		"Approval(uint32 code,uint32 version,bytes32 updateHash)",
	))

	// EIP-712 domain used by ClueMultiSig: EIP712("ClueRegistry", "1").
	domainTypehash = crypto.Keccak256Hash([]byte(
		"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)",
	))
	nameHash    = crypto.Keccak256Hash([]byte("ClueRegistry"))
	versionHash = crypto.Keccak256Hash([]byte("1"))
)

// approvalDigest is the 32-byte hash a validator signs. It is exactly what
// ClueMultiSig.approvalDigest() returns on-chain.
func approvalDigest(
	chainID *big.Int,
	verifying common.Address,
	code, version uint32,
	updateHash common.Hash,
) common.Hash {
	domainSeparator := crypto.Keccak256Hash(pack(
		[]string{"bytes32", "bytes32", "bytes32", "uint256", "address"},
		domainTypehash, nameHash, versionHash, chainID, verifying,
	))

	structHash := crypto.Keccak256Hash(pack(
		[]string{"bytes32", "uint32", "uint32", "bytes32"},
		approvalTypehash, code, version, updateHash,
	))

	// EIP-712: "\x19\x01" || domainSeparator || structHash
	buf := make([]byte, 0, 2+32+32)
	buf = append(buf, 0x19, 0x01)
	buf = append(buf, domainSeparator.Bytes()...)
	buf = append(buf, structHash.Bytes()...)
	return crypto.Keccak256Hash(buf)
}

func pack(types []string, values ...any) []byte {
	args := make(abi.Arguments, len(types))
	for i, t := range types {
		typ, err := abi.NewType(t, "", nil)
		if err != nil {
			panic(err)
		}
		args[i] = abi.Argument{Type: typ}
	}
	encoded, err := args.Pack(values...)
	if err != nil {
		panic(err)
	}
	return encoded
}

// updateHashFromProposal builds the proposal commitment the registry would pass
// in. The real registry can cover whatever it wants; here we use the three
// fields ClueMultiSig documents: previous root, proposed root and batch id.
func updateHashFromProposal(prevRoot, newRoot common.Hash, batchID uint64) common.Hash {
	var id [32]byte
	binary.BigEndian.PutUint64(id[24:], batchID)
	return crypto.Keccak256Hash(prevRoot.Bytes(), newRoot.Bytes(), id[:])
}

// ---------------------------------------------------------------------------
// Off-chain mirror of the relevant on-chain state
// ---------------------------------------------------------------------------

type validator struct {
	id      common.Hash
	account common.Address
	active  bool
}

// session mirrors the storage used by _validatorSignUpdate for one open update.
type session struct {
	code     uint32
	version  uint32
	required int
	signers  map[common.Hash]struct{}
}

type registry struct {
	chainID  *big.Int
	verifier common.Address
	byAddr   map[common.Address]validator
	linked   map[uint32]map[common.Hash]struct{} // code => validator ids
}

// signUpdate mirrors ClueMultiSig._validatorSignUpdate.
//
// The important line is the digest: it commits to (code, version, updateHash).
// Recovering the signer over that digest proves the validator approved that
// exact proposal. If the signature was made over a different proposal, the
// recovered address is not the validator and the call is rejected.
func (r *registry) signUpdate(
	s *session,
	updateHash common.Hash,
	signature []byte,
) (bool, error) {
	digest := approvalDigest(r.chainID, r.verifier, s.code, s.version, updateHash)

	pub, err := crypto.SigToPub(digest.Bytes(), signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %w", err)
	}
	signer := crypto.PubkeyToAddress(*pub)

	v, ok := r.byAddr[signer]
	if !ok || !v.active {
		return false, fmt.Errorf("NotCodeValidator(%d, %s)", s.code, signer.Hex())
	}
	if _, ok := r.linked[s.code][v.id]; !ok {
		return false, fmt.Errorf("NotCodeValidator(%d, %s)", s.code, signer.Hex())
	}
	if _, ok := s.signers[v.id]; ok {
		return false, fmt.Errorf("AlreadySigned(%d, %s)", s.code, v.id.Hex())
	}

	s.signers[v.id] = struct{}{}
	return len(s.signers) >= s.required, nil
}

func main() {
	// Same domain as ClueMultiSig: EIP712("ClueRegistry", "1").
	chainID := big.NewInt(31337)
	verifier := common.HexToAddress("0x000000000000000000000000000000000000C0DE")

	const (
		code    = uint32(10001)
		version = uint32(7)
	)

	// Four validators are linked to the code. The contract's policy is
	// two thirds plus one, capped at the total, so 4 validators need 3 signatures.
	keys := make([]*ecdsa.PrivateKey, 4)
	reg := &registry{
		chainID:  chainID,
		verifier: verifier,
		byAddr:   make(map[common.Address]validator),
		linked:   map[uint32]map[common.Hash]struct{}{code: {}},
	}
	for i := range keys {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		keys[i] = key

		id := crypto.Keccak256Hash([]byte(fmt.Sprintf("validator-%d", i+1)))
		v := validator{id: id, account: crypto.PubkeyToAddress(key.PublicKey), active: true}
		reg.byAddr[v.account] = v
		reg.linked[code][id] = struct{}{}
		fmt.Printf("validator %d: id=%s account=%s\n", i+1, id.Hex(), v.account.Hex())
	}

	s := &session{
		code:     code,
		version:  version,
		required: 3, // 2/3 of 4 + 1
		signers:  make(map[common.Hash]struct{}),
	}

	// The proposal the registry wants approved.
	updateHash := updateHashFromProposal(
		common.HexToHash("0x1111"),
		common.HexToHash("0x2222"),
		42,
	)
	fmt.Printf("\nupdateHash (proposal commitment) = %s\n", updateHash.Hex())
	fmt.Printf("approvalDigest                  = %s\n", approvalDigest(chainID, verifier, code, version, updateHash).Hex())

	// Each validator signs the same digest, i.e. the same updateHash.
	fmt.Println("\n-- validators sign the same proposal --")
	for i := 0; i < 3; i++ {
		digest := approvalDigest(chainID, verifier, code, version, updateHash)
		sig, err := crypto.Sign(digest.Bytes(), keys[i])
		if err != nil {
			panic(err)
		}
		done, err := reg.signUpdate(s, updateHash, sig)
		if err != nil {
			fmt.Printf("validator %d: rejected: %v\n", i+1, err)
			continue
		}
		fmt.Printf("validator %d: accepted, signatures=%d, fulfilled=%v\n", i+1, len(s.signers), done)
	}

	// A signature over a different updateHash does not verify for this proposal.
	fmt.Println("\n-- validator 4 signs a DIFFERENT proposal --")
	otherUpdateHash := updateHashFromProposal(
		common.HexToHash("0x1111"),
		common.HexToHash("0x9999"), // different proposed root
		42,
	)
	otherDigest := approvalDigest(chainID, verifier, code, version, otherUpdateHash)
	wrongSig, err := crypto.Sign(otherDigest.Bytes(), keys[3])
	if err != nil {
		panic(err)
	}
	if _, err := reg.signUpdate(s, updateHash, wrongSig); err != nil {
		fmt.Printf("validator 4: rejected: %v\n", err)
	}

	// The same validator cannot sign twice.
	fmt.Println("\n-- replay: validator 1 signs the same proposal again --")
	digest := approvalDigest(chainID, verifier, code, version, updateHash)
	replaySig, err := crypto.Sign(digest.Bytes(), keys[0])
	if err != nil {
		panic(err)
	}
	if _, err := reg.signUpdate(s, updateHash, replaySig); err != nil {
		fmt.Printf("validator 1: rejected: %v\n", err)
	}

	fmt.Println("\nDone: updateHash binds each signature to one exact proposal.")
}

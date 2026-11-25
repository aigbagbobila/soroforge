package stellar_test

import (
	"fmt"

	"github.com/soroworks/soroforge/internal/stellar"
	"github.com/stellar/go-stellar-sdk/xdr"
)

// ExampleNewMockSigner creates a signer backed by a generated keypair.
func ExampleNewMockSigner() {
	signer := stellar.NewMockSigner()
	fmt.Println(len(signer.Address()) == 56)

	// Output:
	// true
}

// ExampleNewMockSignerFromSeed creates a signer with a stable address from a
// Stellar secret seed, which is useful when a test needs repeatable identities.
func ExampleNewMockSignerFromSeed() {
	signer, err := stellar.NewMockSignerFromSeed("SBFGFF27Y64ZUGFAIG5AMJGQODZZKV2YQKAVUUN4HNE24XZXD2OEUVUP")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(signer.Address()) == 56)

	// Output:
	// true
}

// ExampleNewMultiSigner combines a primary signer with one or more
// co-signers for transactions that require multiple signatures.
func ExampleNewMultiSigner() {
	primary, err := stellar.NewMockSignerFromSeed("SBFGFF27Y64ZUGFAIG5AMJGQODZZKV2YQKAVUUN4HNE24XZXD2OEUVUP")
	if err != nil {
		fmt.Println(err)
		return
	}
	coSigner, err := stellar.NewMockSignerFromSeed("SCN4VT2ELMFTALLNPPJNCK6OGQEWLERZIMEPARTRVDXXL2IJAPWSXZQI")
	if err != nil {
		fmt.Println(err)
		return
	}
	multi, err := stellar.NewMultiSigner(primary, coSigner)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(multi.Signers())

	// Output:
	// 2
}

// ExampleContractIDFromSimulation recovers the contract ID a create-contract
// simulation reports, so a dry run can print the address a real deploy would
// produce before anything is submitted.
func ExampleContractIDFromSimulation() {
	// The simulated return value is the new contract's address, encoded as
	// base64 XDR. ContractIDFromSimulation decodes it back into a C... strkey.
	encoded, err := xdr.MarshalBase64(xdr.ScVal{
		Type: xdr.ScValTypeScvAddress,
		Address: &xdr.ScAddress{
			Type: xdr.ScAddressTypeScAddressTypeContract,
			ContractId: func() *xdr.ContractId {
				raw := make([]byte, 32)
				for i := range raw {
					raw[i] = byte(i + 1)
				}
				id := xdr.ContractId(raw)
				return &id
			}(),
		},
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	got, err := stellar.ContractIDFromSimulation(stellar.SimulateResult{
		Results: []stellar.HostFunctionResult{{ReturnValueXDR: encoded}},
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(got)

	// Output:
	// CAAQEAYEAUDAOCAJBIFQYDIOB4IBCEQTCQKRMFYYDENBWHA5DYPSBFLM
}

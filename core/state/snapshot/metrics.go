// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package snapshot

import "github.com/ethereum/go-ethereum/metrics"

// Metrics in generation
var (
	snapGeneratedAccountMeter     = metrics.NewRegisteredMeter("state/snapshot/generation/account/generated", nil)
	snapRecoveredAccountMeter     = metrics.NewRegisteredMeter("state/snapshot/generation/account/recovered", nil)
	snapWipedAccountMeter         = metrics.NewRegisteredMeter("state/snapshot/generation/account/wiped", nil)
	snapMissallAccountMeter       = metrics.NewRegisteredMeter("state/snapshot/generation/account/missall", nil)
	snapGeneratedStorageMeter     = metrics.NewRegisteredMeter("state/snapshot/generation/storage/generated", nil)
	snapRecoveredStorageMeter     = metrics.NewRegisteredMeter("state/snapshot/generation/storage/recovered", nil)
	snapWipedStorageMeter         = metrics.NewRegisteredMeter("state/snapshot/generation/storage/wiped", nil)
	snapMissallStorageMeter       = metrics.NewRegisteredMeter("state/snapshot/generation/storage/missall", nil)
	snapDanglingStorageMeter      = metrics.NewRegisteredMeter("state/snapshot/generation/storage/dangling", nil)
	snapSuccessfulRangeProofMeter = metrics.NewRegisteredMeter("state/snapshot/generation/proof/success", nil)
	snapFailedRangeProofMeter     = metrics.NewRegisteredMeter("state/snapshot/generation/proof/failure", nil)

	// snapAccountProveCounter measures time spent on the account proving
	snapAccountProveCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/prove", nil)
	// snapAccountTrieReadCounter measures time spent on the account trie iteration
	snapAccountTrieReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/trieread", nil)
	// snapAccountSnapReadCounter measures time spent on the snapshot account iteration
	snapAccountSnapReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/snapread", nil)
	// snapAccountWriteCounter measures time spent on writing/updating/deleting accounts
	snapAccountWriteCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/account/write", nil)
	// snapStorageProveCounter measures time spent on storage proving
	snapStorageProveCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/prove", nil)
	// snapStorageTrieReadCounter measures time spent on the storage trie iteration
	snapStorageTrieReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/trieread", nil)
	// snapStorageSnapReadCounter measures time spent on the snapshot storage iteration
	snapStorageSnapReadCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/snapread", nil)
	// snapStorageWriteCounter measures time spent on writing/updating storages
	snapStorageWriteCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/write", nil)
	// snapStorageCleanCounter measures time spent on deleting storages
	snapStorageCleanCounter = metrics.NewRegisteredCounter("state/snapshot/generation/duration/storage/clean", nil)
)

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

package params

// These are network parameters that need to be constant between clients, but
// aren't necessarily consensus related.

const (
	// BloomBitsBlocks is the number of blocks a single bloom bit section vector
	// contains on the server side.
	BloomBitsBlocks uint64 = 4096

	// BloomBitsBlocksClient is the number of blocks a single bloom bit section vector
	// contains on the light client side
	BloomBitsBlocksClient uint64 = 32768

	// BloomConfirms is the number of confirmation blocks before a bloom section is
	// considered probably final and its rotated bits are calculated.
	BloomConfirms = 256

	// CHTFrequency is the block frequency for creating CHTs
	CHTFrequency = 32768

	// BloomTrieFrequency is the block frequency for creating BloomTrie on both
	// server/client sides.
	BloomTrieFrequency = 32768

	// HelperTrieConfirmations is the number of confirmations before a client is expected
	// to have the given HelperTrie available.
	HelperTrieConfirmations = 2048

	// HelperTrieProcessConfirmations is the number of confirmations before a HelperTrie
	// is generated
	HelperTrieProcessConfirmations = 256

	// CheckpointFrequency is the block frequency for creating checkpoint
	CheckpointFrequency = 32768

	// CheckpointProcessConfirmations is the number before a checkpoint is generated
	CheckpointProcessConfirmations = 256

	// FullImmutabilityThreshold is the number of blocks after which a chain segment is
	// considered immutable (i.e. soft finality). It is used by the downloader as a
	// hard limit against deep ancestors, by the blockchain against deep reorgs, by
	// the freezer as the cutoff threshold and by clique as the snapshot trust limit.
	FullImmutabilityThreshold = 90000

	// LightImmutabilityThreshold is the number of blocks after which a header chain
	// segment is considered immutable for light client(i.e. soft finality). It is used by
	// the downloader as a hard limit against deep ancestors, by the blockchain against deep
	// reorgs, by the light pruner as the pruning validity guarantee.
	LightImmutabilityThreshold = 30000
)

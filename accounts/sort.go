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

package accounts

// AccountsByURL implements sort.Interface for []Account based on the URL field.
type AccountsByURL []Account

func (a AccountsByURL) Len() int           { return len(a) }
func (a AccountsByURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AccountsByURL) Less(i, j int) bool { return a[i].URL.Cmp(a[j].URL) < 0 }

// WalletsByURL implements sort.Interface for []Wallet based on the URL field.
type WalletsByURL []Wallet

func (w WalletsByURL) Len() int           { return len(w) }
func (w WalletsByURL) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WalletsByURL) Less(i, j int) bool { return w[i].URL().Cmp(w[j].URL()) < 0 }

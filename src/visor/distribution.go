package visor

import (
	"github.com/skycoin/skycoin/src/coin"
)

const (
	// MaxCoinSupply is the maximum supply of skycoins
	MaxCoinSupply uint64 = 1e8 // 100,000,000 million

	// DistributionAddressesTotal is the number of distribution addresses
	DistributionAddressesTotal uint64 = 100

	// DistributionAddressInitialBalance is the initial balance of each distribution address
	DistributionAddressInitialBalance uint64 = MaxCoinSupply / DistributionAddressesTotal

	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64 = 25

	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 = 5

	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 = 60 * 60 * 24 * 365 // 1 year
)

func init() {
	if MaxCoinSupply%DistributionAddressesTotal != 0 {
		panic("MaxCoinSupply should be perfectly divisible by DistributionAddressesTotal")
	}
}

// GetDistributionAddresses returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func GetDistributionAddresses() []string {
	addrs := make([]string, len(distributionAddresses))
	for i := range distributionAddresses {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetUnlockedDistributionAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, InitialUnlockedCount)
	for i := range distributionAddresses[:InitialUnlockedCount] {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetLockedDistributionAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddresses[InitialUnlockedCount:] {
		addrs[i] = distributionAddresses[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}

// TransactionIsLocked returns true if the transaction spends locked outputs
func TransactionIsLocked(inUxs coin.UxArray) bool {
	lockedAddrs := GetLockedDistributionAddresses()
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true
		}
	}

	return false
}

var distributionAddresses = [DistributionAddressesTotal]string{
	"28CCuSx9sFrrEzbK9ULqgz8qgiqL8vFsxTN",
	"i4HX3qGsypDAXHnApGKWECZrHnhMftHe69",
	"2HPQhxEXfrCT9satjjUdQWDHk7J64Q2htQ",
	"mcnbaKqXEdSTybh52VJpbs3dWmurjTkpFo",
	"213CtaRVMiSv1aounkkg4aAiuYT3Bzz2GWT",
	"HheZzwoD3YDYr1uCnHAokvDDy9iNv8wTgr",
	"6Seww8Yb5YMRyeFSNRh1j5YVdYWtQHAzDu",
	"7AWxMkzW4Sim1Jt2nt6v8wZJFMoLL1YKHr",
	"8LWfckUvznYuYwivn35SFzr8ZomK5m4YZ6",
	"2Qfunxj1yjCgT2hGtZvU8dBPXhB8RP98Hk8",
	"W9QR1ZSfHEAorzaWRctydf7cyWTrZ84H4y",
	"5GPfZWx2WPVMeXhugaRcirUZxiP4T2kNgD",
	"rUbAFoDgzEuehKJspYPS6ic9ZpsFTEbmhk",
	"2GexsVMzJaBTUhsV2hJ6js6eyEaniehtVsK",
	"6G7fB6tCDXD62KqijHDU9bJ7Hp3dPpACZ4",
	"27H94tpGdkhzXdhhHZPYf7AvbfnAVoXf6W5",
	"FsFZgPnknT2ECNyoDNSyerK1Npwc1nY4mK",
	"2fLU8mSxqBeDvwC5FzVP2G8uxxvWgeVWoSf",
	"28RNnGPRKPWVEEx8SYMJZBhHDhJc2b4MkwG",
	"enw285y4Ni67Zj7qCAAZByKAmLMJciAczr",
	"qiueAfo7PsFFAzLZ9uJDJL6E9trCzTvpSR",
	"8z7vb352nM1mPpBgDaHCFfvqC4nxcgG8Mh",
	"4toGZyo4puvULTBjm62wW4d7PpzhyTDUG6",
	"WupXKMoxyiXu8v1hNrtJ3zGGSwbnQMNj6X",
	"NKg2Pnm7Njj9qP5e3gwPmr5eL8gKZhmhaH",
	"GMW11mHXCdUgLbPsZ7mf8Ku5LPbajiksfG",
	"RGaybkrzarSRKMonssXJtFyi936szVp7Ms",
	"21XQ3XwZBbVxKogy229VMrGbUqffQCWtcm5",
	"VghNSgu7SYqViaTXbmJm4ZW6DDLXVfLepE",
	"YEuVvit948yb347v4bUnbedQm9V3mG57dt",
	"bzfeURqdPHznE3fLSn2xcYD39ygH8Exexs",
	"2XdZ6TB2euYGna5XNZb9pX4A57zdUPh8n46",
	"2dNHsbSTtLP6K7ciWaUGrcuMmS7SzjWhF3C",
	"2iBYPwbhERRpZ7nVqYVjJ79YxhsECai5RYE",
	"2gEfzvVW4YA55hh2H4amV6oPWxyxSrNqPTu",
	"232jEzg9BxqocX9DqXCsEPnoaVQojHKAJpT",
	"HPk49RkzaiSAx5RKc2xv9i2y8oTkuzLFET",
	"td6XoUrzGQXJM9CgvegMg2vZfh23tHQtfP",
	"2cxxm78FZJ7CLz7qLbJ2RLsjEojjHvWgyxi",
	"XDX5Am9odjFCSEsj93kmBUkravc839fPb9",
	"2EwkUiQBDkABDbPaNdJNnq6fCA6gmpoAnF1",
	"b4QDJej8HrNu7zDbxzKgTTxbbYGq9j9BPP",
	"2CNKpzximCzZ3ELT8WA9mVCFjigLdmKD49g",
	"9qY8YAJfauSE8vwC8n1fnoi737kBRYVK8a",
	"o9WGKRY9Jq1DccLLgLZmEzfiDkucvz4Kpm",
	"ZojwhmXhfEnF7wnWBaSBpxFxdrTZtkw1Np",
	"SAiratAUpayEAdRmXusiXSm1VgcwwJyDhx",
	"767DCkuYFtrBk6K9CaNwSRSXQcEaJCoy6g",
	"HkdpeAK1eNWy2Ftd2X5tTNPzAgat25Qs3r",
	"D8ZJy1AKxdE2tD8hbez29mzRvNDGhGkTNS",
	"2QQVynZjAoukTNbeC4ATorYzQzsD6cKJ2uN",
	"Sbdrj6FrQ97Yy7hPYG3K11uhjMa5TRq3cQ",
	"25b8q1zksdv4UfBZRqXhHPCmeX7rBiATz5W",
	"2e475VD7beBtFZCFgKLq4TAqnGWzQWJQGgn",
	"2kKsm4qnUZHpUyYUKe5NGMpZ6JhWCGV5VFJ",
	"2cjiNJsNvHw2ejmaddc9EaiddTD5ir8MU85",
	"8mawM2tU2fepoQrsuW3qtm8erahER8hBcE",
	"2afjiWhctmMuAraVA7bEEVjJ5m2Pr9e8AeZ",
	"zZs6MUj4yUgDRhyBz8aQCyLGZo3cnFt8Kb",
	"2j5F9EzAFs3vGKTubrXj6QnXm58r4emBuXR",
	"vEvYkBa6K6bJKRJt1aLXMDT5nY6diJf9xa",
	"2mjXG5W3gysdzn9FxwuJBU5i4hMpy3sUBkC",
	"SjNrU3MD59Sn7yYFjnJHsRCAspkVDcoRd8",
	"hT9GzX5rNoRn4XJeDL3fgSAZfzQWH5bKfj",
	"2Xt4VbydGKuGqjVjJy3rR48qRwzxV4tq35M",
	"yfZzLTEQvYayqzHE5GF2bstuEkVQY6p9j9",
	"2cfVjqaDyuNn8JgJTWngcNGTwLiagt7JGK2",
	"4ysDqyBnfuHgSMdwAZEoi77EVprCTSeF5W",
	"55wMfEq3hps2fXCG9JDQfCusfFVkGgbPoJ",
	"Fta7nsbNLq5B5fxne3tMzYq6redv1rrdNs",
	"2ZY95NYA8aeWt6jvup8f8Ne2Mi79hJ55ubf",
	"2dxevugNJP5Nz5a9c6NeFbMrceWrWE3vMSm",
	"2AaHHc9Vn9kfZkVLqXvuWSVZMu59Rp3fnu7",
	"2a39tEFHTm6bhSqNzfeSV2xGLwtJGBReRtw",
	"rYMir7DBZr7WK9bfp1NVnzj6d8ZXasLfrs",
	"SuWhNNwRAhNWD8WKFWTuKHWFdKofS1tSev",
	"2Qcofpv49hJn4dvZsCsKWuLdJ1c3LYSsXzW",
	"2FBjxaCu3SiQiK6HgbodSTuLqMymGFyAC2N",
	"TjMnNxqCeZq2ryQLcdAQwHpfV6j1Aj8MtB",
	"2hZCPnDNAhcrwaykWpp4rKZxGQgwimq6TbZ",
	"dtHE76cCKXCqNPpXuDvwnRwUnwCU8nQqyM",
	"2L24Lk91pSH47sBXJNRkhspyVVaA1H2d4Tu",
	"CJUbizmyVpZhPbd69sYURHqXGM2zyZsvRX",
	"25cVCuLEAADn6sPuU3uuVC9wXWaFmhox1r5",
	"nmfnoKYEY65BpBWdTSfzkFuv85V3UBrXan",
	"BG5LAUF9LrdQKdpcCCduUJfDudpDvwk4U6",
	"KKxojNPNAyG5ycorsmrTWyES8wm7wn3mER",
	"FwbRuahhZfPchDMHVYkDb4Svnhqk7WU9PS",
	"KVp2G1CycnS8RZX3Y5vSyAM8xwP8vXckRK",
	"rgCd3y7LmVx1meSx8G5ervkrvKWH2EbjND",
	"29fufhLhmPNgH1d1UrwCKNRAcfBahXLdUV8",
	"kFEwCVMbQp6eLXhcPXBtTyu6xdSK4n99ea",
	"22N6F7MLaFUw2wUijg7oy8BKBP6Pz1ewRRQ",
	"2CaZPUh8imzZPZK2gTTaRxKTrP5GpWsTayF",
	"VPNUHKVgqNmUWgMgX1dauJK3gapmYxc5wX",
	"H8XFoDqZsNmsSpbVsSGBHYehHshtg9Ywop",
	"2BXdT1qccKhcjr1LoR26xis6oVkRgJfSKdh",
	"2V8o4Fs264HuTJoyRrtQv91Ls1zwcL7gMta",
	"7sYXsdjkfSFn8iFYqwVswPD79EPZREFjcR",
	"2XYoumaq6izfbzsWbsCjc1r1XCRmM9DES1H",
}

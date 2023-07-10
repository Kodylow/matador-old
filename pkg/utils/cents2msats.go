package utils

const BTCPRICE float64 = 25000 // $25,000.00 per bitcoin

func CentsToMsats(cents uint64) uint64 {
    btc := float64(cents) / (BTCPRICE * 100) // Convert cents to BTC
    msats := btc * 100000000000 // Convert BTC to msats
    return uint64(msats)
}
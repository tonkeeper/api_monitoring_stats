package dapps

var StonFi = &Dapp{
	name:        "StonFi",
	dAppUrl:     "https://app.ston.fi/swap?ft=jUSDT&tt=STON",
	calcApiUrl:  "https://rpc.ston.fi/",
	calcPayload: pointer(`{"jsonrpc":"2.0","id":1,"method":"dex.simulate_swap","params":{"offer_address":"EQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAM9c","offer_units":"1000000000","ask_address":"EQA2kCVNwVsil2EM2mB0SkXytxCqQjS4mttjDpnXmwG9T6bO","slippage_tolerance":"0.001"}}`),
}

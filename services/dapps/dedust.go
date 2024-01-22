package dapps

var DeDust = &Dapp{
	name:        "DeDust",
	dAppUrl:     "https://dedust.io/swap",
	calcApiUrl:  "https://api.dedust.io/v2/routing/plan",
	calcPayload: pointer(`{"from":"native","to":"jetton:0:65aac9b5e380eae928db3c8e238d9bc0d61a9320fdc2bc7a2f6c87d6fedf9208","amount":"1000000000"}`),
}

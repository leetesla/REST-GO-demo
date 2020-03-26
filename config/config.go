package config

// API KEY
var (
	// todo: replace with your own AccessKey and Secret Key
	ACCESS_KEY string = "39a3669c-xa2b53ggfc-4cfbfba6-1435e"
	SECRET_KEY string = "fcedbbb7-126df733-eb2f96d0-8c9af"

	// default to be disabled, please DON'T enable it unless it's officially announced.
	ENABLE_PRIVATE_SIGNATURE bool = false

	// generated the key by: openssl ecparam -name prime256v1 -genkey -noout -out privatekey.pem
	// only required when Private Signature is enabled
	// todo: replace with your own PrivateKey from privatekey.pem
	PRIVATE_KEY_PRIME_256 string = `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

)

// API请求地址, 不要带最后的/
const (
	//todo: replace with real URLs and HostName
	MARKET_URL string = "https://api-aws.huobi.pro"
	TRADE_URL  string = "https://api-aws.huobi.pro"
	HOST_NAME  string = "api-aws.huobi.pro"  
)

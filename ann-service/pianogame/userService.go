package pianogame

import "fmt"

func init() {
	// This is blank
	// fmt.Println(SysConfig)
	// fmt.Println(Ssl.Path.Cert)
	fmt.Println(authSettings.Secret.Jwt)
}

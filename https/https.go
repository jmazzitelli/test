package main

import (
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
)

func loadCert() (certFile string, err error) {
	certFile = "/tmp/mycert.cer"

	data := `
	-----BEGIN CERTIFICATE-----
	MIIEUDCCArigAwIBAgIQfaOk5w7VpJ4P80/N9jd3HzANBgkqhkiG9w0BAQsFADBX
	MR4wHAYDVQQKExVta2NlcnQgZGV2ZWxvcG1lbnQgQ0ExFjAUBgNVBAsMDWRwYW5p
	Y0BtYXN0ZXIxHTAbBgNVBAMMFG1rY2VydCBkcGFuaWNAbWFzdGVyMB4XDTE5MDcw
	NTA5MTc0N1oXDTI5MDcwNTA5MTc0N1owPzEnMCUGA1UEChMebWtjZXJ0IGRldmVs
	b3BtZW50IGNlcnRpZmljYXRlMRQwEgYDVQQLDAtyb290QG1hc3RlcjCCASIwDQYJ
	KoZIhvcNAQEBBQADggEPADCCAQoCggEBAOh3Fa96rHPhyP8fhMVYGnDI52E0XWF4
	K4M4KBhksRC/pXxG/c64Gci8+ZL1cUy5zV4hpeA7W0Xzi3r39p9ogucx/9xAyIox
	KbcmrrjHMdBWzLB/2pKXdx/7b5Kqj/Z37tQLQHX2AiAyIMmoqg1d8cRpvwk0eME4
	ImKwPEk8xCp4/QusBWJ1kjSAx8Mk4f/Lj4imxt5JuEQptUmjIWtRfikEtmqknJfR
	NvTjDRKV8p0S/cqQGqnqgPSQ+ktWHH4Rj/bpsyKPtqFMfrKKqcrJzoaqLF2kPNE0
	KROn16CAXPdMZtcnGrn9543dVQQlfLy7A1nnhhesuQz1hX2K7toa1esCAwEAAaOB
	rzCBrDAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0T
	AQH/BAIwADAfBgNVHSMEGDAWgBSYBf6xg5HXiLwo8iDQWui72UJ5/DBWBgNVHREE
	TzBNggtleGFtcGxlLmNvbYINKi5leGFtcGxlLmNvbYIMZXhhbXBsZS50ZXN0ggls
	b2NhbGhvc3SHBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAEwDQYJKoZIhvcNAQELBQAD
	ggGBAG9XMVDkELlZCJZ32tdXOEyIufz5XoDlOCHWqXLS57Oo6up0ktushC8sxzdq
	FeoGe/Qu+jDlNdhKHMMneeRc38myZQRjyA7iA8gIG86Fj8di4zF2wpuQ5i8M/5Rq
	rUzexZ3JKHUWh0omlHx1hXdJi9o34iIS9a5cEvwyGU39e8ag2vSO/voCkSCFheS7
	peeuvQg400V0ii5oY39KlvocfKK2yEgCTF1LJjIc6YTV6aemPpBv9g2tNHTJTXKj
	MPdLXU85ArVLED2rNpWHnto7w+pXxQgoHgzDqcKRQXt7B5aovAMC+09ATBTF8Dxr
	iXkl78d6P2SggDBzft6jPXqcQzmaSyiSF4B9UYgNzfObM+pkaV6v17SEDwxJfu+N
	u+cgt9S5nhwxqZkhT6pK8EsgqqrBlJoRAeKISL0oHMOwlvUT7ceI5LKJWY58kxVC
	FlGkxUHu+WvJx2+P+etq1FrCmIHmAzLMGuZyp+IXiGHtPEwzT7amnKhJ7EDRrfEG
	zG4RJQ==
	-----END CERTIFICATE-----`

	data = strings.ReplaceAll(data, "\t", "")
	d1 := []byte(data)
	err = ioutil.WriteFile(certFile, d1, 0600)

	return
}

func loadKey() (keyFile string, err error) {
	keyFile = "/tmp/mykey.key"

	data := `
	-----BEGIN PRIVATE KEY-----
	MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDodxWveqxz4cj/
	H4TFWBpwyOdhNF1heCuDOCgYZLEQv6V8Rv3OuBnIvPmS9XFMuc1eIaXgO1tF84t6
	9/afaILnMf/cQMiKMSm3Jq64xzHQVsywf9qSl3cf+2+Sqo/2d+7UC0B19gIgMiDJ
	qKoNXfHEab8JNHjBOCJisDxJPMQqeP0LrAVidZI0gMfDJOH/y4+IpsbeSbhEKbVJ
	oyFrUX4pBLZqpJyX0Tb04w0SlfKdEv3KkBqp6oD0kPpLVhx+EY/26bMij7ahTH6y
	iqnKyc6GqixdpDzRNCkTp9eggFz3TGbXJxq5/eeN3VUEJXy8uwNZ54YXrLkM9YV9
	iu7aGtXrAgMBAAECggEADVKl34S8VXffOR/pUBYYjdY1zJBfubJVbBPP2HYM39Tb
	+x9mdG6Aq8yI0S9X6vnLF1X+V7ePJ5cpq0aCz+gBeJaY/1qHI8Rli6Wf5d8kr7gJ
	yyPItxYPMboLTvCPh6Sf/28Vpq0OuiGlV2lfNZzoukUFOdXUBd7duaI4Ekp1Q6nH
	gAQjB2iSauwpq4FGMZDvkAApJdC8FSYbZg8G3jSMYWMj6/gW0YMJZFH9A6mSq2WF
	dKIrK6EjTNTNHI8DrZ6KPgGPX5jpzy88zqSFdu7VPErjhQvX+bW6Ypjfv1PUoVcI
	UoTmFMnpLj5JmE+7q+1WsxGkPTX/dgn6U2o6GI5N4QKBgQDxIwC/rf6h4PAvv85f
	A7erl5BY9+QesKUui0oGreqbkdmuwPp3U7HvqjzsMW1v6jdpN18rgTn1yUcyMf35
	ylT1KVKwnTvVNCiCFT3QDp7mEdtAJtZE9gLyxV48xuD0YxevDuAESfTW1dApZ+c1
	424DfpBdXht5/szkjRbiDkPgRwKBgQD2yz/rxJO2kC9/Oe+dPHpSkN4wbGcbOE9T
	ekw6R/AdWMS0TGydHUg/bdxKer7X90Zkt7/OUhvCpshsIT3yZz+/1/ab7tDsDWA9
	GYd1cbgefyjnT6hGimCksvBLgDqJwpN64UbGK1VEghtbsFPQ0rBUnfmxqd/9nxmb
	797Rch3zPQKBgELvhHWwxs4IsqOOiqq1TXbES71mklwyjKeu4o2YGVe11Mc9qkkV
	Yn80slSeI9K9IUSDqldZN82SYcD9P5LnJ04mel2sR7+XCueRHedzJ4iVzFaycSgT
	Yh4hy1bznd4444okhuqp3N0F3RKhVP0QdKljqI9CYD4tDJMk1wVJEG5hAoGASmVC
	w6Pik2orp0KjxNZyWWlqUVacTkxPPW7kg70j2PTldySCqWomWViYy6rs1NWp1rq9
	i0idLbRxPodW0Tfms8I6iQ8Y08/Ebya++txpEGhswC33ICyerYdzgI8LFnQdWTGH
	0D1H2vsNnDovSgf5N8jXeIMpDp9jbOqGVMT92lECgYAUzSp7QB/8zyVCyoToXck0
	flIs5pHjSB8xmtJ+hx9J5ulN39DkNaHXiPDcSEo0RCw/YO1CTpsvchtA2eo2edEC
	TPeyjdznmrgIGM+XloNuRiQE+VpuGDo+TqXSMuiJhtL9/JIO2CnPGGrTXcsWQtcT
	mS27nYYWXhe9zSekgrxCPg==
	-----END PRIVATE KEY-----`

	data = strings.ReplaceAll(data, "\t", "")
	d1 := []byte(data)
	err = ioutil.WriteFile(keyFile, d1, 0600)

	return
}

func main() {
	certFile, err := loadCert()
	if err != nil {
		panic(err)
	}

	keyFile, err := loadKey()
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS10,
		NextProtos: []string{"h2", "http/1.1"},
	}

	httpServer := &http.Server{
		Addr:      "0.0.0.0:8443",
		TLSConfig: cfg,
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello World!")
	})

	fmt.Println("Connect to server at: https://127.0.0.1:8443")
	if err := httpServer.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}

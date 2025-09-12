package api

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"golang.org/x/oauth2"
)

func Test_ExtractToken(t *testing.T) {
	tokenCookie := "eyJhY2Nlc3NfdG9rZW4iOiJ5YTI5LmEwQWNNNjEyeEJsV3JpSjE4eWt5T1VrbDFEZ1p2UEV2WkxMZF96TXY3eDc5NXREQXFYbjAwUlFiT3lEVS1fanBUSHZXd0ZhMDl6cUlXUjQ2SkQ3bXhPVFpGd0NUNzEyTnl2T3JiX3pUZnVQaWFiU09PUnc1cEstVEt6b1BIc25vQW9XaHRSR1NUMlA0ZTR0SlJseEpvV2hfUmwzcnYwOUFiUnBJcGxlbjZDcXZ5T1prRGdhV2lLYUNnWUtBZG9TQVJFU0ZRSEdYMk1pX0libmxZcjJDUFhLTEFqdkJpdnRCdzAxODciLCJ0b2tlbl90eXBlIjoiQmVhcmVyIiwicmVmcmVzaF90b2tlbiI6IjEvLzA2by1RRXY3UlB0emNDZ1lJQVJBQUdBWVNOd0YtTDlJcldncFZCSGNpODE4VHdrQktjTExaRWdGdFplaXZYTjhSQXE5c2Z3RjNwWWNUclZpcS1BOEI0NHZ4MUdNMW53OGdXWTgiLCJleHBpcnkiOiIyMDI0LTA5LTE3VDA2OjE4OjA5LjE3NDUxMzEwNFoifQ=="
	tokenObj := oauth2.Token{}

	tokenJson, _ := base64.StdEncoding.DecodeString(tokenCookie)
	_ = json.Unmarshal([]byte(tokenJson), &tokenObj)

	/* if isTokenValid(tokenObj) {
		if time.Now().Unix() >= (tokenObj.Expiry.Unix() - 60) {
			_, err := RefreshToken(tokenObj.RefreshToken)
			if err != nil {

				fmt.Print("asasd")
			}
			//update cookie
			//storeTokenInCookie(c, token)
			//return &token
		}
		url := "https://www.googleapis.com/youtube/v3/channels?part=id&mine=true"

		// Create a Bearer string by appending string access token
		var bearer = tokenObj.TokenType + " " + tokenObj.AccessToken

		// Create a new request using http
		req, err := http.NewRequest("GET", url, nil)

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			output := map[string]string{}
			output["loggedIn"] = "true"
			//return c.JSON(http.StatusOK, output)
		}
	}
	*/
}

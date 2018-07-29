package download

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
)

// Implementing the new PivNet API token system
// The below function extract the token ( UAA )
func (r *ResponseBody) getToken(uaa_token string) error {

	log.Info("Getting the access token from the UAA token provided.")
	r.Token.RefreshToken = uaa_token
	b, err := json.Marshal(&r.Token)
	if err != nil {
		return fmt.Errorf("failed to marshal API token request body: %s", err.Error())
	}

	// Placing request for access token.
	req, err := http.NewRequest("POST", RefreshToken, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to construct API token request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("API token request failed: %s", err.Error())
	}

	defer resp.Body.Close()

	// Is it a success or not
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch API token - received status %v", resp.StatusCode)
	}

	// Let store it for the rest of the program
	err = json.NewDecoder(resp.Body).Decode(&r.Token)
	if err != nil {
		return fmt.Errorf("failed to decode API token response: %s", err.Error())
	}

	return nil

}

package download

import (
	"io"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strconv"
	"os"
)

func (r *ResponseBody) ResponseChecker(statusCode int, body io.ReadCloser, url string) error {

	// If the status code is not 200, then error out
	if statusCode != http.StatusOK {

		// If we get 451 even after accepting the EULA, then the user would manually
		// connect to the webpage and accept the EULA as its pivotal legal policy
		// to accept the EULA if you are downloading the product for the first time.
		if statusCode == 451 {

			// Print to why we got this error
			fmt.Println("\n\x1b[33;1mReceived Status code 451 when access the API, this means as per the pivotal " +
				"legal policy if you are attempting to download the product for the first time you are requested to " +
				"to manually accept the end user license agreement (only one time). please connect " +
				"to PivNet and accept the end user license agreement and then try again, as this step cannot be avoided. Click on the link " +
				"mentioned below to redirect you to website to accept the EULA \x1b[0m")

			// Read the error text and store it
			bodyText, _ := ioutil.ReadAll(body)
			defer body.Close()
			var f interface{}
			_ = json.Unmarshal(bodyText, &f)
			m := f.(map[string]interface{})

			// Obtain the URL that the user can access to accept the EULA
			for k, v := range m {
				if k == "message" {
					fmt.Println("\n\x1b[35;1m" + v.(string) + "\x1b[0m\n")
				}
			}

		}
		return errors.New("API ERROR: HTTP Status code expected (" + strconv.Itoa(http.StatusOK) + ") / received (" + strconv.Itoa(statusCode) +  "): URL (" + url + ")")
	}
	return nil
}

// Get the Json from the provided URL or download the file
// if requested.
func (r *ResponseBody) GetApi(method string, urlLink string) ([]byte, error) {

	// Get the request
	req, err := http.NewRequest(method, urlLink, nil)

	// Add Header to the Http Request
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Adding Token header
	// Initial kickoff will not have any valid token
	// so we will request for the token
	if r.Token.ReceivedToken == "" {
		req.Header.Set("Authorization", "Bearer ")
		err := r.getToken(r.Env.Download.ApiToken)
		if err != nil || r.Token.RefreshToken == "" {
			return []byte(""), fmt.Errorf("authentication failure: %s", err.Error())
		}
		return []byte(""), nil
	} else { // hey we have the access token, so lets move on.
		req.Header.Set("Authorization", "Bearer " + r.Token.ReceivedToken)
	}

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}

	// If the status code is not 200, then error out
	err = r.ResponseChecker(resp.StatusCode, resp.Body, urlLink)
	if err != nil {
		return []byte(""), err
	}

	// Close the body once its done
	defer resp.Body.Close()

	// If its to download the software then download it
	if r.download {
		// The Size of the file
		size := r.UserRequest.ProductFileSize

		// Fully qualifies path
		r.UserRequest.ProductFileName = r.Env.Download.DownloadDir + r.UserRequest.ProductFileName

		// Create th file
		out, err := os.Create(r.UserRequest.ProductFileName)
		if err != nil {return []byte(""), err}

		// Initalize progress bar
		done := make(chan int64)
		go PrintDownloadPercent(done, r.UserRequest.ProductFileName, int64(size))
		defer out.Close()

		// Start Downloading
		n, err := io.Copy(out, resp.Body)
		if err != nil {return []byte(""), err}
		done <- n
	}

	// Read the json
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return bodyText, nil
}

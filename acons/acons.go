/*
	autils / acons - generates a URL to open the AWS Console
	Copyright (C) 2019 cloudninja.cloud

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type SigninTokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

func main() {

	var profile string
	var sessionName string
	var defaultSessionName string
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(CommandLine.Output(), "autils / acons - generates a URL to open the AWS Console\n© 2019 cloudninja.cloud\nSee https://ninja.cloudninja.cloud/utils/autils/acons/ for more details.\n\n")
		fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	username, err := user.Current()
	if err != nil {
		defaultSessionName = "acons"
	} else {
		defaultSessionName = username.Username
	}

	flag.StringVar(&profile, "p", "default", "Identifies the profile to use to generate the URL")
	flag.StringVar(&sessionName, "s", defaultSessionName, "Specifies the session name")
	flag.Parse()

	svc := sts.New(session.New())

	policy := aws.String("{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}")

	input := &sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(3600),
		Name:            &sessionName,
		Policy:          policy,
	}

	result, err := svc.GetFederationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(sts.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case sts.ErrCodePackedPolicyTooLargeException:
				fmt.Println(sts.ErrCodePackedPolicyTooLargeException, aerr.Error())
			case sts.ErrCodeRegionDisabledException:
				fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	accessKeyId := *result.Credentials.AccessKeyId
	secretAccessKey := *result.Credentials.SecretAccessKey
	sessionToken := *result.Credentials.SessionToken

	// Create the sign-in token using temporary credentials,
	// including the access key ID,  secret access key, and security token.

	requestParams := "{"
	requestParams += "\"sessionId\":\"" + accessKeyId + "\","
	requestParams += "\"sessionKey\":\"" + secretAccessKey + "\","
	requestParams += "\"sessionToken\":\"" + sessionToken + "\""
	requestParams += "}"

	request_parameters := "?Action=getSigninToken"
	request_parameters += "&DurationSeconds=3600"
	request_parameters += "&SessionType=json"
	request_parameters += "&Session=" + url.QueryEscape(requestParams)

	request_url := "https://signin.aws.amazon.com/federation" + request_parameters

	response, err := http.Get(request_url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}

		var s = new(SigninTokenResponse)
		err = json.Unmarshal(contents, &s)
		if err != nil {
			fmt.Println("whoops:", err)
		}

		consoleURL := "https://console.aws.amazon.com/"
		signInURL := "https://signin.aws.amazon.com/federation"

		signinTokenParameter := "&SigninToken=" + url.QueryEscape(s.SigninToken)
		destinationParameter := "&Destination=" + url.QueryEscape(consoleURL)

		loginURL := signInURL + "?Action=login" +
			signinTokenParameter + destinationParameter

		fmt.Println(loginURL)
	}
}

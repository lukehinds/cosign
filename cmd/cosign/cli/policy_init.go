package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/mail"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type MainPolicyStruct struct {
	BodyStruct   BodyStruct `json:"body"`
	SignedStruct Signed     `json:"signed"`
}

type BodyStruct struct {
	Maintainers       []string  `json:"maintainers"`
	RegistryNamespace string    `json:"registry_namespace"`
	Threshold         int       `json:"threshold"`
	Expires           time.Time `json:"expires"`
}

type SignedStruct struct {
	Email      string `json:"email"`
	FulcioCert string `json:"fuclio_cert"`
	Signature  string `json:"signature"`
}

type Signed []*SignedStruct

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func PolicyInit() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("cosign policy-init", flag.ExitOnError)

		nameSpace   = flagset.String("ns", "", "The registry namespace")
		mainTainers = flagset.String("maintainers", "", "Comma separated list of maintainers")
		threshHold  = flagset.Int("threshold", 2, "Threshold")
	)

	return &ffcli.Command{
		Name:       "policy-init",
		ShortUsage: "generate a new keyless policy",
		ShortHelp:  "policy-init is used to generate a root.json policy for keyless signing delegation",
		LongHelp: `policy-init is used to generate a root.json policy
for keyless signing delegation. This is used to establish a policy for a registry namespace,
a signing threshold and a list of maintainers who can sign over the body section.

EXAMPLES
  # extract public key from private key to a specified out file.
  cosign policy-init -ns <project_namespace> --maintainers {email_addresses} --threshold <int> --expires <int>(days)`,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {

			var emailList []string

			// Process the list of maintainers by
			// 1. Ensure each entry is a correctly formatted email address
			// 2. If 1 is true, then remove surplus whitespace (caused by gaps between commas)
			emails := strings.Split(*mainTainers, ",")
			for _, email := range emails {
				if !validEmail(email) {
					panic(fmt.Sprintf("Invalid email format: %s", email))
				} else {
					emailList = append(emailList, email)
					// strip out whitespace if there is any
					for i := range emailList {
						emailList[i] = strings.TrimSpace(emailList[i])
					}
				}
			}

			// The body constitutes the main body section of the policy
			// These kv's contain security guarantees
			body := BodyStruct{
				Maintainers:       emailList,
				RegistryNamespace: *nameSpace,
				Threshold:         *threshHold,
				Expires:           time.Now(),
			}

			// Signed is empty on first initialization of the policy
			// We signed over this as maintainers
			signed := Signed{}

			// Construct the complete policy
			policyJSON := MainPolicyStruct{
				BodyStruct:   body,
				SignedStruct: signed,
			}

			byteArray, err := json.MarshalIndent(policyJSON, "", "  ")

			if err != nil {
				return errors.Wrapf(err, "failed to marshal policy json")
			}

			fmt.Println(string(byteArray))
			err = ioutil.WriteFile("root.json", byteArray, 0600)
			if err != nil {
				return errors.Wrapf(err, "error writing to root.json")
			}
			return nil
		},
	}
}

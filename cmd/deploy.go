package cmd

import (
	"fmt"
	"github.com/jakeslee/dogecloud-cli/pkg/doge"
	"github.com/spf13/cobra"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [cert path] [cert key path]",
	Short: "Deploy cert to DogeCloud",
	Long:  `Deploy cert to DogeCloud`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cert, key := args[0], args[1]
		run(cert, key)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func run(cert, key string) {
	cert, err := readFile(cert)
	if err != nil {
		log.Fatal(err)
	}

	key, err = readFile(key)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fetching domains...")
	domains, err := doge.ListDomains()
	if err != nil {
		log.Fatal(err)
	}

	if len(domains) == 0 {
		log.Fatalf("no online domains, exit")
	}

	log.Printf("uploading cert and key...")
	certId, err := doge.UploadCert(fmt.Sprintf("%d", rand.Intn(100000000)), cert, key)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("new cert id: %d", certId)

	log.Printf("deploying cert to domains: %s", strings.Join(domains, ","))
	doge.DomainCertDeploy(certId, domains)

	log.Printf("cleaning unused certs...")
	for _, certObj := range doge.ListCerts() {
		if certObj.Count == 0 {
			log.Printf("deleting cert %d[%s], SAN: %s", certObj.Id, certObj.Name, strings.Join(certObj.Info.SAN, ","))
			doge.DeleteCert(certObj.Id)
		}
	}

	_, _ = doge.ListDomains()
}

func readFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

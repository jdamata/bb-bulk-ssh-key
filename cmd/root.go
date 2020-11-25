package cmd

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/ssh"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	bbBaseUri = "https://api.bitbucket.org/2.0"
	rootCmd   = &cobra.Command{
		Use:   "bb-bulk-ssh-key",
		Short: "Upload multiple ssh keys to your bitbucket account",
		Long:  "bb-bulk-ssh-key is a CLI that allows uploading multiple ssh keys to your bitbucket account.",
		Run:   main,
	}
)

func Execute(version string) error {
	rootCmd.Version = version
	rootCmd.PersistentFlags().StringP("username", "u", "", "Bitbucket username")
	rootCmd.MarkFlagRequired("username")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Bitbucket app password")
	rootCmd.MarkFlagRequired("password")
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "Directory to search for ssh keys")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	return rootCmd.Execute()
}

func main(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Get bitbucket user id
	getUserUUID()
	log.Fatalf("yo")
	directory := viper.GetString("directory")

	// Iterate through all files in directory and
	items, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatalf("Failed to get directory contents: %v", err)
	}

	for _, item := range items {
		// Skip if item is a directory
		if item.IsDir() {
			continue
		}

		//Read in file contents
		fileName := item.Name()
		out, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Warnf("Failed to read file %v", fileName)
			continue
		}

		// Check if file is a valid ssh key
		key, err := ssh.ParsePublicKey(out)
		if err != nil {
			log.Infof("File is not a valid ssh key: %v", fileName)
		}

		log.Info(key)
		// curl -X POST -H "Content-Type: application/json" -d '{"key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY user@myhost"}' https://api.bitbucket.org/2.0/users/{ed08f5e1-605b-4f4a-aee4-6c97628a673e}/ssh-keys

	}
}

func uploadSSHKey() {}

func getUserUUID() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", bbBaseUri+"/user", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(viper.GetString("username"), viper.GetString("password"))
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get user UUID: %v", err)
	}
	log.Info(resp)
	return ""
}

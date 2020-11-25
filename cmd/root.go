package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/ssh"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	client    = &http.Client{}
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
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
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
	uuid := getUserUUID()
	directory := viper.GetString("directory")

	// Iterate through all files in directory and
	items, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatalf("Failed to get directory contents. ERROR: %v", err)
	}

	for _, item := range items {
		// Skip if item is a directory
		if item.IsDir() {
			continue
		}

		if item.Name() == "known_hosts" {
			continue
		}

		//Read in file contents
		fileName := directory + item.Name()
		out, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Warnf("Failed to read %v. ERROR: %v", fileName, err)
			continue
		}

		// Check if file is a valid ssh key
		_, _, _, _, err = ssh.ParseAuthorizedKey(out)
		if err != nil {
			log.Infof("%v is not a valid ssh key. ERROR: %v", fileName, err)
			continue
		}

		body := map[string]string{"key": string(out)}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Warnf("Failed to marshal ssh key into json body, %v. ERROR: %v", fileName, err)
		}

		req, err := http.NewRequest("POST", bbBaseUri+"/users/"+uuid+"/ssh-keys", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Warnf("Failed to prepare http post request, %v. ERROR: %v", fileName, err)
		}

		req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Warnf("Failed to upload ssh key, %v. ERROR: %v", fileName, err)
		}
		defer resp.Body.Close()

		//out, err = ioutil.ReadAll(resp.Body)
		//var result map[string]interface{}
		//json.Unmarshal([]byte(out), &result)
		//log.Info(result)
	}
}

func uploadSSHKey() {}

func getUserUUID() string {
	req, err := http.NewRequest("GET", bbBaseUri+"/user", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
	resp, err := client.Do(req)
	if err != nil {

		log.Fatalf("Failed to get user UUID. ERROR: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	return result["uuid"].(string)
}

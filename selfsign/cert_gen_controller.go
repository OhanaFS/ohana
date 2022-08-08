package selfsign

import (
	"errors"
	"fmt"
	"github.com/OhanaFS/ohana/config"
	"gopkg.in/yaml.v3"
	"os"
)

// ProcessFlags processes the flags provided by the user.
func ProcessFlags(flagsConfig *config.FlagsConfig, debugMode bool) error {

	if *flagsConfig.GenCA {
		fmt.Println("Generating CA certs")

		err := GenCA(*flagsConfig.GenCAPath, debugMode)
		if err != nil {
			return err
		}

		if !*flagsConfig.GenCerts { // exit if gencerts is not set
			return nil
		}
	}
	if *flagsConfig.GenCerts {
		fmt.Println("Generating certs")

		if *flagsConfig.AllHosts == "" {
			fmt.Println("No hosts file found. See certhosts.example.yaml for example file to provide system")
			return nil
		}

		// Process all hosts

		hostFile, err := os.Open(*flagsConfig.AllHosts)
		if err != nil {
			return errors.New("Unable to open hosts file. See certhosts.example.yaml for example file to provide system")
		}
		defer hostFile.Close()

		var hosts Hosts
		decoder := yaml.NewDecoder(hostFile)
		if err = decoder.Decode(&hosts); err != nil {
			return err
		}

		// verify that numOfCerts is more than 0
		if *flagsConfig.NumOfCerts <= 0 {
			fmt.Println("Number of certs must be more than 0")
			return nil
		}

		for i := 0; i < *flagsConfig.NumOfCerts; i++ {
			err = GenCerts(*flagsConfig.CsrPath, *flagsConfig.CertPath, *flagsConfig.PkPath, *flagsConfig.GenCertsPath, hosts.Hosts)
			if err != nil {
				return err
			}
		}

		return nil
	}
	return nil

}

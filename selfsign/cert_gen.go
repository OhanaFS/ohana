package selfsign

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cfssl/api/generator"
	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type outputFile struct {
	Filename string
	Contents string
	IsBinary bool
	Perms    os.FileMode
}

type Hosts struct {
	Hosts []string `yaml:"hosts"`
}

// GenCA generates a new CA certificate and private key.
// It returns the certificate and private key in PEM format and the CSR in JSON format.
func GenCA(pathName string) error {

	if pathName == "" {
		pathName = "main"
	}

	// asking user for csr info
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your country: ")
	country, _ := reader.ReadString('\n')
	fmt.Print("Enter your state: ")
	state, _ := reader.ReadString('\n')
	fmt.Print("Enter your locality: ")
	locality, _ := reader.ReadString('\n')
	fmt.Print("Enter your organization: ")
	organization, _ := reader.ReadString('\n')
	fmt.Print("Enter your organization unit: ")
	organizationalUnit, _ := reader.ReadString('\n')

	req := csr.CertificateRequest{
		KeyRequest: csr.NewKeyRequest(),
		Names: []csr.Name{
			{
				C:  country,
				ST: state,
				L:  locality,
				O:  organization,
				OU: organizationalUnit,
			},
		},
		Hosts: []string{"*.local", "localhost"},
	}

	req.KeyRequest.A = "rsa"
	req.KeyRequest.S = 2048

	// Create a new CA certificate and private key
	var key, csrPEM, cert []byte
	cert, csrPEM, key, err := initca.New(&req)
	if err != nil {
		return err
	}

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cli.PrintCert(key, csrPEM, cert)

	outChannel := make(chan string)
	// copy the output to outChannel
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			log.Error(err)
		}
		outChannel <- buf.String()
	}()

	err = w.Close()
	if err != nil {
		return err
	}
	os.Stdout = old // restoring the real stdout
	outputString := <-outChannel
	err = JSONCertWriter(outputString, pathName, true)
	if err != nil {
		return err
	}

	// json marshal the csr
	// write to file
	jsonString, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		panic(err)
	}
	// write to file
	return writeFile(pathName+"_csr.json", string(jsonString), 0644)

}

// GenCerts generates a new certificate and private key for the nodes/servers to use.
// It returns the certificate and private key in PEM format and the CSR in JSON format.
func GenCerts(csrPath, certPath, pkPath, output string, hosts []string) error {

	if output == "" {
		output = "output"
	}

	csrJSONFileBytes, err := cli.ReadStdin(csrPath) // "main_csr.json"
	if err != nil {
		return err
	}

	req := csr.CertificateRequest{
		KeyRequest: csr.NewKeyRequest(),
	}
	err = json.Unmarshal(csrJSONFileBytes, &req)

	if len(hosts) > 0 {
		req.Hosts = append(req.Hosts, hosts...)
	}

	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		key = nil
		return err
	}

	c := cli.Config{
		CAFile:    certPath, // This is the CA certificate "main_GLOBAL_CERTIFICATE.pem"
		CAKeyFile: pkPath,   // This is the CA private key "main_PRIVATE_KEY.pem"
		CFG:       nil,
	}

	s, err := sign.SignerFromConfig(c)
	if err != nil {
		return err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
		Label:   c.Label,
	}

	if c.CRL != "" {
		signReq.CRLOverride = c.CRL
	}
	cert, err = s.Sign(signReq)
	if err != nil {
		return err
	}

	if len(signReq.Hosts) == 0 && len(req.Hosts) == 0 {
		log.Warning(generator.CSRNoHostMessage)
	}

	// Capturing stdout and piping it to a string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cli.PrintCert(key, csrBytes, cert)

	outChannel := make(chan string)
	// copy the output to outChannel
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			log.Error(err)
		}
		outChannel <- buf.String()
	}()

	err = w.Close()
	if err != nil {
		return err
	}
	os.Stdout = old // restoring the real stdout
	outputString := <-outChannel
	return JSONCertWriter(outputString, output, false)

}

func writeFile(filespec, contents string, perms os.FileMode) error {

	// if file exists, add a number to the end of the file name
	var err error
	i := 0
	// split filename and extension
	file := strings.Split(filespec, ".")[0]
	ext := strings.Split(filespec, ".")[1]
	err = nil

	for err == nil {
		_, err = os.Stat(filespec)
		if err == nil {
			i++
			filespec = file + "_" + strconv.Itoa(i) + "." + ext
		}
	}

	if _, err := os.Stat(filespec); err == nil {
		filespec = filespec + ".new"

	}
	err = ioutil.WriteFile(filespec, []byte(contents), perms)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%v\n", err)
		if err != nil {
			return err
		}
	}
	return err
}

// JSONCertWriter helps converts the output from cfssl to a json file
// in an array that writeFile can write with
func JSONCertWriter(in string, baseName string, caMode bool) error {
	// Copied from cfssljson.go from cfssl package by cloudflare

	var input = map[string]interface{}{}
	var outs []outputFile
	var cert, key string

	if baseName == "" {
		baseName = "Ohana"
	}

	err := json.Unmarshal([]byte(in), &input)
	if err != nil {
		return err
	}

	if contents, ok := input["cert"]; ok {
		cert = contents.(string)
	} else if contents, ok = input["certificate"]; ok {
		cert = contents.(string)
	}
	if cert != "" {
		var filename string
		if caMode {
			filename = baseName + "_GLOBAL_CERTIFICATE.pem"
		} else {
			filename = baseName + "_cert.pem"
		}
		outs = append(outs, outputFile{
			Filename: filename,
			Contents: cert,
			Perms:    0664,
		})
	}

	if contents, ok := input["key"]; ok {
		key = contents.(string)
	} else if contents, ok = input["private_key"]; ok {
		key = contents.(string)
	}
	if key != "" {
		var filename string
		if caMode {
			filename = baseName + "_PRIVATE_KEY.pem"
		} else {
			filename = baseName + "_key.pem"
		}
		outs = append(outs, outputFile{
			Filename: filename,
			Contents: key,
			Perms:    0600,
		})
	}

	/*
		// Commented out because we don't use it.
		var csr string
		if contents, ok := input["csr"]; ok {
			csr = contents.(string)
		} else if contents, ok = input["certificate_request"]; ok {
			csr = contents.(string)
		}
		if csr != "" {
			outs = append(outs, outputFile{
				Filename: baseName + ".csr",
				Contents: csr,
				Perms:    0644,
			})
		}

	*/

	for _, e := range outs {
		err := writeFile(e.Filename, e.Contents, e.Perms)
		if err != nil {
			return err
		}
	}

	return nil
}

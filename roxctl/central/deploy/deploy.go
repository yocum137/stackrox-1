package deploy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	cflog "github.com/cloudflare/cfssl/log"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/docker"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/zip"
	"github.com/stackrox/rox/roxctl/central/deploy/renderer"
)

func init() {
	// The cfssl library prints logs at Info level when it processes a
	// Certificate Signing Request (CSR) or issues a new certificate.
	// These logs do not help the user understand anything, so here
	// we adjust the log level to exclude them.
	cflog.Level = cflog.LevelWarning
}

var (
	generatedMonitoringPassword = renderer.CreatePassword()
)

func generateJWTSigningKey(fileMap map[string][]byte) error {
	// Generate the private key that we will use to sign JWTs for API keys.
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("couldn't generate private key: %s", err)
	}
	fileMap["jwt-key.der"] = x509.MarshalPKCS1PrivateKey(privateKey)
	return nil
}

func generateMTLSFiles(fileMap map[string][]byte) (cert, key []byte, err error) {
	// Add MTLS files
	req := csr.CertificateRequest{
		CN:         "StackRox Certificate Authority",
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	cert, _, key, err = initca.New(&req)
	if err != nil {
		err = fmt.Errorf("could not generate keypair: %s", err)
		return
	}
	fileMap["ca.pem"] = cert
	fileMap["ca-key.pem"] = key
	return
}

func generateMonitoringFiles(fileMap map[string][]byte, caCert, caKey []byte) error {
	monitoringCert, err := mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_DB_SERVICE, Identifier: "Monitoring DB"},
		caCert, caKey)
	if err != nil {
		return err
	}
	fileMap["monitoring-db-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-db-key.pem"] = monitoringCert.KeyPEM

	monitoringCert, err = mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_UI_SERVICE, Identifier: "Monitoring UI"},
		caCert, caKey)
	if err != nil {
		return err
	}

	fileMap["monitoring-ui-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-ui-key.pem"] = monitoringCert.KeyPEM

	monitoringCert, err = mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_CLIENT_SERVICE, Identifier: "Monitoring Client"},
		caCert, caKey)
	if err != nil {
		return err
	}

	fileMap["monitoring-client-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-client-key.pem"] = monitoringCert.KeyPEM

	return nil
}

func outputZip(config renderer.Config) error {
	fmt.Fprint(os.Stderr, "Generating deployment bundle... ")

	wrapper := zip.NewWrapper()

	d, ok := renderer.Deployers[config.ClusterType]
	if !ok {
		return fmt.Errorf("undefined cluster deployment generator: %s", config.ClusterType)
	}

	config.SecretsByteMap = make(map[string][]byte)
	if err := generateJWTSigningKey(config.SecretsByteMap); err != nil {
		return err
	}

	if features.HtpasswdAuth.Enabled() {
		if config.Environment == nil {
			config.Environment = make(map[string]string)
		}
		config.Environment[features.HtpasswdAuth.EnvVar()] = "true"

		htpasswd, err := renderer.GenerateHtpasswd(&config)
		if err != nil {
			return err
		}

		config.SecretsByteMap["htpasswd"] = htpasswd
		wrapper.AddFiles(zip.NewFile("password", []byte(config.Password+"\n"), zip.Sensitive))
	}

	cert, key, err := generateMTLSFiles(config.SecretsByteMap)
	if err != nil {
		return err
	}

	if config.K8sConfig != nil && config.K8sConfig.MonitoringType.OnPrem() {
		generateMonitoringFiles(config.SecretsByteMap, cert, key)

		if config.K8sConfig.MonitoringPassword == "" {
			config.K8sConfig.MonitoringPassword = renderer.CreatePassword()
			config.K8sConfig.MonitoringPasswordAuto = true
		}

		config.SecretsByteMap["monitoring-password"] = []byte(config.K8sConfig.MonitoringPassword)
		wrapper.AddFiles(zip.NewFile("monitoring/password", []byte(config.K8sConfig.MonitoringPassword+"\n"), zip.Sensitive))
	}

	config.SecretsBase64Map = make(map[string]string)
	for k, v := range config.SecretsByteMap {
		config.SecretsBase64Map[k] = base64.StdEncoding.EncodeToString(v)
	}

	files, err := d.Render(config)
	if err != nil {
		return fmt.Errorf("could not render files: %s", err)
	}
	wrapper.AddFiles(files...)

	var outputPath string
	if docker.IsContainerized() {
		bytes, err := wrapper.Zip()
		if err != nil {
			return fmt.Errorf("error generating zip file: %v", err)
		}
		_, err = os.Stdout.Write(bytes)
		if err != nil {
			return fmt.Errorf("couldn't write zip file: %v", err)
		}
	} else {
		var err error
		outputPath, err = wrapper.Directory(config.OutputDir)
		if err != nil {
			return fmt.Errorf("error generating directory for Central output: %v", err)
		}
	}

	fmt.Fprintln(os.Stderr, "Done!")
	fmt.Fprintln(os.Stderr)

	if outputPath != "" {
		fmt.Fprintf(os.Stderr, "Wrote central bundle to %q\n", outputPath)
		fmt.Fprintln(os.Stderr)
	}

	config.WriteInstructions(os.Stderr)
	return nil
}

func interactive() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Interactive runs the CLI in interactive mode with user prompts",
		RunE: func(c *cobra.Command, args []string) error {
			c = Command()
			c.SilenceUsage = true
			return runInteractive(c)
		},
		SilenceUsage: true,
	}
}

// Command defines the deploy command tree
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "generate",
		Short: "Generate creates the required YAML files to deploy StackRox Central.",
		Long:  "Generate creates the required YAML files to deploy StackRox Central.",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}
	if features.HtpasswdAuth.Enabled() {
		c.PersistentFlags().StringVar(&cfg.Password, "password", "", "administrator password (default: autogenerated)")
	}

	c.PersistentFlags().Var(&featureValue{&cfg.Features}, "flags", "Feature flags to enable")
	c.PersistentFlags().MarkHidden("flags")
	c.AddCommand(interactive())

	c.AddCommand(k8s())
	c.AddCommand(openshift())
	return c
}

func runInteractive(cmd *cobra.Command) error {
	// Overwrite os.Args because cobra uses them
	os.Args = walkTree(cmd)
	return cmd.Execute()
}

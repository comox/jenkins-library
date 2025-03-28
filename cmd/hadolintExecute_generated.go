// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/gcp"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type hadolintExecuteOptions struct {
	ConfigurationURL          string   `json:"configurationUrl,omitempty"`
	ConfigurationUsername     string   `json:"configurationUsername,omitempty"`
	ConfigurationPassword     string   `json:"configurationPassword,omitempty"`
	DockerFile                string   `json:"dockerFile,omitempty"`
	ConfigurationFile         string   `json:"configurationFile,omitempty"`
	ReportFile                string   `json:"reportFile,omitempty"`
	CustomTLSCertificateLinks []string `json:"customTlsCertificateLinks,omitempty"`
}

// HadolintExecuteCommand Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.
func HadolintExecuteCommand() *cobra.Command {
	const STEP_NAME = "hadolintExecute"

	metadata := hadolintExecuteMetadata()
	var stepConfig hadolintExecuteOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createHadolintExecuteCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.",
		Long: `Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.
The linter is parsing the Dockerfile into an abstract syntax tree (AST) and performs rules on top of the AST.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, err := os.Getwd()
			if err != nil {
				return err
			}
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err = PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.ConfigurationUsername)
			log.RegisterSecret(stepConfig.ConfigurationPassword)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 || len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			vaultClient := config.GlobalVaultClient()
			if vaultClient != nil {
				defer vaultClient.MustRevokeToken()
			}

			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.LogStepTelemetryData()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if GeneralConfig.HookConfig.GCPPubSubConfig.Enabled {
					err := gcp.NewGcpPubsubClient(
						vaultClient,
						GeneralConfig.HookConfig.GCPPubSubConfig.ProjectNumber,
						GeneralConfig.HookConfig.GCPPubSubConfig.IdentityPool,
						GeneralConfig.HookConfig.GCPPubSubConfig.IdentityProvider,
						GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.OIDCConfig.RoleID,
					).Publish(GeneralConfig.HookConfig.GCPPubSubConfig.Topic, telemetryClient.GetDataBytes())
					if err != nil {
						log.Entry().WithError(err).Warn("event publish failed")
					}
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(STEP_NAME)
			hadolintExecute(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addHadolintExecuteFlags(createHadolintExecuteCmd, &stepConfig)
	return createHadolintExecuteCmd
}

func addHadolintExecuteFlags(cmd *cobra.Command, stepConfig *hadolintExecuteOptions) {
	cmd.Flags().StringVar(&stepConfig.ConfigurationURL, "configurationUrl", os.Getenv("PIPER_configurationUrl"), "URL pointing to the .hadolint.yaml exclude configuration to be used for linting. Also have a look at `configurationFile` which could avoid central configuration download in case the file is part of your repository.")
	cmd.Flags().StringVar(&stepConfig.ConfigurationUsername, "configurationUsername", os.Getenv("PIPER_configurationUsername"), "The username to authenticate")
	cmd.Flags().StringVar(&stepConfig.ConfigurationPassword, "configurationPassword", os.Getenv("PIPER_configurationPassword"), "The password to authenticate")
	cmd.Flags().StringVar(&stepConfig.DockerFile, "dockerFile", `./Dockerfile`, "Dockerfile to be used for the assessment.")
	cmd.Flags().StringVar(&stepConfig.ConfigurationFile, "configurationFile", `.hadolint.yaml`, "Name of the configuration file used locally within the step. If a file with this name is detected as part of your repo downloading the central configuration via `configurationUrl` will be skipped. If you change the file's name make sure your stashing configuration also reflects this.")
	cmd.Flags().StringVar(&stepConfig.ReportFile, "reportFile", `hadolint.xml`, "Name of the result file used locally within the step.")
	cmd.Flags().StringSliceVar(&stepConfig.CustomTLSCertificateLinks, "customTlsCertificateLinks", []string{}, "List of download links to custom TLS certificates. This is required to ensure trusted connections between Piper and the system where the configuration file is to be downloaded from.")

}

// retrieve step metadata
func hadolintExecuteMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "hadolintExecute",
			Aliases:     []config.Alias{},
			Description: "Executes the Haskell Dockerfile Linter which is a smarter Dockerfile linter that helps you build [best practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "configurationCredentialsId", Description: "Jenkins 'Username with password' credentials ID containing username/password for access to your remote configuration file.", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "configurationUrl",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_configurationUrl"),
					},
					{
						Name: "configurationUsername",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "configurationCredentialsId",
								Param: "username",
								Type:  "secret",
							},

							{
								Name:    "hadolintConfigSecretName",
								Type:    "vaultSecret",
								Default: "hadolintConfig",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "username"}},
						Default:   os.Getenv("PIPER_configurationUsername"),
					},
					{
						Name: "configurationPassword",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "configurationCredentialsId",
								Param: "password",
								Type:  "secret",
							},

							{
								Name:    "hadolintConfigSecretName",
								Type:    "vaultSecret",
								Default: "hadolintConfig",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "password"}},
						Default:   os.Getenv("PIPER_configurationPassword"),
					},
					{
						Name:        "dockerFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "dockerfile"}},
						Default:     `./Dockerfile`,
					},
					{
						Name:        "configurationFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `.hadolint.yaml`,
					},
					{
						Name:        "reportFile",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `hadolint.xml`,
					},
					{
						Name:        "customTlsCertificateLinks",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     []string{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "hadolint", Image: "hadolint/hadolint:latest-alpine"},
			},
		},
	}
	return theMetaData
}

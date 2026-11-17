// Package cmd provides the Cobra command-line interface for kube-cost-lens.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dawidbera/kube-cost-lens/internal/finops"
	"github.com/dawidbera/kube-cost-lens/internal/k8s"
	"github.com/dawidbera/kube-cost-lens/internal/tui"
	"github.com/spf13/cobra"
)

var (
	// CLI flag variables
	namespaceFlag       string
	demoFlag            bool
	ratesPathFlag       string
	cpuRateFlag         float64
	ramRateFlag         float64
	refreshIntervalFlag time.Duration
	kubeconfigFlag      string
	contextFlag         string
)

// rootCmd defines the base command when called without subcommands.
var rootCmd = &cobra.Command{
	Use:   "kube-cost-lens",
	Short: "Interactive FinOps Terminal UI for Kubernetes resource & cost optimization",
	Long: `kube-cost-lens is an interactive, single-binary Terminal User Interface (TUI)
for Kubernetes SREs, FinOps practitioners, and platform engineers.

It inspects cluster workloads, compares CPU/Memory requests vs. actual consumption,
calculates idle waste and financial overruns in real-time, and generates
optimized resource patch manifests with 1 click.`,
	RunE: runRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Target specific Kubernetes namespace (default: all namespaces)")
	rootCmd.Flags().BoolVar(&demoFlag, "demo", false, "Run in synthetic offline demo mode with simulated cluster workloads")
	rootCmd.Flags().StringVar(&ratesPathFlag, "rates", "", "Path to custom JSON pricing rate card file")
	rootCmd.Flags().Float64Var(&cpuRateFlag, "cpu-rate", 0, "Hourly cost per vCPU core (overrides default/file pricing)")
	rootCmd.Flags().Float64Var(&ramRateFlag, "ram-rate", 0, "Hourly cost per GB of memory (overrides default/file pricing)")
	rootCmd.Flags().DurationVar(&refreshIntervalFlag, "refresh-interval", 10*time.Second, "Background cluster scrape interval")
	rootCmd.Flags().StringVar(&kubeconfigFlag, "kubeconfig", "", "Path to the kubeconfig file (overrides KUBECONFIG env)")
	rootCmd.Flags().StringVar(&contextFlag, "context", "", "Override default kubeconfig context name")
}

// runRoot executes the main application lifecycle.
func runRoot(cmd *cobra.Command, _ []string) error {
	// 1. Resolve pricing configuration
	pricingCard := finops.DefaultPricingCard()
	if ratesPathFlag != "" {
		loaded, err := finops.LoadPricingCard(ratesPathFlag)
		if err != nil {
			return fmt.Errorf("error loading pricing file: %w", err)
		}
		pricingCard = loaded
	}

	if cpuRateFlag > 0 {
		pricingCard.CPURatePerHour = cpuRateFlag
	}
	if ramRateFlag > 0 {
		pricingCard.MemoryRatePerGBHour = ramRateFlag
	}

	// 2. Resolve cluster provider (Mock vs Live)
	var provider k8s.ClusterProvider
	if demoFlag {
		provider = k8s.NewMockClusterProvider()
	} else {
		liveOpts := k8s.ClientOptions{
			KubeconfigPath: kubeconfigFlag,
			Context:        contextFlag,
		}
		liveProvider, err := k8s.NewLiveClusterProvider(liveOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: unable to connect to live Kubernetes cluster: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to --demo mode. Pass --demo explicitly to silence this warning.\n\n")
			provider = k8s.NewMockClusterProvider()
		} else {
			provider = liveProvider
		}
	}
	defer func() {
		_ = provider.Close()
	}()

	// 3. Launch Bubbletea TUI program
	return tui.RunApp(tui.AppOptions{
		Provider:        provider,
		PricingCard:     pricingCard,
		NamespaceFilter: namespaceFlag,
		RefreshInterval: refreshIntervalFlag,
	})
}

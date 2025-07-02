package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thegostev/go-kubernetes-controllers/pkg/controller"
	"github.com/thegostev/go-kubernetes-controllers/pkg/k8s"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	disableLeaderElection bool
	metricsPort           int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the controller manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
		mgr, err := ctrl.NewManager(k8s.NewConfigOrDie(), manager.Options{
			Scheme:                     k8s.NewScheme(),
			LeaderElection:             !disableLeaderElection,
			LeaderElectionID:           "go-k8s-ctrl-leader-election",
			LeaderElectionResourceLock: "leases",
		})
		if err != nil {
			return err
		}
		if err := controller.SetupDeploymentController(mgr); err != nil {
			return err
		}
		return mgr.Start(ctrl.SetupSignalHandler())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVar(&disableLeaderElection, "disable-leader-election", false, "Disable leader election for controller manager")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "The port the metric endpoint binds to")
}

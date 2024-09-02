package approve

import (
	"fmt"
	"io"

	"github.com/openshift-pipelines/manual-approval-gate/pkg/actions"
	cli "github.com/openshift-pipelines/manual-approval-gate/pkg/cli"
	"github.com/openshift-pipelines/manual-approval-gate/pkg/cli/flags"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	taskGroupResource = schema.GroupVersionResource{Group: "openshift-pipelines.org", Resource: "approvaltasks"}
)

func Command(p cli.Params) *cobra.Command {
	opts := &cli.Options{}
	c := &cobra.Command{
		Use:   "approve",
		Short: "Approve the approvaltask",
		Long:  `This command approves the approvaltask.`,
		Annotations: map[string]string{
			"commandType": "main",
		},
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: flags.PersistentPreRunE(p),
		RunE: func(cmd *cobra.Command, args []string) error {
			cs, err := p.Clients()
			if err != nil {
				return err
			}

			ns := p.Namespace()
			if opts.AllNamespaces {
				ns = ""
			}

			username, err := p.GetUserInfo()
			if err != nil {
				return err
			}

			message := opts.Message

			opts = &cli.Options{
				Name:      args[0],
				Namespace: ns,
				Input:     "approve",
				Username:  username,
				Message:   message,
			}

			if err := actions.Update(taskGroupResource, cs, opts); err != nil {
				return fmt.Errorf("failed to approve approvalTask from namespace %s: %v", ns, err)
			}

			res := fmt.Sprintf("ApprovalTask %s is approved in %s namespace\n", args[0], ns)
			io.WriteString(cmd.OutOrStdout(), res)

			return nil

		},
	}

	c.Flags().StringVarP(&opts.Message, "message", "m", "", "message while approving the approvalTask")

	flags.AddOptions(c)

	return c
}

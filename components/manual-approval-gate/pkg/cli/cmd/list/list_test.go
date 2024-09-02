package list

import (
	"fmt"
	"strings"
	"testing"

	"github.com/openshift-pipelines/manual-approval-gate/pkg/apis/approvaltask/v1alpha1"
	"github.com/openshift-pipelines/manual-approval-gate/pkg/test"
	cb "github.com/openshift-pipelines/manual-approval-gate/pkg/test/builder"
	testDynamic "github.com/openshift-pipelines/manual-approval-gate/pkg/test/dynamic"
	"github.com/spf13/cobra"
	"gotest.tools/v3/golden"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

func TestListApprovalTasks(t *testing.T) {
	approvaltasks := []*v1alpha1.ApprovalTask{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "at-1",
				Namespace: "foo",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "reject",
					},
					{
						Name:  "cli",
						Input: "pending",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				ApproversResponse: []v1alpha1.ApproverState{
					{
						Name:     "tekton",
						Response: "rejected",
					},
				},
				State: "rejected",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "at-2",
				Namespace: "foo",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "approve",
					},
					{
						Name:  "cli",
						Input: "approve",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				ApproversResponse: []v1alpha1.ApproverState{
					{
						Name:     "tekton",
						Response: "approve",
					},
					{
						Name:     "cli",
						Response: "approve",
					},
				},
				State: "approved",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "at-3",
				Namespace: "foo",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "pending",
					},
					{
						Name:  "cli",
						Input: "pending",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				State: "pending",
			},
		},
	}

	approvaltasksMultipleNs := []*v1alpha1.ApprovalTask{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "mango",
				Namespace: "test-1",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "reject",
					},
					{
						Name:  "cli",
						Input: "pending",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				ApproversResponse: []v1alpha1.ApproverState{
					{
						Name:     "tekton",
						Response: "rejected",
					},
				},
				State: "rejected",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apple",
				Namespace: "test-2",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "approve",
					},
					{
						Name:  "cli",
						Input: "approve",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				ApproversResponse: []v1alpha1.ApproverState{
					{
						Name:     "tekton",
						Response: "approve",
					},
					{
						Name:     "cli",
						Response: "approve",
					},
				},
				State: "approved",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "banana",
				Namespace: "test-3",
			},
			Spec: v1alpha1.ApprovalTaskSpec{
				Approvers: []v1alpha1.ApproverDetails{
					{
						Name:  "tekton",
						Input: "pending",
					},
					{
						Name:  "cli",
						Input: "pending",
					},
				},
				NumberOfApprovalsRequired: 2,
			},
			Status: v1alpha1.ApprovalTaskStatus{
				Approvers: []string{
					"tekton",
					"cli",
				},
				State: "pending",
			},
		},
	}

	ns := []*corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "namespace",
			},
		},
	}

	dc, err := testDynamic.Client(
		cb.UnstructuredV1alpha1(approvaltasks[0], "v1alpha1"),
		cb.UnstructuredV1alpha1(approvaltasks[1], "v1alpha1"),
		cb.UnstructuredV1alpha1(approvaltasks[2], "v1alpha1"),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	dc2, err := testDynamic.Client(
		cb.UnstructuredV1alpha1(approvaltasksMultipleNs[0], "v1alpha1"),
		cb.UnstructuredV1alpha1(approvaltasksMultipleNs[1], "v1alpha1"),
		cb.UnstructuredV1alpha1(approvaltasksMultipleNs[2], "v1alpha1"),
	)
	if err != nil {
		t.Errorf("unable to create dynamic client: %v", err)
	}

	tests := []struct {
		name      string
		command   *cobra.Command
		args      []string
		wantError bool
	}{
		{
			name:      "no approval tasks found",
			command:   command(t, approvaltasks, ns, dc),
			args:      []string{"list", "-n", "invalid"},
			wantError: true,
		},
		{
			name:      "all in namespace",
			command:   command(t, approvaltasks, ns, dc),
			args:      []string{"list", "-n", "foo"},
			wantError: false,
		},
		{
			name:      "in all namespaces",
			command:   command(t, approvaltasksMultipleNs, ns, dc2),
			args:      []string{"list", "--all-namespaces"},
			wantError: false,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			output, err := test.ExecuteCommand(td.command, td.args...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err != nil && !td.wantError {
				t.Errorf("Unexpected error: %v", err)
			}

			golden.Assert(t, output, strings.ReplaceAll(fmt.Sprintf("%s.golden", t.Name()), "/", "-"))
		})
	}

}

func command(t *testing.T, approvaltasks []*v1alpha1.ApprovalTask, ns []*corev1.Namespace, dc dynamic.Interface) *cobra.Command {
	cs, _ := test.SeedTestData(t, test.Data{Approvaltasks: approvaltasks, Namespaces: ns})
	p := &test.Params{ApprovalTask: cs.ApprovalTask, Kube: cs.Kube, Dynamic: dc}
	cs.ApprovalTask.Resources = cb.APIResourceList("v1alpha1", []string{"approvaltask"})

	return Command(p)
}

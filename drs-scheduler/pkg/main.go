package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "DQNPlugin"
	// ErrReason returned when node name doesn't match.
	ErrReason = "this node is not the result given by the DRL scheduler"
)

type DQNPlugin struct {
	handle framework.Handle
}

var _ framework.FilterPlugin = DQNPlugin{}

func (dp DQNPlugin) Name() string {
	return Name
}

func (dp DQNPlugin) Filter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	fmt.Printf("[INFO] Filtering pod: %v, node: %v\n", pod.Name, nodeInfo.Node().Name)
	node := nodeInfo.Node()
	if node == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}

	// get result from rl algorithm
	schedulerUrl := "http://172.31.234.111:1234/choose"
	urlValues := url.Values{}
	urlValues.Add("podname", pod.Name)
	resp, err := http.PostForm(schedulerUrl, urlValues)
	if err != nil {
		fmt.Printf("[ERROR] Get choose from %v failed\n", schedulerUrl)
		fmt.Printf("[INFO] DRL scheduling failed, all node will pass the filter\n\n")
		return framework.NewStatus(framework.Success, "")
	}
	body, _ := io.ReadAll(resp.Body)
	choose := string(body)
	fmt.Printf("[INFO] Get choose from %v successfully: %v\n", schedulerUrl, choose)

	if node.Name == choose {
		fmt.Printf("[INFO] Filter pod successfully: %v, node: %v\n\n", pod.Name, nodeInfo.Node().Name)
		return framework.NewStatus(framework.Success, "")
	} else {
		fmt.Printf("[ERROR] Filter pod failed: %v, node: %v\n\n", pod.Name, nodeInfo.Node().Name)
		return framework.NewStatus(framework.Error, ErrReason)
	}
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &DQNPlugin{handle: h}, nil
}

func main() {
	fmt.Printf("[INFO] Start DQN plugin\n")

	// create a new scheduler
	command := app.NewSchedulerCommand(
		app.WithPlugin(Name, New),
	)

	// run the scheduler
	if err := command.Execute(); err != nil {
		fmt.Printf("[ERROR] Run scheduler failed: %v\n", err)
		os.Exit(1)
	}
}

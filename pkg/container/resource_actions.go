package container

import (
	"github.com/openspacee/ospagent/pkg/container/resource"
	"github.com/openspacee/ospagent/pkg/kubernetes"
	"github.com/openspacee/ospagent/pkg/utils"
	"github.com/openspacee/ospagent/pkg/websocket"
)

const (
	LIST       = "list"
	GET        = "get"
	DELETE     = "delete"
	UPDATEYAML = "update_yaml"
	UPDATEOBJ  = "update_obj"
	EXEC       = "exec"
	STDIN      = "stdin"
	OPENLOG    = "openLog"
	CLOSELOG   = "closeLog"
)

type Handler func(interface{}) *utils.Response

type ActionHandler map[string]Handler

type ResourceActions struct {
	KubeClient            *kubernetes.KubeClient
	ResourceActionHandler map[string]ActionHandler
}

func NewResourceActions(kubeClient *kubernetes.KubeClient, sendResponse websocket.SendResponse) *ResourceActions {
	actionHandlers := make(map[string]ActionHandler)

	watch := resource.NewWatchResource(sendResponse)
	watchActions := ActionHandler{
		GET: watch.WatchAction,
	}
	actionHandlers["watch"] = watchActions

	pod := resource.NewPod(kubeClient, sendResponse, watch)
	podActions := ActionHandler{
		LIST:       pod.List,
		GET:        pod.Get,
		EXEC:       pod.Exec,
		STDIN:      pod.ExecStdIn,
		OPENLOG:    pod.OpenLog,
		CLOSELOG:   pod.CloseLog,
		DELETE:     pod.Delete,
		UPDATEYAML: pod.UpdateYaml,
	}
	actionHandlers["pod"] = podActions

	ns := resource.NewNamespace(kubeClient, sendResponse, watch)
	nsActions := ActionHandler{
		LIST: ns.List,
	}
	actionHandlers["namespace"] = nsActions

	node := resource.NewNode(kubeClient, sendResponse)
	nodeActions := ActionHandler{
		LIST: node.List,
	}
	actionHandlers["node"] = nodeActions

	event := resource.NewEvent(kubeClient, watch)
	eventActions := ActionHandler{
		LIST: event.List,
	}
	actionHandlers["event"] = eventActions

	deployment := resource.NewDeployment(kubeClient, watch)
	deploymentActions := ActionHandler{
		LIST:       deployment.List,
		GET:        deployment.Get,
		DELETE:     deployment.Delete,
		UPDATEYAML: deployment.UpdateYaml,
		UPDATEOBJ:  deployment.UpdateObj,
	}
	actionHandlers["deployment"] = deploymentActions

	statefulset := resource.NewStatefulSet(kubeClient, watch)
	statefulsetActions := ActionHandler{
		LIST:       statefulset.List,
		GET:        statefulset.Get,
		DELETE:     statefulset.Delete,
		UPDATEYAML: statefulset.UpdateYaml,
		UPDATEOBJ:  statefulset.UpdateObj,
	}
	actionHandlers["statefulset"] = statefulsetActions

	daemonset := resource.NewDaemonSet(kubeClient, watch)
	daemonsetActions := ActionHandler{
		LIST:       daemonset.List,
		GET:        daemonset.Get,
		DELETE:     daemonset.Delete,
		UPDATEYAML: daemonset.UpdateYaml,
		UPDATEOBJ:  daemonset.UpdateObj,
	}
	actionHandlers["daemonset"] = daemonsetActions

	job := resource.NewJob(kubeClient, watch)
	jobActions := ActionHandler{
		LIST:       job.List,
		GET:        job.Get,
		DELETE:     job.Delete,
		UPDATEYAML: job.UpdateYaml,
		UPDATEOBJ:  job.UpdateObj,
	}
	actionHandlers["job"] = jobActions

	cronjob := resource.NewCronJob(kubeClient, watch)
	cronjobActions := ActionHandler{
		LIST:       cronjob.List,
		GET:        cronjob.Get,
		DELETE:     cronjob.Delete,
		UPDATEYAML: cronjob.UpdateYaml,
		UPDATEOBJ:  cronjob.UpdateObj,
	}
	actionHandlers["cronjob"] = cronjobActions

	configMap := resource.NewConfigMap(kubeClient, sendResponse)
	configMapActions := ActionHandler{
		LIST: configMap.List,
		GET:  configMap.Get,
	}
	actionHandlers["configMap"] = configMapActions

	persistentVolume := resource.NewPersistentVolume(kubeClient, sendResponse)
	persistentVolumeActions := ActionHandler{
		LIST: persistentVolume.List,
	}
	actionHandlers["persistentVolume"] = persistentVolumeActions

	return &ResourceActions{
		KubeClient:            kubeClient,
		ResourceActionHandler: actionHandlers,
	}
}

func (r *ResourceActions) GetRequestHandler(resource string, action string) Handler {
	return r.ResourceActionHandler[resource][action]
}

package main

import (
	"github.com/weibeld/kaliases"
	"os"
)

/* Suites */

var suites = []kaliases.Suite{
	kaliases.Suite{},
	kaliases.Suite{actionsGet, resources, optionsGet},
	kaliases.Suite{actionsDelete, resources, optionsDelete},
	kaliases.Suite{actionsDescribe, resources, optionsDescribe},
	kaliases.Suite{actionsExec, optionsExec},
	kaliases.Suite{actionsLogs, optionsLogs},
	kaliases.Suite{actionsEdit, resources},
	kaliases.Suite{actionsOther},
}

/* Groups */

var resources = kaliases.Group{
	[]kaliases.Segment{resourcePod, resourceDeployment, resourceService, resourceNode,
		resourceIngress, resourceRole, resourceRoleBinding, resourceClusterRole,
		resourceClusterRoleBinding},
	false,
}

var actionsGet = kaliases.Group{
	[]kaliases.Segment{actionGet},
	false,
}
var optionsGet = kaliases.Group{
	[]kaliases.Segment{optionWatch, optionOutput, optionAllNamespaces},
	true,
}

var actionsDelete = kaliases.Group{
	[]kaliases.Segment{actionDelete},
	false,
}
var optionsDelete = kaliases.Group{
	[]kaliases.Segment{optionAll, optionAllNamespaces},
	true,
}

var actionsDescribe = kaliases.Group{
	[]kaliases.Segment{actionDescribe},
	false,
}
var optionsDescribe = kaliases.Group{
	[]kaliases.Segment{optionAllNamespaces},
	true,
}

var actionsExec = kaliases.Group{
	[]kaliases.Segment{actionExec},
	false,
}
var optionsExec = kaliases.Group{
	[]kaliases.Segment{optionInteractive},
	true,
}

var actionsLogs = kaliases.Group{
	[]kaliases.Segment{actionLogs},
	false,
}
var optionsLogs = kaliases.Group{
	[]kaliases.Segment{optionFollow},
	true,
}

var actionsEdit = kaliases.Group{
	[]kaliases.Segment{actionEdit},
	false,
}

var actionsOther = kaliases.Group{
	[]kaliases.Segment{actionApply, actionPortForward, actionExplain},
	false,
}

/* Action segments */

var actionGet = kaliases.Segment{
	{Short: "g", Long: "get"},
}
var actionDelete = kaliases.Segment{
	{Short: "d", Long: "delete"},
}
var actionDescribe = kaliases.Segment{
	{Short: "s", Long: "decribe"},
}
var actionEdit = kaliases.Segment{
	{Short: "e", Long: "edit"},
}
var actionExec = kaliases.Segment{
	{Short: "x", Long: "exec"},
}
var actionLogs = kaliases.Segment{
	{Short: "l", Long: "logs"},
}
var actionApply = kaliases.Segment{
	{Short: "a", Long: "apply"},
}
var actionPortForward = kaliases.Segment{
	{Short: "p", Long: "port-forward"},
}
var actionExplain = kaliases.Segment{
	{Short: "ex", Long: "explain"},
}

/* Resource segments */

var resourcePod = kaliases.Segment{
	{Short: "p", Long: "pod"},
}
var resourceDeployment = kaliases.Segment{
	{Short: "d", Long: "deployment"},
}
var resourceService = kaliases.Segment{
	{Short: "s", Long: "service"},
}
var resourceNode = kaliases.Segment{
	{Short: "n", Long: "node"},
}
var resourceIngress = kaliases.Segment{
	{Short: "i", Long: "ingress"},
}
var resourceRole = kaliases.Segment{
	{Short: "r", Long: "role"},
}
var resourceRoleBinding = kaliases.Segment{
	{Short: "rb", Long: "rolebinding"},
}
var resourceClusterRole = kaliases.Segment{
	{Short: "cr", Long: "clusterrole"},
}
var resourceClusterRoleBinding = kaliases.Segment{
	{Short: "crb", Long: "clusterrolebinding"},
}

/* Option segments */

var optionWatch = kaliases.Segment{
	{Short: "w", Long: "-w"},
}
var optionOutput = kaliases.Segment{
	{Short: "y", Long: "-o yaml"},
	{Short: "j", Long: "-o json"},
}
var optionAllNamespaces = kaliases.Segment{
	{Short: "a", Long: "--all-namespaces"},
}
var optionAll = kaliases.Segment{
	{Short: "A", Long: "--all"},
}
var optionInteractive = kaliases.Segment{
	{Short: "i", Long: "-it"},
}
var optionFollow = kaliases.Segment{
	{Short: "f", Long: "-f"},
}

func main() {
	kaliases.Generate(suites, os.Stdout)
}

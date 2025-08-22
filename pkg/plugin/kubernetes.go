package plugin

import (
	"context"
	"strings"
)

// KubernetesPlugin handles kubectl and Kubernetes errors
type KubernetesPlugin struct {
	*BasePlugin
}

// NewKubernetesPlugin creates a new Kubernetes plugin
func NewKubernetesPlugin() *KubernetesPlugin {
	return &KubernetesPlugin{
		BasePlugin: NewBasePlugin(
			"kubernetes",
			"Handles kubectl and Kubernetes cluster errors",
			85,
		),
	}
}

// Match determines if this plugin should handle the command
func (p *KubernetesPlugin) Match(cmd string, output string, exitCode int) bool {
	if exitCode == 0 {
		return false
	}

	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	k8sCommands := []string{"kubectl ", "k ", "helm ", "minikube "}
	
	for _, k8sCmd := range k8sCommands {
		if strings.HasPrefix(cmdLower, k8sCmd) {
			return true
		}
	}

	return false
}

// Suggest generates a suggestion for fixing the Kubernetes error
func (p *KubernetesPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	outputLower := strings.ToLower(output)
	
	// Check for common Kubernetes errors
	
	// Command typos
	if correctedCmd, explanation := p.checkForTypos(cmd); correctedCmd != "" {
		return correctedCmd, explanation, nil
	}
	
	// Connection refused / cluster not accessible
	if strings.Contains(outputLower, "connection refused") ||
	   strings.Contains(outputLower, "server could not be reached") {
		return "kubectl cluster-info", "Check cluster connectivity", nil
	}
	
	// No kubeconfig
	if strings.Contains(outputLower, "kubeconfig") && 
	   strings.Contains(outputLower, "not found") {
		return "kubectl config view", "Check kubeconfig configuration", nil
	}
	
	// CrashLoopBackOff
	if strings.Contains(outputLower, "crashloopbackoff") {
		podName := extractResourceName(cmd, "pod")
		if podName != "" {
			return "kubectl logs " + podName + " --previous", 
				   "Check logs from previous container instance", nil
		}
		return "kubectl get pods", "Check pod status and identify failing pods", nil
	}
	
	// ImagePullBackOff
	if strings.Contains(outputLower, "imagepullbackoff") ||
	   strings.Contains(outputLower, "errimagepull") {
		podName := extractResourceName(cmd, "pod")
		if podName != "" {
			return "kubectl describe pod " + podName, 
				   "Check pod events for image pull errors", nil
		}
		return "kubectl get pods", "Check pods with image pull issues", nil
	}
	
	// Resource not found
	if strings.Contains(outputLower, "not found") {
		resourceType := extractResourceType(cmd)
		if resourceType != "" {
			return "kubectl get " + resourceType, "List all " + resourceType + " resources", nil
		}
		return "kubectl get all", "List all resources in current namespace", nil
	}
	
	// Permission denied / RBAC
	if strings.Contains(outputLower, "forbidden") ||
	   strings.Contains(outputLower, "permission denied") {
		return "kubectl auth can-i '*' '*'", "Check current user permissions", nil
	}
	
	// Invalid YAML
	if strings.Contains(outputLower, "yaml") && 
	   (strings.Contains(outputLower, "parse") || strings.Contains(outputLower, "unmarshal")) {
		return "kubectl apply --dry-run=client -f <filename>", 
			   "Validate YAML syntax before applying", nil
	}
	
	// Context not set
	if strings.Contains(outputLower, "context") && 
	   strings.Contains(outputLower, "not found") {
		return "kubectl config get-contexts", "List available contexts", nil
	}
	
	// Namespace not found
	if strings.Contains(outputLower, "namespace") && 
	   strings.Contains(outputLower, "not found") {
		return "kubectl get namespaces", "List available namespaces", nil
	}
	
	// Pod pending
	if strings.Contains(outputLower, "pending") {
		return "kubectl describe nodes", "Check node resources and constraints", nil
	}
	
	// Service account issues
	if strings.Contains(outputLower, "serviceaccount") && 
	   strings.Contains(outputLower, "not found") {
		return "kubectl get serviceaccounts", "List available service accounts", nil
	}
	
	return "", "No specific Kubernetes suggestion available", nil
}

// checkForTypos checks for common typos in kubectl commands
func (p *KubernetesPlugin) checkForTypos(cmd string) (string, string) {
	corrections := map[string]string{
		"get":         "get",
		"gt":          "get",
		"ge":          "get",
		"gett":        "get",
		"apply":       "apply",
		"aply":        "apply",
		"appl":        "apply",
		"aplly":       "apply",
		"delete":      "delete",
		"delet":       "delete",
		"delte":       "delete",
		"deleet":      "delete",
		"describe":    "describe",
		"describ":     "describe",
		"desribe":     "describe",
		"descrbe":     "describe",
		"logs":        "logs",
		"log":         "logs",
		"lgs":         "logs",
		"exec":        "exec",
		"exe":         "exec",
		"exce":        "exec",
		"create":      "create",
		"creat":       "create",
		"crete":       "create",
		"craete":      "create",
		"edit":        "edit",
		"edi":         "edit",
		"eidt":        "edit",
		"config":      "config",
		"confi":       "config",
		"confg":       "config",
		"port-forward": "port-forward",
		"portforward": "port-forward",
		"port-foward": "port-forward",
		"rollout":     "rollout",
		"rolout":      "rollout",
		"roolout":     "rollout",
		"scale":       "scale",
		"scal":        "scale",
		"scle":        "scale",
		"patch":       "patch",
		"patc":        "patch",
		"ptach":       "patch",
		"pods":        "pods",
		"pod":         "pods",
		"po":          "pods",
		"podz":        "pods",
		"services":    "services",
		"service":     "services",
		"svc":         "services",
		"deployments": "deployments",
		"deployment":  "deployments",
		"deploy":      "deployments",
		"nodes":       "nodes",
		"node":        "nodes",
		"namespaces":  "namespaces",
		"namespace":   "namespaces",
		"ns":          "namespaces",
	}

	words := strings.Fields(cmd)
	changed := false
	
	for i, word := range words {
		cleanWord := strings.ToLower(word)
		if correction, exists := corrections[cleanWord]; exists {
			words[i] = correction
			changed = true
		}
	}

	if changed {
		newCmd := strings.Join(words, " ")
		return newCmd, "Fixed typos in kubectl command"
	}

	return "", ""
}

// extractResourceType extracts the Kubernetes resource type from command
func extractResourceType(cmd string) string {
	words := strings.Fields(cmd)
	
	for i, word := range words {
		if (word == "get" || word == "describe" || word == "delete") && i+1 < len(words) {
			return words[i+1]
		}
	}
	
	return ""
}

// extractResourceName extracts the resource name from command
func extractResourceName(cmd, resourceType string) string {
	words := strings.Fields(cmd)
	
	for i, word := range words {
		if word == resourceType && i+1 < len(words) {
			return words[i+1]
		}
	}
	
	return ""
}

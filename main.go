package main

import (
    "context"
    "fmt"
    "os"

    configv1 "github.com/openshift/api/config/v1" // Import the config/v1 API
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/controller-runtime/pkg/client" // Use controller-runtime client
)

// getOpenshiftPlatformType uses the OpenShift API to get the infrastructure platform type.
func getOpenshiftPlatformType(kubeconfigPath string) (string, error) {
    // 1. Load the Kubernetes configuration.
    config, err := loadConfig(kubeconfigPath)
    if err != nil {
        return "", fmt.Errorf("failed to load Kubernetes config: %v", err)
    }

    // 2. Create a Kubernetes client.
    k8sClient, err := createClient(config)
    if err != nil {
        return "", fmt.Errorf("failed to create Kubernetes client: %v", err)
    }

    // 3. Get the Infrastructure resource.
    infrastructure := &configv1.Infrastructure{}
    if err := k8sClient.Get(context.TODO(), client.ObjectKey{Name: "cluster"}, infrastructure); err != nil {
        return "", fmt.Errorf("failed to get Infrastructure 'cluster' resource: %v", err)
    }

    // 4. Extract and return the platform type.
    platformType := string(infrastructure.Status.PlatformStatus.Type) // Directly access the type
    if platformType == "" {
        return "", fmt.Errorf("platformType is empty in the Infrastructure status")
    }
    return platformType, nil
}

// loadConfig loads the Kubernetes configuration from a kubeconfig file or the default location.
func loadConfig(kubeconfigPath string) (*rest.Config, error) {
    if kubeconfigPath != "" {
        // Use the specified kubeconfig path.
        return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    }
    // Use the default kubeconfig path or in-cluster config.
    return rest.InClusterConfig() //check for in-cluster config first
}

// createClient creates a Kubernetes client from the given configuration.
func createClient(config *rest.Config) (client.Client, error) {
    scheme := runtime.NewScheme()
    if err := configv1.AddToScheme(scheme); err != nil { //register openshift scheme
        return nil, err
    }

    k8sClient, err := client.New(config, client.Options{Scheme: scheme})
    return k8sClient, err
}

func main() {
    // Allow the user to specify the kubeconfig path as a command-line argument.
    var kubeconfigPath string
    if len(os.Args) > 1 {
        kubeconfigPath = os.Args[1]
    }

    // Call the function to get the platform type.
    platformType, err := getOpenshiftPlatformType(kubeconfigPath)
    if err != nil {
        // Print the error to standard error and exit with a non-zero exit code.
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Print the platform type to standard output.
    fmt.Println(platformType)
}




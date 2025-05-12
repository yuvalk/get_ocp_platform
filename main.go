package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// Import the OpenShift config API types
	configv1 "github.com/openshift/api/config/v1"

	// Import Kubernetes client-go libraries
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Determine Kubernetes configuration path or use in-cluster config
	// This logic is standard client-go practice, not directly from sources,
	// but necessary for a runnable program.
	kubeconfigPath := os.Getenv("KUBECONFIG")
	var config *rest.Config
	var err error

	if kubeconfigPath == "" {
		// Try in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			// Fallback to default config path if in-cluster fails
			// Using clientcmd.NewDefaultClientConfigLoadingRules() is standard
			// client-go for loading ~/.kube/config
			kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().Get </dev/null>
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				kubeconfig,
				&clientcmd.ConfigOverrides{},
			).ClientConfig()
			if err != nil {
				log.Fatalf("Error building Kubernetes config: %v", err)
			}
		}
	} else {
		// Use specified kubeconfig file
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Fatalf("Error building Kubernetes config from %s: %v", kubeconfigPath, err)
		}
	}


	// Create a client for the config.openshift.io/v1 API group
	// This uses the config package imported above to add the scheme.
	// Source [1] shows getting client instances like configClient.
	client, err := rest.NewRESTClient(config)
	if err != nil {
		log.Fatalf("Error creating REST client: %v", err)
	}

	// Add the configv1 scheme to the client's serializer.
	// This allows the client to understand Infrastructure objects.
	scheme := configv1.AddToScheme
	if err := scheme(client.NewRequest("").SetHeader("Accept", "application/json").Verb("").Name("").Namespace("").Resource("").SubResource("").Param("", "").Timeout(0).Do(context.TODO()).Into(nil)); err != nil {
       log.Fatalf("Error adding configv1 scheme: %v", err)
    }
    // NOTE: The scheme registration is typically done via a *Scheme struct
    // and added to a dynamic client or a scheme-aware client-go client.
    // A more standard way (less error prone with RESTClient) is below using a dynamic client or typed client.

	// --- Alternative and more standard way using a typed client ---
	// This is a more common pattern shown in client-go examples for custom resources.
	// Source [1] implies use of typed clients like configClient.
	configClient, err := configv1.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating configv1 client: %v", err)
	}

	// Get the "cluster" Infrastructure resource
	// Source [1] explicitly shows getting infraConfig using client.Get or configClient...Get
	infraConfig, err := configClient.Infrastructures().Get(context.TODO(), "cluster", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error getting Infrastructure resource: %v", err)
	}

	// Access the platform type from the status
	// The `.status.platform` field is deprecated [Source from prior turn]
	// The preferred field is `.status.platformStatus.type` [Source from prior turn, implied by 57]
	platformType := "Unknown" // Default
	// Deprecated field access (optional, for demonstration)
	deprecatedPlatform := infraConfig.Status.Platform
	if deprecatedPlatform != "" {
		log.Printf("Deprecated .status.platform: %s", deprecatedPlatform) // Log as deprecated
	}

	// Preferred field access
	if infraConfig.Status.PlatformStatus != nil {
		platformType = string(infraConfig.Status.PlatformStatus.Type) // Convert configv1.PlatformType to string
	}


	// Print the result
	fmt.Printf("Cluster Platform Type: %s\n", platformType)
}



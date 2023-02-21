package main

import (
	"os"
	"fmt"
	"flag"
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/tools/clientcmd"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
)

func getk8sSecretClient() coreV1Types.SecretInterface {

	kc := os.Getenv("HOME") + "/.kube/config"
	c, err := clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		panic("config building failed:" + err.Error())
	}
	nc, err := kubernetes.NewForConfig(c) 
	if err != nil {
		panic("creating new config failed:" + err.Error())
	}
	sc := nc.CoreV1().Secrets("default")
	return sc
}


func main() {

	cp := flag.String("cert-path", "./build/secrets/server_crt.pem", "The path to your end-entity certificate")
	kp := flag.String("key-path", "./build/secrets/server_key.pem", "The path to your end-entity private key")
    flag.Parse()

	cert, err := os.ReadFile(*cp)
	if err != nil {
    	panic("failed to read certificate: " + err.Error())
    } 

	key, err := os.ReadFile(*kp)
	if err != nil {
    	panic("failed to parse private key: " + err.Error())
    }

	sc := getk8sSecretClient()

	// Retry updating secret until you no longer get a conflict error. 
	// This way, you can preserve changes made by other clients between.
	// Ref: https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		s, getErr := sc.Get(context.TODO(), "tls-secret", metaV1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of secret tls-secret: %v", getErr))
		}
		s.Data["tls.crt"] = []byte(cert)
		s.Data["tls.key"] = []byte(key)

		_, updateErr := sc.Update(context.TODO(), s, metaV1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Secret tls-secret is successfully updated")

}
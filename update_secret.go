package main

import (
	"os"
	"fmt"
	"flag"
	"context"
	"encoding/pem"
	"encoding/base64"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
)


func getk8sSecretClient() coreV1Types.SecretInterface {

	kc := os.Getenv("HOME") + "/.kube/config"
	c, err := clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		fmt.Printf("config building failed: %v\n", err.Error())
	}
	nc, err := kubernetes.NewForConfig(c) 
	if err != nil {
		fmt.Printf("creating new config failed: %v\n", err.Error())
	}
	sc := nc.CoreV1().Secrets("default")
	return sc
}


func main() {

	cp := flag.String("cert-path", "./build/secrets/server_crt.pem", "The path to your end-entity certificate")
	kp := flag.String("key-path", "./build/secrets/server_key.pem", "The path to your end-entity private key")
  flag.Parse()

	cpem, err := os.ReadFile(*cp)
	if err != nil {
    	fmt.Printf("failed to read certificate: %v\n", err.Error())
  } 

	kpem, err := os.ReadFile(*kp)
	if err != nil {
			fmt.Printf("failed to parse private key: %v\n", err.Error())
	}

	cblock, _ := pem.Decode(cpem)
	if cblock == nil || cblock.Type != "CERTIFICATE" {
		fmt.Println("failed to decode PEM block containing certificate")
	}

	kblock, _ := pem.Decode(kpem)
	if kblock == nil || kblock.Type != "PRIVATE KEY" {
		fmt.Println("failed to decode PEM block containing private key")
	}

	cert := base64.StdEncoding.EncodeToString(cblock.Bytes)
	key := base64.StdEncoding.EncodeToString(kblock.Bytes)

	sc := getk8sSecretClient()

	// Retry updating secret until you no longer get a conflict error. 
	// This way, you can preserve changes made by other clients between.
	// Ref: https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		s, getErr := sc.Get(context.TODO(), "tls-secret", metaV1.GetOptions{})

		// If secret "tls-secret" is not found
		if getErr != nil {
			sd := make(map[string]string)
			sd["tls.crt"] = cert
			sd["tls.key"] = key

			s_ := &v1.Secret{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "tls-secret",
					Namespace: "default",
				},
				StringData: sd,
			}

			_, createErr := sc.Create(context.TODO(), s_, metaV1.CreateOptions{})
			if createErr != nil {
				fmt.Printf("Update failed: %v\n", createErr)
			}
			return createErr
		}

		s.StringData["tls.crt"] = cert
		s.StringData["tls.key"] = key

		_, updateErr := sc.Update(context.TODO(), s, metaV1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		fmt.Printf("Update failed: %v\n", retryErr)
	}
	fmt.Println("Secret tls-secret is successfully updated")

}

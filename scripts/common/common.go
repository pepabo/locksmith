package common

import (
	"os"
	"log"
	"fmt"
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
)


func getk8sSecretClient(namespace string) coreV1Types.SecretInterface {

	kc := os.Getenv("HOME") + "/.kube/config"
	c, err := clientcmd.BuildConfigFromFlags("", kc)

	if err != nil {
		log.Fatalf("config building failed: %v\n", err.Error())
	}

	nc, err := kubernetes.NewForConfig(c) 
	if err != nil {
		log.Fatalf("creating new config failed: %v\n", err.Error())
	}

	sc := nc.CoreV1().Secrets(namespace)
	return sc
}


func Updatek8sSecret(namespace string, secretname string, cpem []byte, kpem []byte) {

	sc := getk8sSecretClient(namespace)
	cert := string(cpem)
	key := string(kpem)

	// Retry updating secret until you no longer get a conflict error. 
	// This way, you can preserve changes made by other clients between.
	// Ref: https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		s, getErr := sc.Get(context.TODO(), secretname, metaV1.GetOptions{})

		// If secret "tls-secret" is not found
		if errors.IsNotFound(getErr) {

			sd := make(map[string]string)
			sd["tls.crt"] = cert
			sd["tls.key"] = key

			s_ := &v1.Secret{
				Type: v1.SecretTypeTLS,
				ObjectMeta: metaV1.ObjectMeta{
					Name: secretname,
					Namespace: namespace,
				},
				StringData: sd,
			}

			_, createErr := sc.Create(context.TODO(), s_, metaV1.CreateOptions{})
			if createErr != nil {
				log.Fatalf("Update failed: %v\n", createErr)
			}
			return createErr
		}

		// If you forget to add this statement, you will get an error (panic: assignment to entry in nil map)
		if s.StringData == nil {
			s.StringData = map[string]string{}
		}

		s.StringData["tls.crt"] = cert
		s.StringData["tls.key"] = key

		_, updateErr := sc.Update(context.TODO(), s, metaV1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		log.Fatalf("Update failed: %v\n", retryErr)
	}
	
	fmt.Println("Secret tls-secret is successfully updated")
	os.Exit(0)

}

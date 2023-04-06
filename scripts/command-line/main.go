package main

import (
	"os"
	"log"
	"flag"
	"encoding/pem"
	"github.com/pepabo/locksmith/scripts/common"
)

func main() {

	cp := flag.String("cert-path", "../build/secrets/server_crt.pem", "The path to your end-entity certificate")
	kp := flag.String("key-path", "../build/secrets/server_key.pem", "The path to your end-entity private key")
	ns := flag.String("namespace", "default", "Your namespace")
	sn := flag.String("secret-name", "tls-secret", "Your secret name")
  flag.Parse()

	cpem, err := os.ReadFile(*cp)
	if err != nil {
    	log.Fatalf("failed to read certificate: %v\n", err.Error())
  } 

	kpem, err := os.ReadFile(*kp)
	if err != nil {
			log.Fatalf("failed to parse private key: %v\n", err.Error())
	}

	cblock, _ := pem.Decode(cpem)
	if cblock == nil || cblock.Type != "CERTIFICATE" {
		log.Fatal("failed to decode PEM block containing certificate")
	}

	kblock, _ := pem.Decode(kpem)
	if kblock == nil || kblock.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	common.Updatek8sSecret(*ns, *sn, cpem, kpem)
	
}

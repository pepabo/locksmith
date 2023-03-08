# locksmith
Open source software that provides AWS temporary credentials with k8s clusters running outside AWS

## Usage

### Create your own root CA certificate and private key
- Your new root CA certificate is stored in locksmith/build/secrets/ca_crt.pem
- Your new root CA private key is stored in locksmith/build/secrets/ca_key.pem

```
cd locksmith/build
chmod +x ./create_secret_ca.sh
./create_secret_ca.sh
```
### Create your own server certificate and private key
- Your new server certificate is stored in locksmith/build/secrets/server_crt.pem
- Your new server private key is stored in locksmith/build/secrets/server_key.pem

```
cd locksmith/build
chmod +x ./create_secret_server.sh
./create_secret_server.sh
```

### Create a trustanchor on AWS

### Create a special IAM role on AWS

### Create an AWS profile on AWS

### Create a docker image for locksmith

### Add locksmith to your manifest file

### Run your deployment on your k8s cluster

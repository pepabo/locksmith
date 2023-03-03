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

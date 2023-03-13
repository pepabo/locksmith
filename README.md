# locksmith
Open source software that provides AWS temporary credentials with k8s clusters running outside AWS

### Setup

#### (Optional) Create your own root CA certificate and private key
- Your new root CA certificate is stored in locksmith/build/secrets/ca_crt.pem
- Your new root CA private key is stored in locksmith/build/secrets/ca_key.pem

```
cd locksmith/build
chmod +x ./create_secret_ca.sh
./create_secret_ca.sh
```
#### (Optional) Create your own server certificate and private key
- Your new server certificate is stored in locksmith/build/secrets/server_crt.pem
- Your new server private key is stored in locksmith/build/secrets/server_key.pem

```
cd locksmith/build
chmod +x ./create_secret_server.sh
./create_secret_server.sh
```

### Create a trustanchor on AWS

#### 1. Click "Create a trust anchor"
![trust-anchor](/images/trust-anchor.png)

#### 2. Paste your Root CA key to External certificate bundle
![create-trust-anchor](/images/create-trust-anchor.png)


#### 3. Create a special IAM role for IAM Roles Anywhere
`rolesanywhere-trust-policy.json`
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "rolesanywhere.amazonaws.com"
            },
            "Action": [
                "sts:AssumeRole",
                "sts:SetSourceIdentity",
                "sts:TagSession"
            ]
        }
    ]
}
```

#### 4. Create an AWS role that you are planning to assume

#### 5. Create an AWS profile on AWS
![profile](/images/profile.png)

#### 6. Create a docker image for locksmith

##### 6.1 Execute commands reuired for building the docker image 
```
cd locksmith
export AWS_TRUST_ANCHOR_ARN=(ARN for your trust anchor)
export AWS_PROFILE_ARN=(ARN for your AWS profile)
export AWS_ROLE_ARN=(ARN for the AWS role that you are going to assume)
export AWS_REGION=(your AWS region)
```

##### 6.2 Build your docker image
```
docker compose up -d
```

#### 7. Create k8s secret for server certificate and private key
```
kubectl create secret tls tls-secret \
  --cert=(path to your server certificate) \
  --key=(path to your private key)
```

#### 6. Create k8s secret for AWS ARNS
```
kubectl create secret generic aws-config \
        --from-literal="aws-trust-anchor-arn=$AWS_TRUST_ANCHOR_ARN" \
        --from-literal="aws-profile-arn=$AWS_PROFILE_ARN" \
        --from-literal="aws-role-arn=$AWS_ROLE_ARN" \
        --from-literal="aws-region=$AWS_REGION"
```

#### 8. Add locksmith to your manifest file
See an [example](k8s/deployment.yaml) 

#### 9. Run your deployment on your k8s cluster

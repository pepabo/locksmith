#!/bin/bash

cp /usr/local/bin/aws_signing_helper /opt/locksmith/bin
cat <<EOF > /root/.aws/config
[default]
    credential_process = /opt/locksmith/bin/aws_signing_helper credential-process --certificate /opt/pki/server/crt/server_crt.pem --private-key /opt/pki/server/key/server_key.pem --region ${AWS_REGION} --trust-anchor-arn ${AWS_TRUST_ANCHOR_ARN} --profile-arn ${AWS_PROFILE_ARN} --role-arn ${AWS_ROLE_ARN}
EOF
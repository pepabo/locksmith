docker build ./ -t ca_image -f dockerfile_ca
docker run --name ca -it ca_image:latest
export CA_CONTAINER_ID=$(docker ps -a | grep ca_image:latest | awk '{print $1}')

docker cp $CA_CONTAINER_ID:/opt/pki/RootCA/ca_crt.pem secrets/ca_crt.pem
docker cp $CA_CONTAINER_ID:/opt/pki/RootCA/ca_key.pem secrets/ca_key.pem

echo "Certificate for your root CA is now stored in locksmith/build/secrets/ca_crt.pem"
echo "Private key for your root CA is now stored in locksmith/build/secrets/ca_key.pem"

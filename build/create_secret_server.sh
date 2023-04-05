docker build ./ -t server_image -f dockerfile_server
docker run --name server -it server_image:latest
export SERVER_CONTAINER_ID=$(docker ps -a | grep server_image:latest | awk '{print $1}')

docker cp $SERVER_CONTAINER_ID:/opt/pki/RootCA/server_crt.pem secrets/server_crt.pem
docker cp $SERVER_CONTAINER_ID:/opt/pki/RootCA/server_key.pem secrets/server_key.pem

echo "Certificate for your server is now stored in locksmith/build/secrets/server_crt.pem"
echo "Private key for your server is now stored in locksmith/build/secrets/server_key.pem"

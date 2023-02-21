cd /opt/pki/RootCA

mkdir newcerts
echo "01" > serial
echo "00" > crlnumber
touch index.txt

openssl genrsa -out server_key.pem 2048

openssl req -new \
 -subj "/C=JP/ST=Tokyo/O=EXAMPLE/CN=EXAMPLE Client" \
 -out server_csr.pem  \
 -key server_key.pem

 openssl ca -config openssl_client.cnf  \
 -batch -extensions v3_client \
 -out server_crt.pem \
 -in server_csr.pem \
 -cert ca_crt.pem \
 -keyfile ca_key.pem \
 -passin pass:rootcaprivkeypass

 openssl x509 -in server_crt.pem -out server_crt.pem
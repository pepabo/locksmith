
cd /opt/pki/RootCA

mkdir newcerts
echo "01" > serial
echo "00" > crlnumber
touch index.txt

openssl genrsa -out ca_key.pem -aes256 -passout pass:rootcaprivkeypass 2048

openssl req -new \
 -subj "/C=JP/ST=Tokyo/O=EXAMPLE/CN=EXAMPLE Root CA" \
 -out ca_csr.pem \
 -key ca_key.pem \
 -passin pass:rootcaprivkeypass

openssl ca -config openssl_ca.cnf -batch -extensions v3_ca \
 -out ca_crt.pem \
 -in ca_csr.pem \
 -selfsign \
 -keyfile ca_key.pem \
 -passin pass:rootcaprivkeypass

openssl x509 -in ca_crt.pem -out ca_crt.pem
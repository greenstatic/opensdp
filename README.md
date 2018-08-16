# OpenSDP
A proof of concept Software Defined Perimeter (SDP) implementation using OpenSPA for service hiding


## How To

### Keys

```bash
# Create CA
openssl genrsa -out ignore/ca.key 2048
openssl req -new -x509 -key ignore/ca.key -out ignore/ca.crt

# Create server key
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -out server.crt -days 365 -CAcreateserial

# Create client key
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -out client.crt -days 365 -CAcreateserial
```


## License
This software is licensed under: [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

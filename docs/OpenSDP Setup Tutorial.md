# OpenSDP Setup Tutorial
Being a PoC, we currently do not provide binary releases.
This means you will need to compile the software yourself.
But don't worry it's easy, it's [Go](https://golang.org/)!

## Requirements
* OpenSPA server setup (with a test client)
* Debian/Ubuntu based system

## RSA Keys
Clients authenticate with the server via mutual TLS.
This means that the servers and clients will need a signed certificate by a commonly trusted CA.

### Creating Our Own (Self Signed) CA
```bash
mkdir -p ~/opensdp/keys
cd ~/opensdp/keys

openssl genrsa -out ./ca.key 2048

# Fill out the certificate info as you like
openssl req -new -x509 -key ./ca.key -out ./ca.crt
```

### Creating the Server and Client(s) Keys
Let's create our server keypair
```bash
openssl genrsa -out server.key 2048

# Fill out the certificate info as you like EXCEPT the common name (CN)!
# The CN should be "OpenSDP-server"
openssl req -new -key server.key -out server.csr

# Sign the CSR with our CA to create a 365 day valid cert
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -out server.crt -days 365 -CAcreateserial

# We don't need the CSR anymore
rm server.csr
```

Next, let's create a client's keypair (this step is identical for all clients).
For this step you need to have generated an OpenSPA client OSPA file for EACH client you wish to create a keypair here.
```bash
openssl genrsa -out client.key 2048

# Fill out the certificate info as you like EXCEPT the common name (CN)!
# The CN should be the client's device ID (UUID), eg; 9f84fbb8-10e8-4b8a-abd2-bb91cbf484df.
openssl req -new -key client.key -out client.csr

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -out client.crt -days 365 -CAcreateserial
rm client.csr
```

### Build the Software
Let's build the server.
First get the source code.
````bash
go get github.com/greenstatic/opensdp
````
Don't worry if you get a *package github.com/greenstatic/opensdp: no Go files in /home/ubuntu/go/src/github.com/greenstatic/opensdp* error message.

Now let's build it.
```bash
cd ~/go/src/github.com/greenstatic/opensdp/cmd/opensdp-server
go get -u ./... # this may take some time

go build -o ~/opensdp/opensdp-server
```

Great, we have our executable now.
Let's build the client now.

```bash
cd ~/go/src/github.com/greenstatic/opensdp/cmd/opensdp-client
go get -u ./... # this may take some time

go build -o ~/opensdp/opensdp-client
```

Next copy the example config files.
```bash
cd ~/go/src/github.com/greenstatic/opensdp/config
mkdir ~/opensdp/clients

cp client/config.yaml ~/opensdp/clients/config.yaml
cp server/*.yaml ~/opensdp/
```

### Configure the Server & Client
First the server.
```bash
cd ~/opensdp
# Update the ca-cert, certificate and key fields to reflect the
# ~/openspa/keys dir (we recommend using absolute paths to avoid errors)
nano config.yaml

# Configure this to your liking
nano services.yaml 

# Configure this to your liking,
# Please only use services that you defined earlier otherwise
# the server will reject the file.
nano clients.yaml
```

Next the client.
```bash
# Move the previously generated keys to the clients dir
mv keys/client.* clients/

# For now just update the server field (the default port is 33311)
nano clients/config.yaml
```

You have successfully created the necessary client files.
Copy the following to your client:
* `~/opensdp/clients/client.crt`
* `~/opensdp/clients/client.key`
* `~/opensdp/clients/config.yaml`
* `~/opensdp/keys/ca.crt`

And of course the `opensdp-client`.

### Starting the Server
Since it's only a PoC we'll run the OpenSDP-server using screen.
```bash
cd ~/opensdp
screen -d -m ~/opensdp/opensdp-server
```

### Starting the Clients
If you had a working OpenSPA installation before and the required client OSPA file, edit your client's `config.yaml` file.
The path of the OpenSPA-client should be specified (absolute path) and the path of the OSPA file.
Make sure your keypairs and ca cert are specified as well.
Once you have everything setup, run:
````bash
opensdp services --config ./config.yaml --verbose
````
This should print out a list of available services.
If there is absolutely no response from the OpenSDP server, then you probably have an OpenSPA related configuration issue.
In case OpenSPA is not the issue, check all the OpenSDP yaml files.
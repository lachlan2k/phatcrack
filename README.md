# Phatcrack

Distributed hash cracking for chonkers. Phatcrack is a frontend for Hashcat that provides:
* Modern web-interface, supporting most common Hashcat attacks
* Multi-user and project access control
* Work distribution across multiple workers

### Deployment

#### Server
Docker is the only supported deployment method for the server. 

```sh
# Ideally the container processes should be run rootless, so we'll create an unprivileged user.
adduser --system --no-create-home phatcrack-server

mkdir -p /srv/containers/phatcrack
cd /srv/containers/phatcrack

wget https://raw.githubusercontent.com/lachlan2k/phatcrack/main/docker-compose.prod.yml

echo "HOST_NAME=phatcrack.lan" >> .env
echo "DB_PASS=$(openssl rand -hex 16)" >> .env
echo "PHATCRACK_USER=phatcrack-server" >> .env

# Make a directory to persist files in
mkdir -p /srv/containers/phatcrack/filerepo
chown phatcrack-server:phatcrack-server /srv/containers/phatcrack/filerepo

# If you chose a hostname that is publicly accessible and expose this to the world (not recommended), Caddy will automatically deploy TLS.

## Otherwise, use the following for self-signed TLS
# echo "TLS_OPTS=tls internal" >> .env

## If you want to supply custom certificates, place them in a directory called `certs`
## And add ./certs:/etc/caddy/certs:ro as a mount in docker-compose.prod.yml for 
# echo "TLS_OPTS=tls /etc/caddy/certs/cert.pem /etc/caddy/certs/key.pem" >> .env

docker compose up -d
```

The default credentials are `admin:changeme`.

#### Agents

For each agent you want to enroll, visit the admin GUI, and create a new agent. Note down the generated API key, as this won't be shown again.

On each agent, you can manually set up the agent as follows:

```sh
# Create a user for the phatcrack agent
adduser --system --no-create-home phatcrack-agent

# Depending on your distro, you may need to the phatcrack-agent to a group
# usermod -aG video phatcrack-agent

# Place the compiled agent binary in /opt/phatcrack/
mkdir -p /opt/phatcrack/hashcat
mkdir -p /opt/phatcrack/listfiles

cd /opt/phatcrack/hashcat/
wget https://github.com/hashcat/hashcat/releases/download/v6.2.6/hashcat-6.2.6.7z -q -O hashcat.7z
7z x hashcat.7z
rm hashcat.7z
mv hashcat-6.2.6 hashcat

mkdir -p /etc/phatcrack-agent/
wget https://raw.githubusercontent.com/lachlan2k/phatcrack/main/agent/example_config.json -O /etc/phatcrack-agent/config.json
echo -n "API KEY HERE" > /etc/phatcrack-agent/auth.key

# Update api_endpoint to point to your phatcrack installation. Ensure to use HTTPS if you have deployed it.
vi /etc/phatcrack-agent/config.json

# Set secure permissions
chown -R phatcrack-agent:phatcrack-agent /etc/phatcrack-agent
chmod 700 /etc/phatcrack-agent
chmod 600 /etc/phatcrack-agent/auth.key
chmod 600 /etc/phatcrack-agent/config.json

# Now, for testing, you can just manually invoke it. TODO: Systemd
/opt/phatcrack/agent
```

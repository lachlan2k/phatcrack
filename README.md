# Phatcrack

Phatcrack is a modern solution for distributed hash cracking, designed for hackers and other information security professionals.

Key features include:
* Built on [Hashcat](https://hashcat.net), supporting most common Hashcat attacks and almost all hash types.
* Excellent UX for manging projects, configuring attack settings, and viewing results.
* Distributes attacks, allowing both dictionary-based and mask-based attacks to be split across multiple workers.
* Automatically synchronises wordlists & rulefiles to all workers. Low-privileged users can be granted permission to upload wordlists & rulefiles.
* Modern web-interface, with multi-user support and project-based access control.

![image](https://github.com/lachlan2k/phatcrack/assets/4683714/b10df9ec-ed5a-4678-9442-89003636bbce)

### Deployment

#### Server
Docker is the only supported deployment method for the server. The following instructions assume you already have [Docker installed on your server](https://docs.docker.com/engine/install/), and are logged in as root (`sudo su`).

```sh
# Ideally the container processes should be run rootless, so we'll create an unprivileged user.
adduser --system --home /opt/phatcrack-server phatcrack-server

cd /opt/phatcrack-server

wget https://github.com/lachlan2k/phatcrack/releases/download/v0.2.0/docker-compose.yml

# Update your hostname here:
echo "HOST_NAME=phatcrack.lan" >> .env
echo "DB_PASS=$(openssl rand -hex 16)" >> .env
echo "PHATCRACK_USER=$(id -u phatcrack-server):$(id -g phatcrack-server)" >> .env
chmod 600 .env

# If you chose a hostname that is publicly accessible and expose this to the world (not recommended), Caddy will automatically deploy TLS.

## Otherwise, use the following for self-signed TLS
# echo "TLS_OPTS=tls internal" >> .env

## If you want to supply custom certificates, place them in a directory called `certs`
## And add ./certs:/etc/caddy/certs:ro as a mount in docker-compose.prod.yml for 
# echo "TLS_OPTS=tls /etc/caddy/certs/cert.pem /etc/caddy/certs/key.pem" >> .env

# Make a directory to persist files in
mkdir filerepo
chown phatcrack-server:phatcrack-server filerepo

docker compose up -d
```

You can then visit your local installation. The default credentials are `admin:changeme`.

#### Agents

For each agent you want to enroll, visit the admin GUI, and create a new agent. Note down the generated API key, as this won't be shown again.

On each agent, you can manually set up the agent as follows:

```sh
# Create a user for the phatcrack agent
adduser --system --home /opt/phatcrack-agent phatcrack-agent

# Depending on your distro, you may need to the phatcrack-agent to a group
usermod -aG video phatcrack-agent

cd /opt/phatcrack-agent

# Download hashcat
wget https://github.com/hashcat/hashcat/releases/download/v6.2.6/hashcat-6.2.6.7z -q -O hashcat.7z
7z x hashcat.7z
rm hashcat.7z
mv hashcat-6.2.6 hashcat
chown -R phatcrack-agent:phatcrack-agent ./hashcat

# Download the phatcrack-agent program from the local server
wget https://phatcrack.lan/phatcrack-agent
# Or, you can download from https://github.com/lachlan2k/phatcrack/releases/download/v0.2.0/phatcrack-agent

chmod +x ./phatcrack-agent
./phatcrack-agent install -defaults -api-endpoint https://phatcrack.lan/api/v1 -auth-key API_KEY_FROM_SERVER_HERE 

systemctl enable phatcrack-agent
systemctl start phatcrack-agent
```
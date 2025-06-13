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
useradd --system --create-home --home-dir /opt/phatcrack-server phatcrack-server

cd /opt/phatcrack-server

wget https://github.com/lachlan2k/phatcrack/releases/download/v0.6.8/docker-compose.yml

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

To enroll an agent, visit the admin GUI, and click "Register Agent". The web interface will provide a script that you can run on the agent to enroll it.

However, if you want to set up an agent manually, you can do so with the following commands. You will need to replace `REGISTRATION_KEY_FROM_SERVER_HERE` with the key from the server.

```sh
# Create a user for the phatcrack agent
useradd --system --create-home --home-dir /opt/phatcrack-agent phatcrack-agent

# Depending on your distro, you may need to the phatcrack-agent to a group
usermod -aG video phatcrack-agent

cd /opt/phatcrack-agent

# Download the phatcrack-agent program from the local server
wget https://phatcrack.lan/agent-assets/phatcrack-agent
# Or, you can download from https://github.com/lachlan2k/phatcrack/releases/download/v0.6.8/phatcrack-agent

chmod +x ./phatcrack-agent
# Optionally add -disable-tls-verification if you are using self-signed certs
./phatcrack-agent install -defaults -api-endpoint https://phatcrack.lan/api/v1 -registration-key REGISTRATION_KEY_FROM_SERVER_HERE 

systemctl enable phatcrack-agent
systemctl start phatcrack-agent
```
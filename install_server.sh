#!/bin/bash

export PHATCRACK_VERSION_TAG=v0.5.2

set -e

is_yes() {
    [[ "$1" =~ ^([yY][eE][sS]|[yY])$ ]]
}

if [ "$EUID" -ne 0 ]; then
    echo "This script must be run as root. Exiting..."
    exit 1
fi

if id "phatcrack-server" &>/dev/null || [ -d "/opt/phatcrack-server" ]; then
    echo "Warning: It appears that there is an existing installation of Phatcrack."
    echo "Please clean up by ensuring the phatcrack-server user and /opt/phatcrack-server directory do not exist."
    echo "(userdel --remove phatcrack-server)"
    exit 1
fi

if ! command -v docker &>/dev/null; then
    echo "Docker is not installed on this system."
    read -p "Do you want to install Docker? (yes/no): " install_docker

    if is_yes "$install_docker"; then
        echo "Installing Docker..."

        case "$(. /etc/os-release && echo "$ID")" in

            # EL-derivatvie distros not supported by get.docker.com
            rocky|almalinux)
                sudo yum install -y yum-utils
                sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
                sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
            ;;
            
            *)
                curl -fsSL https://get.docker.com | bash
            ;;

        esac

        if ! command -v docker &>/dev/null; then
            echo "Docker installation failed. Exiting..."
            exit 1
        else
            echo "Docker installed successfully."
            systemctl enable --now docker
        fi
    else
        echo "Docker is required to run this script. Exiting..."
        exit 1
    fi
else
    echo "Docker is already installed."
fi

echo "Creating phatcrack-server user..."
useradd --system --create-home --home-dir /opt/phatcrack-server phatcrack-server

cd /opt/phatcrack-server

echo "Downloading docker-compose.yml..."
curl -qLO https://github.com/lachlan2k/phatcrack/releases/download/$PHATCRACK_VERSION_TAG/docker-compose.yml

echo "PHATCRACK_VERSION_TAG=${PHATCRACK_VERSION_TAG}"
chmod 600 .env

read -p "What DNS hostname will resolve to your Phatcrack instance (leave blank for anything)?: " server_hostname
if [ "$server_hostname" == "" ]; then

    echo "HOST_NAME=:443" >> .env
    echo "TLS_OPTS=\"tls internal {\\non_demand\\n}\"" >> .env

else
    echo "HOST_NAME=$server_hostname" >> .env
    read -p "Would you like to use self-signed certificates? (yes/no): " use_self_signed
    if is_yes "$use_self_signed"; then
        echo "TLS_OPTS=tls internal" >> .env
    else
        read -p "Would you like to provide your own certificates? (yes/no): " provide_certs

        if is_yes "$provide_certs"; then
            mkdir ./certs
            sed -i '/^\s*# - \.\/certs:\/etc\/caddy\/Certs:ro/s/^# //' docker-compose.yml
            echo "TLS_OPTS=tls /etc/caddy/certs/cert.pem /etc/caddy/certs/key.pem" >> .env

            echo "Please provide your certificates files cert.pem and key.pem in /opt/phatcrack-server/certs/"
            echo "You may need to restart the server (docker compose restart)"
        else

            read -p "Would you like to use Let's Encrypt to provision certificates ($server_hostname must be publicly accessible) ?" use_letsencrypt

            if is_yes "$use_letsencrypt"; then
                # Default, doesnt need anything
                :
            else
                echo "No supported TLS configuration was accepted"
                exit 1
            fi

        fi
    fi

fi


echo "DB_PASS=$(openssl rand -hex 16)" >> .env
echo "PHATCRACK_USER=$(id -u phatcrack-server):$(id -g phatcrack-server)" >> .env

chmod 600 .env

mkdir filerepo
chown phatcrack-server:phatcrack-server filerepo

echo "Starting Phatcrack"
docker compose up -d
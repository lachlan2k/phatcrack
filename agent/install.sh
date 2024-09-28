#!/bin/sh

set -xe

if [ -z "$PHATCRACK_HOST" ]; then
  echo "PHATCRACK_HOST is not set. Exiting."
  exit 1
fi

if [ -z "$PHATCRACK_API_KEY" ]; then
  echo "PHATCRACK_API_KEY is not set. Exiting."
  exit 1
fi

if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Exiting."
  exit 1
fi

if ! command -v 7z &> /dev/null; then
    echo "7z is not installed. Please install 7-Zip to continue."
    exit 1
fi

download_file() {
    local url="$1"
    local output="$2"
    
    if command -v curl &> /dev/null; then
        if [[ -n "$DISABLE_TLS_VERIFICATION" ]]; then
            curl -k "$url" -o "$output"
        else
            curl "$url" -o "$output"
        fi
    elif command -v wget &> /dev/null; then
        if [[ -n "$DISABLE_TLS_VERIFICATION" ]]; then
            wget --no-check-certificate "$url" -O "$output"
        else
            wget "$url" -O "$output"
        fi
    else
        echo "Neither curl nor wget is installed. Cannot download file."
        return 1
    fi
}

echo "Adding phatcrack-agent user..."
useradd --system --create-home --home-dir /opt/phatcrack-agent phatcrack-agent || true

echo "Adding phatcrack-agent user to video group (might error if it doesn't exist)"
usermod -aG video phatcrack-agent || true

cd /opt/phatcrack-agent

echo "Downloading binaries"
download_file $PHATCRACK_HOST/agent-assets/hashcat.7z hashcat.7z
download_file $PHATCRACK_HOST/agent-assets/phatcrack-agent phatcrack-agent

7z x hashcat.7z
rm hashcat.7z
mv hashcat-* hashcat

chown -R phatcrack-agent:phatcrack-agent ./hashcat
chmod +x ./phatcrack-agent

local tls_arg=""
if [[ -n "$DISABLE_TLS_VERIFICATION" ]]; then
    tls_arg="-disable-tls-verification"
fi

./phatcrack-agent install -defaults -api-endpoint $PHATCRACK_HOST/api/v1 -auth-key $PHATCRACK_API_KEY $tls_arg

systemctl daemon-reload
systemctl enable --now phatcrack-agent
#!/bin/bash

# Format v1.2.3
RELEASE_VERSION=$1

sed -i "s@lachlan2k/phatcrack/releases/download/.*/@lachlan2k/phatcrack/releases/download/${RELEASE_VERSION}/@" README.md
sed -i "s/export PHATCRACK_VERSION_TAG=.*/export PHATCRACK_VERSION_TAG=${RELEASE_VERSION}/" install_server.sh

git commit --allow-empty -am "chore: Version ${RELEASE_VERSION}"

git push
git tag ${RELEASE_VERSION}
git push --tags
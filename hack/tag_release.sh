#!/bin/bash

# Format v1.2.3
RELEASE_VERSION=$1

sed -i "s@lachlan2k/phatcrack/releases/download/.*/@lachlan2k/phatcrack/releases/download/${RELEASE_VERSION}/@" README.md
sed -i "s@lachlan2k/phatcrack/releases/download/.*/@lachlan2k/phatcrack/releases/download/${RELEASE_VERSION}/@" deploy_server.sh

git commit --allow-empty -am "chore: Version ${RELEASE_VERSION}"

git push
git tag ${RELEASE_VERSION}
git push --tags
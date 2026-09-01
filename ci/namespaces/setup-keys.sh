#!/bin/bash
set -e

# Generate WireGuard keys for node-a
mkdir -p ci/namespaces/node-a
wg genkey | tee ci/namespaces/node-a/node-a.key | wg pubkey > ci/namespaces/node-a/node-a.pub

# Generate WireGuard keys for node-b
mkdir -p ci/namespaces/node-b
wg genkey | tee ci/namespaces/node-b/node-b.key | wg pubkey > ci/namespaces/node-b/node-b.pub

# Copy public keys to opposite nodes for peer configuration
cp ci/namespaces/node-a/node-a.pub ci/namespaces/node-b/node-a.pub
cp ci/namespaces/node-b/node-b.pub ci/namespaces/node-a/node-b.pub

# Set proper permissions
chmod 600 ci/namespaces/node-a/node-a.key ci/namespaces/node-b/node-b.key
chmod 644 ci/namespaces/node-a/*.pub ci/namespaces/node-b/*.pub
chmod +x ci/namespaces/node-a/start.sh ci/namespaces/node-b/start.sh

echo "WireGuard keys generated successfully!"
echo "Node-A public key: $(cat ci/namespaces/node-a/node-a.pub)"
echo "Node-B public key: $(cat ci/namespaces/node-b/node-b.pub)"

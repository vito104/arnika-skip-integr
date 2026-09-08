#!/bin/bash
set -e

echo "====== Namespace Integration Test - PSK Verification ====="

MAX_ATTEMPTS=5
ATTEMPT=0
PSK_A=""
PSK_B=""

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    PSK_A=$(docker exec clab-arnika-ci-test-netns-node-a ip netns exec arnika-ns-a wg show wg0 preshared-keys 2>/dev/null | awk '{print $2}')
    PSK_B=$(docker exec clab-arnika-ci-test-netns-node-b ip netns exec arnika-ns-b wg show wg0 preshared-keys 2>/dev/null | awk '{print $2}')
    
    if [ -n "$PSK_A" ] && [ "$PSK_A" != "(none)" ] && \
       [ -n "$PSK_B" ] && [ "$PSK_B" != "(none)" ] && \
       [ "$PSK_A" = "$PSK_B" ]; then
        break
    fi
    
    ATTEMPT=$((ATTEMPT + 1))
    if [ $ATTEMPT -lt $MAX_ATTEMPTS ]; then
        sleep 2
    fi
done

if [ "$PSK_A" = "$PSK_B" ] && [ -n "$PSK_A" ] && [ "$PSK_A" != "(none)" ]; then
    echo "SUCCESS: Namespace PSKs match"
    exit 0
else
    echo "FAILED: Namespace PSKs do not match"
    echo "Node-A PSK: ${PSK_A}"
    echo "Node-B PSK: ${PSK_B}"
    exit 1
fi

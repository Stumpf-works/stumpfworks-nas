#!/bin/bash
# Setup firewall rules for bridge networking
# This script is executed at boot to ensure containers/VMs can communicate through bridges

# Allow forwarding for all bridge interfaces
for bridge in $(ip -o link show type bridge | awk -F': ' '{print $2}'); do
    echo "Setting up firewall rules for bridge: $bridge"

    # Allow traffic within the bridge
    iptables -C FORWARD -i "$bridge" -o "$bridge" -j ACCEPT 2>/dev/null || \
        iptables -I FORWARD -i "$bridge" -o "$bridge" -j ACCEPT

    # Allow traffic from bridge to external
    iptables -C FORWARD -i "$bridge" -j ACCEPT 2>/dev/null || \
        iptables -I FORWARD -i "$bridge" -j ACCEPT

    # Allow traffic to bridge from external
    iptables -C FORWARD -o "$bridge" -j ACCEPT 2>/dev/null || \
        iptables -I FORWARD -o "$bridge" -j ACCEPT

    echo "  ✓ Firewall rules configured for $bridge"
done

echo "Bridge firewall setup complete"

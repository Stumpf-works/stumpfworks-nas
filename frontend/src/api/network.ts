import client, { ApiResponse } from './client';

// Types
export interface NetworkInterface {
  name: string;
  index: number;
  hardwareAddr: string;
  flags: string[];
  mtu: number;
  addresses: string[];
  isUp: boolean;
  speed: string;
  type: string;
}

export interface InterfaceStats {
  name: string;
  rxBytes: number;
  txBytes: number;
  rxPackets: number;
  txPackets: number;
  rxErrors: number;
  txErrors: number;
  rxDropped: number;
  txDropped: number;
}

export interface Route {
  destination: string;
  gateway: string;
  genmask: string;
  flags: string;
  metric: number;
  iface: string;
}

export interface DNSConfig {
  nameservers: string[];
  searchDomains: string[];
}

export interface FirewallRule {
  number: number;
  action: string;
  from: string;
  to: string;
  protocol?: string;
  port?: string;
  description?: string;
}

export interface FirewallStatus {
  enabled: boolean;
  defaultIncoming: string;
  defaultOutgoing: string;
  defaultRouted: string;
  rules: FirewallRule[];
}

export interface NetworkBridge {
  id: string;
  name: string;
  description: string;
  ports: string;
  // IPv4
  ip_address?: string;
  gateway?: string;
  // IPv6 (Proxmox-style)
  ipv6_address?: string;
  ipv6_gateway?: string;
  // Bridge settings
  vlan_aware?: boolean;
  autostart: boolean;
  status: string;
  last_error?: string;
  // Pending changes
  has_pending_changes: boolean;
  pending_ports?: string;
  pending_ip_address?: string;
  pending_gateway?: string;
  pending_ipv6_address?: string;
  pending_ipv6_gateway?: string;
  pending_vlan_aware?: boolean;
  created_at: string;
  updated_at: string;
}

export interface PendingNetworkChange {
  id: string;
  change_type: string;
  action: string;
  resource_id: string;
  current_config?: string;
  pending_config: string;
  description?: string;
  created_by?: string;
  priority: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface PendingChangesResponse {
  has_pending: boolean;
  count: number;
  changes: PendingNetworkChange[];
}

export interface DiagnosticResult {
  command: string;
  output: string;
  success: boolean;
  error?: string;
}

// API
export const networkApi = {
  // Interfaces
  async listInterfaces(): Promise<ApiResponse<NetworkInterface[]>> {
    const response = await client.get('/network/interfaces');
    return response.data;
  },

  async getInterfaceStats(): Promise<ApiResponse<InterfaceStats[]>> {
    const response = await client.get('/network/interfaces/stats');
    return response.data;
  },

  async setInterfaceState(name: string, state: 'up' | 'down'): Promise<ApiResponse<any>> {
    const response = await client.post(`/network/interfaces/${name}/state`, { state });
    return response.data;
  },

  async configureInterface(
    name: string,
    mode: 'static' | 'dhcp',
    address?: string,
    netmask?: string,
    gateway?: string
  ): Promise<ApiResponse<any>> {
    const response = await client.post(`/network/interfaces/${name}/configure`, {
      mode,
      address,
      netmask,
      gateway,
    });
    return response.data;
  },

  // Routes
  async getRoutes(): Promise<ApiResponse<Route[]>> {
    const response = await client.get('/network/routes');
    return response.data;
  },

  async addRoute(destination: string, gateway?: string, iface?: string, metric?: number): Promise<ApiResponse<any>> {
    const response = await client.post('/network/routes', {
      destination,
      gateway,
      interface: iface,
      metric,
    });
    return response.data;
  },

  async deleteRoute(destination: string, gateway?: string, iface?: string): Promise<ApiResponse<any>> {
    const response = await client.delete('/network/routes', {
      data: {
        destination,
        gateway,
        interface: iface,
      },
    });
    return response.data;
  },

  // DNS
  async getDNS(): Promise<ApiResponse<DNSConfig>> {
    const response = await client.get('/network/dns');
    return response.data;
  },

  async setDNS(nameservers: string[], searchDomains: string[]): Promise<ApiResponse<any>> {
    const response = await client.post('/network/dns', { nameservers, searchDomains });
    return response.data;
  },

  // Firewall
  async getFirewallStatus(): Promise<ApiResponse<FirewallStatus>> {
    const response = await client.get('/network/firewall');
    return response.data;
  },

  async setFirewallState(enabled: boolean): Promise<ApiResponse<any>> {
    const response = await client.post('/network/firewall/state', { enabled });
    return response.data;
  },

  async addFirewallRule(
    action: string,
    port: string,
    protocol: string,
    from: string,
    to: string
  ): Promise<ApiResponse<any>> {
    const response = await client.post('/network/firewall/rules', { action, port, protocol, from, to });
    return response.data;
  },

  async deleteFirewallRule(number: number): Promise<ApiResponse<any>> {
    const response = await client.delete(`/network/firewall/rules/${number}`);
    return response.data;
  },

  async setDefaultPolicy(direction: string, policy: string): Promise<ApiResponse<any>> {
    const response = await client.post('/network/firewall/default', { direction, policy });
    return response.data;
  },

  async resetFirewall(): Promise<ApiResponse<any>> {
    const response = await client.post('/network/firewall/reset', {});
    return response.data;
  },

  // Diagnostics
  async ping(host: string, count: number = 4): Promise<ApiResponse<DiagnosticResult>> {
    const response = await client.post('/network/diagnostics/ping', { host, count });
    return response.data;
  },

  async traceroute(host: string): Promise<ApiResponse<DiagnosticResult>> {
    const response = await client.post('/network/diagnostics/traceroute', { host });
    return response.data;
  },

  async netstat(options: string = ''): Promise<ApiResponse<DiagnosticResult>> {
    const response = await client.post('/network/diagnostics/netstat', { options });
    return response.data;
  },

  // Wake-on-LAN
  async wakeOnLAN(macAddress: string): Promise<ApiResponse<any>> {
    const response = await client.post('/network/wol', { macAddress });
    return response.data;
  },

  // Bridge management
  async listBridges(): Promise<ApiResponse<string[]>> {
    const response = await client.get('/network/bridges');
    return response.data;
  },

  async createBridge(name: string, ports: string[]): Promise<ApiResponse<any>> {
    const response = await client.post('/network/bridges', { name, ports });
    return response.data;
  },

  async deleteBridge(name: string): Promise<ApiResponse<any>> {
    const response = await client.delete(`/network/bridges/${name}`);
    return response.data;
  },

  async attachPortToBridge(bridgeName: string, port: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/network/bridges/${bridgeName}/attach`, { port });
    return response.data;
  },

  async detachPortFromBridge(port: string): Promise<ApiResponse<any>> {
    const response = await client.post(`/network/bridges/detach`, { port });
    return response.data;
  },

  // Proxmox-style pending changes workflow
  async getPendingChanges(): Promise<ApiResponse<PendingChangesResponse>> {
    const response = await client.get('/network/pending-changes');
    return response.data;
  },

  async createBridgeWithPendingChanges(
    name: string,
    description: string,
    ports: string[],
    ipAddress?: string,
    gateway?: string,
    ipv6Address?: string,
    ipv6Gateway?: string,
    vlanAware: boolean = false,
    autostart: boolean = true
  ): Promise<ApiResponse<NetworkBridge>> {
    const response = await client.post('/network/bridges/pending', {
      name,
      description,
      ports,
      ip_address: ipAddress,
      gateway,
      ipv6_address: ipv6Address,
      ipv6_gateway: ipv6Gateway,
      vlan_aware: vlanAware,
      autostart,
    });
    return response.data;
  },

  async updateBridgeWithPendingChanges(
    name: string,
    description?: string,
    ports?: string[],
    ipAddress?: string,
    gateway?: string,
    ipv6Address?: string,
    ipv6Gateway?: string,
    vlanAware?: boolean,
    autostart?: boolean
  ): Promise<ApiResponse<NetworkBridge>> {
    const response = await client.put(`/network/bridges/${name}/pending`, {
      description,
      ports,
      ip_address: ipAddress,
      gateway,
      ipv6_address: ipv6Address,
      ipv6_gateway: ipv6Gateway,
      vlan_aware: vlanAware,
      autostart,
    });
    return response.data;
  },

  async updateInterfaceWithPendingChanges(
    name: string,
    ipAddress?: string,
    gateway?: string,
    ipv6Address?: string,
    ipv6Gateway?: string,
    autostart: boolean = true,
    comment?: string
  ): Promise<ApiResponse<any>> {
    const response = await client.put(`/network/interfaces/${name}/pending`, {
      ip_address: ipAddress,
      gateway,
      ipv6_address: ipv6Address,
      ipv6_gateway: ipv6Gateway,
      autostart,
      comment,
    });
    return response.data;
  },

  async applyPendingChanges(): Promise<ApiResponse<any>> {
    const response = await client.post('/network/apply-changes', {});
    return response.data;
  },

  async discardPendingChanges(changeId?: string): Promise<ApiResponse<any>> {
    const response = await client.post('/network/discard-changes', { change_id: changeId });
    return response.data;
  },

  async rollbackToSnapshot(snapshotId: string): Promise<ApiResponse<any>> {
    const response = await client.post('/network/rollback', { snapshot_id: snapshotId });
    return response.data;
  },

  async getStoredBridges(): Promise<ApiResponse<NetworkBridge[]>> {
    const response = await client.get('/network/bridges/stored');
    return response.data;
  },
};

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
    return client.get('/network/interfaces');
  },

  async getInterfaceStats(): Promise<ApiResponse<InterfaceStats[]>> {
    return client.get('/network/interfaces/stats');
  },

  async setInterfaceState(name: string, state: 'up' | 'down'): Promise<ApiResponse<any>> {
    return client.post(`/network/interfaces/${name}/state`, { state });
  },

  async configureInterface(
    name: string,
    mode: 'static' | 'dhcp',
    address?: string,
    netmask?: string,
    gateway?: string
  ): Promise<ApiResponse<any>> {
    return client.post(`/network/interfaces/${name}/configure`, {
      mode,
      address,
      netmask,
      gateway,
    });
  },

  // Routes
  async getRoutes(): Promise<ApiResponse<Route[]>> {
    return client.get('/network/routes');
  },

  // DNS
  async getDNS(): Promise<ApiResponse<DNSConfig>> {
    return client.get('/network/dns');
  },

  async setDNS(nameservers: string[], searchDomains: string[]): Promise<ApiResponse<any>> {
    return client.post('/network/dns', { nameservers, searchDomains });
  },

  // Firewall
  async getFirewallStatus(): Promise<ApiResponse<FirewallStatus>> {
    return client.get('/network/firewall');
  },

  async setFirewallState(enabled: boolean): Promise<ApiResponse<any>> {
    return client.post('/network/firewall/state', { enabled });
  },

  async addFirewallRule(
    action: string,
    port: string,
    protocol: string,
    from: string,
    to: string
  ): Promise<ApiResponse<any>> {
    return client.post('/network/firewall/rules', { action, port, protocol, from, to });
  },

  async deleteFirewallRule(number: number): Promise<ApiResponse<any>> {
    return client.delete(`/network/firewall/rules/${number}`);
  },

  async setDefaultPolicy(direction: string, policy: string): Promise<ApiResponse<any>> {
    return client.post('/network/firewall/default', { direction, policy });
  },

  async resetFirewall(): Promise<ApiResponse<any>> {
    return client.post('/network/firewall/reset', {});
  },

  // Diagnostics
  async ping(host: string, count: number = 4): Promise<ApiResponse<DiagnosticResult>> {
    return client.post('/network/diagnostics/ping', { host, count });
  },

  async traceroute(host: string): Promise<ApiResponse<DiagnosticResult>> {
    return client.post('/network/diagnostics/traceroute', { host });
  },

  async netstat(options: string = ''): Promise<ApiResponse<DiagnosticResult>> {
    return client.post('/network/diagnostics/netstat', { options });
  },

  // Wake-on-LAN
  async wakeOnLAN(macAddress: string): Promise<ApiResponse<any>> {
    return client.post('/network/wol', { macAddress });
  },
};

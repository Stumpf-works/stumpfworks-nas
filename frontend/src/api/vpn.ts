import client, { ApiResponse } from './client';

// VPN Protocol Types
export type VPNProtocol = 'wireguard' | 'openvpn' | 'pptp' | 'l2tp';

export interface ProtocolStatus {
  protocol: string;
  installed: boolean;
  enabled: boolean;
  running: boolean;
  connections: number;
  error?: string;
}

export interface WireGuardPeer {
  id: string;
  name: string;
  public_key: string;
  allowed_ips: string;
  endpoint?: string;
  latest_handshake?: string;
  bytes_received: number;
  bytes_sent: number;
  enabled: boolean;
}

export interface OpenVPNCertificate {
  id: string;
  common_name: string;
  serial_number: string;
  valid_from: string;
  valid_to: string;
  status: 'valid' | 'revoked' | 'expired';
}

export const vpnApi = {
  // Get status of all VPN protocols
  getStatus: async (): Promise<ApiResponse<ProtocolStatus[]>> => {
    const response = await client.get<ApiResponse<ProtocolStatus[]>>('/vpn/status');
    return response.data;
  },

  // Get status of a specific protocol
  getProtocolStatus: async (protocol: VPNProtocol): Promise<ApiResponse<ProtocolStatus>> => {
    const response = await client.get<ApiResponse<ProtocolStatus>>(`/vpn/protocols/${protocol}/status`);
    return response.data;
  },

  // Install a protocol (installs packages and initializes)
  installProtocol: async (protocol: VPNProtocol): Promise<ApiResponse<{ message: string; protocol: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; protocol: string }>>(`/vpn/protocols/${protocol}/install`);
    return response.data;
  },

  // Enable a protocol (start services)
  enableProtocol: async (protocol: VPNProtocol): Promise<ApiResponse<{ message: string; protocol: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; protocol: string }>>(`/vpn/protocols/${protocol}/enable`);
    return response.data;
  },

  // Disable a protocol (stop services)
  disableProtocol: async (protocol: VPNProtocol): Promise<ApiResponse<{ message: string; protocol: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; protocol: string }>>(`/vpn/protocols/${protocol}/disable`);
    return response.data;
  },

  // WireGuard-specific endpoints
  wireguard: {
    // Get all peers
    getPeers: async (): Promise<ApiResponse<WireGuardPeer[]>> => {
      const response = await client.get<ApiResponse<WireGuardPeer[]>>('/vpn/wireguard/peers');
      return response.data;
    },

    // Create a new peer
    createPeer: async (name: string, allowedIPs: string): Promise<ApiResponse<WireGuardPeer>> => {
      const response = await client.post<ApiResponse<WireGuardPeer>>('/vpn/wireguard/peers', {
        name,
        allowed_ips: allowedIPs,
      });
      return response.data;
    },

    // Delete a peer
    deletePeer: async (peerId: string): Promise<ApiResponse<{ message: string }>> => {
      const response = await client.delete<ApiResponse<{ message: string }>>(`/vpn/wireguard/peers/${peerId}`);
      return response.data;
    },

    // Get peer configuration
    getPeerConfig: async (peerId: string): Promise<ApiResponse<{ config: string }>> => {
      const response = await client.get<ApiResponse<{ config: string }>>(`/vpn/wireguard/peers/${peerId}/config`);
      return response.data;
    },

    // Get peer QR code
    getPeerQRCode: async (peerId: string): Promise<ApiResponse<{ qrcode: string }>> => {
      const response = await client.get<ApiResponse<{ qrcode: string }>>(`/vpn/wireguard/peers/${peerId}/qrcode`);
      return response.data;
    },
  },

  // OpenVPN-specific endpoints
  openvpn: {
    // Get all certificates
    getCertificates: async (): Promise<ApiResponse<OpenVPNCertificate[]>> => {
      const response = await client.get<ApiResponse<OpenVPNCertificate[]>>('/vpn/openvpn/certificates');
      return response.data;
    },

    // Create a new certificate
    createCertificate: async (commonName: string): Promise<ApiResponse<OpenVPNCertificate>> => {
      const response = await client.post<ApiResponse<OpenVPNCertificate>>('/vpn/openvpn/certificates', {
        common_name: commonName,
      });
      return response.data;
    },

    // Revoke a certificate
    revokeCertificate: async (certId: string): Promise<ApiResponse<{ message: string }>> => {
      const response = await client.post<ApiResponse<{ message: string }>>(`/vpn/openvpn/certificates/${certId}/revoke`);
      return response.data;
    },
  },
};

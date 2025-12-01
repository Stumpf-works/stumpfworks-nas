// VPN Protocol Types
export type VPNProtocol = 'wireguard' | 'openvpn' | 'pptp' | 'l2tp';

export interface VPNUser {
  id: string;
  username: string;
  email: string;
  createdAt: string;
  lastConnection?: string;
  protocols: {
    wireguard: boolean;
    openvpn: boolean;
    pptp: boolean;
    l2tp: boolean;
  };
  enabled: boolean;
}

export interface VPNConnection {
  id: string;
  userId: string;
  protocol: VPNProtocol;
  ipAddress: string;
  connectedAt: string;
  bytesReceived: number;
  bytesSent: number;
}

// WireGuard Types
export interface WireGuardPeer {
  id: string;
  name: string;
  publicKey: string;
  privateKey?: string;
  allowedIPs: string;
  endpoint?: string;
  latestHandshake?: string;
  bytesReceived: number;
  bytesSent: number;
  enabled: boolean;
  userId?: string;
}

export interface WireGuardServerConfig {
  enabled: boolean;
  running: boolean;
  publicKey: string;
  privateKey: string;
  listenPort: number;
  endpoint: string;
  subnet: string;
  dns: string;
  peers: WireGuardPeer[];
}

// OpenVPN Types
export interface OpenVPNClientCertificate {
  id: string;
  commonName: string;
  serialNumber: string;
  validFrom: string;
  validTo: string;
  status: 'valid' | 'revoked' | 'expired';
  userId?: string;
}

export interface OpenVPNConnection {
  id: string;
  commonName: string;
  realAddress: string;
  virtualAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

export interface OpenVPNServerConfig {
  enabled: boolean;
  running: boolean;
  protocol: 'UDP' | 'TCP';
  port: number;
  subnet: string;
  cipher: string;
  auth: string;
  compression: string;
  tlsVersion: string;
  certificates: OpenVPNClientCertificate[];
  connections: OpenVPNConnection[];
}

// PPTP Types
export interface PPTPConnection {
  id: string;
  username: string;
  ipAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

export interface PPTPServerConfig {
  enabled: boolean;
  running: boolean;
  port: number;
  subnet: string;
  encryption: 'MPPE-128' | 'MPPE-40';
  authentication: 'MS-CHAPv2' | 'CHAP' | 'PAP';
  connections: PPTPConnection[];
}

// L2TP/IPsec Types
export interface L2TPConnection {
  id: string;
  username: string;
  ipAddress: string;
  connectedSince: string;
  bytesReceived: number;
  bytesSent: number;
}

export interface L2TPServerConfig {
  enabled: boolean;
  running: boolean;
  port: number;
  ipsecPort: number;
  subnet: string;
  psk: string;
  encryption: 'AES-256' | 'AES-192' | 'AES-128' | '3DES';
  authentication: 'SHA2-256' | 'SHA2-512' | 'SHA1';
  natTraversal: boolean;
  connections: L2TPConnection[];
}

// General Settings Types
export interface VPNGeneralSettings {
  defaultInterface: string;
  accountSource: 'local' | 'ldap' | 'ad' | 'radius';
  enableLogging: boolean;
  logLevel: 'error' | 'warning' | 'info' | 'debug';
  maxConcurrentConnections: number;
  connectionTimeout: number;
  enableIPv6: boolean;
  dnsServers: string;
  defaultGateway: string;
  enableNAT: boolean;
  forwardingRules: string;
}

// Statistics Types
export interface VPNStatistics {
  totalConnections: number;
  activeProtocols: number;
  totalBytesIn: number;
  totalBytesOut: number;
  connectionsByProtocol: {
    wireguard: number;
    openvpn: number;
    pptp: number;
    l2tp: number;
  };
}

// API Response Types
export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface VPNServerStatus {
  wireguard: WireGuardServerConfig;
  openvpn: OpenVPNServerConfig;
  pptp: PPTPServerConfig;
  l2tp: L2TPServerConfig;
  statistics: VPNStatistics;
  settings: VPNGeneralSettings;
}

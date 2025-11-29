import { ApiResponse } from './client';

// VM Types
export interface VM {
  uuid: string;
  name: string;
  state: string;
  memory: number;
  vcpus: number;
  disk_size: number;
  created_at: string;
  autostart: boolean;
}

export interface VMCreateRequest {
  name: string;
  memory: number;
  vcpus: number;
  disk_size: number;
  os_type?: string;
  os_variant?: string;
  iso_path?: string;
  network?: string;
  autostart?: boolean;
}

export const vmsApi = {
  // List all VMs
  listVMs: async (): Promise<ApiResponse<VM[]>> => {
    const response = await fetch('/api/v1/vms/');
    return response.json();
  },

  // Get VM details
  getVM: async (vmId: string): Promise<ApiResponse<VM>> => {
    const response = await fetch(`/api/v1/vms/${encodeURIComponent(vmId)}`);
    return response.json();
  },

  // Create a new VM
  createVM: async (data: VMCreateRequest): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await fetch('/api/v1/vms/', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Start a VM
  startVM: async (vmId: string): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await fetch(`/api/v1/vms/${encodeURIComponent(vmId)}/start`, {
      method: 'POST',
    });
    return response.json();
  },

  // Stop a VM
  stopVM: async (vmId: string, force: boolean = false): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await fetch(`/api/v1/vms/${encodeURIComponent(vmId)}/stop?force=${force}`, {
      method: 'POST',
    });
    return response.json();
  },

  // Delete a VM
  deleteVM: async (vmId: string, deleteDisks: boolean = false): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await fetch(`/api/v1/vms/${encodeURIComponent(vmId)}?delete_disks=${deleteDisks}`, {
      method: 'DELETE',
    });
    return response.json();
  },

  // Get VNC port for VM console access
  getVNCPort: async (vmId: string): Promise<ApiResponse<{ vm_id: string; port: number }>> => {
    const response = await fetch(`/api/v1/vms/${encodeURIComponent(vmId)}/vnc`);
    return response.json();
  },
};

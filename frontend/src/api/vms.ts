import client, { ApiResponse } from './client';

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
    const response = await client.get<ApiResponse<VM[]>>('/vms/');
    return response.data;
  },

  // Get VM details
  getVM: async (vmId: string): Promise<ApiResponse<VM>> => {
    const response = await client.get<ApiResponse<VM>>(`/vms/${encodeURIComponent(vmId)}`);
    return response.data;
  },

  // Create a new VM
  createVM: async (data: VMCreateRequest): Promise<ApiResponse<{ message: string; name: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; name: string }>>('/vms/', data);
    return response.data;
  },

  // Start a VM
  startVM: async (vmId: string): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; vm_id: string }>>(`/vms/${encodeURIComponent(vmId)}/start`);
    return response.data;
  },

  // Stop a VM
  stopVM: async (vmId: string, force: boolean = false): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await client.post<ApiResponse<{ message: string; vm_id: string }>>(`/vms/${encodeURIComponent(vmId)}/stop?force=${force}`);
    return response.data;
  },

  // Delete a VM
  deleteVM: async (vmId: string, deleteDisks: boolean = false): Promise<ApiResponse<{ message: string; vm_id: string }>> => {
    const response = await client.delete<ApiResponse<{ message: string; vm_id: string }>>(`/vms/${encodeURIComponent(vmId)}?delete_disks=${deleteDisks}`);
    return response.data;
  },

  // Get VNC port for VM console access
  getVNCPort: async (vmId: string): Promise<ApiResponse<{ vm_id: string; port: number }>> => {
    const response = await client.get<ApiResponse<{ vm_id: string; port: number }>>(`/vms/${encodeURIComponent(vmId)}/vnc`);
    return response.data;
  },
};

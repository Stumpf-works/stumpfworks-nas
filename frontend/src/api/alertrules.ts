import client, { ApiResponse } from './client';

export interface AlertRule {
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  description: string;
  enabled: boolean;
  metricType: 'cpu' | 'memory' | 'disk' | 'network' | 'health' | 'temperature' | 'iops';
  condition: 'gt' | 'lt' | 'eq' | 'gte' | 'lte';
  threshold: number;
  duration: number; // seconds
  cooldownMins: number;
  severity: 'info' | 'warning' | 'critical';
  notifyEmail: boolean;
  notifyWebhook: boolean;
  notifyChannels: string;
  lastTriggered?: string;
  triggerCount: number;
  isActive: boolean;
  activatedAt?: string;
}

export interface AlertRuleExecution {
  id: number;
  createdAt: string;
  ruleId: number;
  rule?: AlertRule;
  metricValue: number;
  threshold: number;
  triggered: boolean;
  acknowledged: boolean;
  acknowledgedAt?: string;
  acknowledgedBy?: string;
  acknowledgeNote?: string;
  notificationsSent: boolean;
  message: string;
}

export const alertRulesApi = {
  // Get all alert rules
  listRules: async () => {
    const response = await client.get<ApiResponse<AlertRule[]>>('/alert-rules');
    return response.data;
  },

  // Get single alert rule
  getRule: async (id: number) => {
    const response = await client.get<ApiResponse<AlertRule>>(`/alert-rules/${id}`);
    return response.data;
  },

  // Create new alert rule
  createRule: async (rule: Omit<AlertRule, 'id' | 'createdAt' | 'updatedAt' | 'lastTriggered' | 'triggerCount' | 'isActive' | 'activatedAt'>) => {
    const response = await client.post<ApiResponse<AlertRule>>('/alert-rules', rule);
    return response.data;
  },

  // Update alert rule
  updateRule: async (id: number, rule: Partial<AlertRule>) => {
    const response = await client.put<ApiResponse<AlertRule>>(`/alert-rules/${id}`, rule);
    return response.data;
  },

  // Delete alert rule
  deleteRule: async (id: number) => {
    const response = await client.delete<ApiResponse<{ message: string }>>(`/alert-rules/${id}`);
    return response.data;
  },

  // Get executions for a specific rule
  getExecutions: async (ruleId: number, limit = 50) => {
    const response = await client.get<ApiResponse<AlertRuleExecution[]>>(
      `/alert-rules/${ruleId}/executions?limit=${limit}`
    );
    return response.data;
  },

  // Get recent executions across all rules
  getRecentExecutions: async (limit = 100) => {
    const response = await client.get<ApiResponse<AlertRuleExecution[]>>(
      `/alert-rules/executions/recent?limit=${limit}`
    );
    return response.data;
  },

  // Acknowledge an execution
  acknowledgeExecution: async (executionId: number, acknowledgedBy: string, note?: string) => {
    const response = await client.post<ApiResponse<{ message: string }>>(
      `/alert-rules/executions/${executionId}/acknowledge`,
      { acknowledgedBy, note }
    );
    return response.data;
  },
};

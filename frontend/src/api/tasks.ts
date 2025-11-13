import client from './client';

export interface ScheduledTask {
  id?: number;
  name: string;
  description?: string;
  taskType: string;
  cronExpression: string;
  enabled: boolean;
  lastRun?: string;
  nextRun?: string;
  lastStatus?: string;
  lastError?: string;
  config?: string;
  timeoutSeconds: number;
  retryOnFailure: boolean;
  runCount?: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface TaskExecution {
  id: number;
  taskId: number;
  startedAt: string;
  completedAt?: string;
  duration: number;
  status: string;
  output?: string;
  error?: string;
  triggeredBy: string;
  createdAt: string;
}

export interface TaskListResponse {
  tasks: ScheduledTask[];
  total: number;
  offset: number;
  limit: number;
}

export interface TaskExecutionsResponse {
  executions: TaskExecution[];
  total: number;
  offset: number;
  limit: number;
}

export interface CronValidationRequest {
  expression: string;
}

export interface CronValidationResponse {
  valid: boolean;
  error?: string;
  nextRuns?: string[];
}

export const tasksApi = {
  /**
   * List all scheduled tasks with pagination
   */
  listTasks: async (offset = 0, limit = 50): Promise<TaskListResponse> => {
    const response = await client.get('/api/v1/tasks', {
      params: { offset, limit },
    });
    return response.data.data;
  },

  /**
   * Get a specific task by ID
   */
  getTask: async (id: number): Promise<ScheduledTask> => {
    const response = await client.get(`/api/v1/tasks/${id}`);
    return response.data.data;
  },

  /**
   * Create a new scheduled task
   */
  createTask: async (task: Omit<ScheduledTask, 'id'>): Promise<ScheduledTask> => {
    const response = await client.post('/api/v1/tasks', task);
    return response.data.data;
  },

  /**
   * Update an existing task
   */
  updateTask: async (id: number, task: ScheduledTask): Promise<ScheduledTask> => {
    const response = await client.put(`/api/v1/tasks/${id}`, task);
    return response.data.data;
  },

  /**
   * Delete a task
   */
  deleteTask: async (id: number): Promise<void> => {
    await client.delete(`/api/v1/tasks/${id}`);
  },

  /**
   * Run a task immediately
   */
  runTaskNow: async (id: number): Promise<void> => {
    await client.post(`/api/v1/tasks/${id}/run`);
  },

  /**
   * Get execution history for a task
   */
  getTaskExecutions: async (
    id: number,
    offset = 0,
    limit = 50
  ): Promise<TaskExecutionsResponse> => {
    const response = await client.get(`/api/v1/tasks/${id}/executions`, {
      params: { offset, limit },
    });
    return response.data.data;
  },

  /**
   * Validate a cron expression
   */
  validateCron: async (expression: string): Promise<CronValidationResponse> => {
    const response = await client.post('/api/v1/tasks/validate-cron', {
      expression,
    });
    return response.data.data;
  },
};

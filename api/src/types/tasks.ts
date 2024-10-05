export interface Task {
  id: number;
  function_name: string;
  priority: number;
  max_retries: number;
}

export interface TaskResponse {
  success: boolean;
  message: string;
}

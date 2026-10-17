export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
  error?: any;
}

export interface PaginatedResponse<T = any> {
  success: boolean;
  message: string;
  data: T[];
  pagination: {
    current_page: number;
    limit: number;
    total_items: number;
    total_pages: number;
  };
}

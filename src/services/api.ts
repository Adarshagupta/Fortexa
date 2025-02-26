import axios, { AxiosRequestConfig } from 'axios';
import { ApiResponse, Payment, Settlement, Merchant, DashboardStats } from '../types';

// Create an axios instance with default configuration
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:3001/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for API calls
api.interceptors.request.use(
  (config) => {
    // Get the API key from localStorage
    const apiKey = localStorage.getItem('apiKey');
    
    if (apiKey) {
      config.headers['X-API-Key'] = apiKey;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for API calls
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access (logout user, redirect to login, etc.)
      localStorage.removeItem('apiKey');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Common error handler
const handleApiError = (error: any): any => {
  if (error.response) {
    // The request was made and the server responded with a status code
    // that falls out of the range of 2xx
    console.error('API Error Response:', error.response.data);
    return error.response.data;
  } else if (error.request) {
    // The request was made but no response was received
    console.error('API Error Request:', error.request);
    return {
      success: false,
      message: 'No response received from server',
    };
  } else {
    // Something happened in setting up the request that triggered an Error
    console.error('API Error:', error.message);
    return {
      success: false,
      message: error.message,
    };
  }
};

// Authentication API calls
export const authApi = {
  login: async (email: string, password: string): Promise<ApiResponse<{ apiKey: string }>> => {
    try {
      const response = await api.post<ApiResponse<{ apiKey: string }>>('/auth/login', { email, password });
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

// Payments API calls
export const paymentsApi = {
  getPayments: async (page = 1, limit = 10, filters = {}): Promise<ApiResponse<{ payments: Payment[], total: number }>> => {
    try {
      const response = await api.get<ApiResponse<{ payments: Payment[], total: number }>>('/payments', {
        params: { page, limit, ...filters },
      });
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  getPaymentById: async (id: string): Promise<ApiResponse<Payment>> => {
    try {
      const response = await api.get<ApiResponse<Payment>>(`/payments/${id}`);
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  createPayment: async (paymentData: Partial<Payment>): Promise<ApiResponse<Payment>> => {
    try {
      const response = await api.post<ApiResponse<Payment>>('/payments', paymentData);
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  refundPayment: async (id: string, amount: number, reason?: string): Promise<ApiResponse<Payment>> => {
    try {
      const response = await api.post<ApiResponse<Payment>>(`/payments/${id}/refund`, { 
        amount, 
        reason,
      });
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

// Settlements API calls
export const settlementsApi = {
  getSettlements: async (page = 1, limit = 10, filters = {}): Promise<ApiResponse<{ settlements: Settlement[], total: number }>> => {
    try {
      const response = await api.get<ApiResponse<{ settlements: Settlement[], total: number }>>('/settlements', {
        params: { page, limit, ...filters },
      });
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  getSettlementById: async (id: string): Promise<ApiResponse<Settlement>> => {
    try {
      const response = await api.get<ApiResponse<Settlement>>(`/settlements/${id}`);
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

// Merchant API calls
export const merchantApi = {
  getMerchantProfile: async (): Promise<ApiResponse<Merchant>> => {
    try {
      const response = await api.get<ApiResponse<Merchant>>('/merchants/profile');
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
  
  updateMerchantProfile: async (profileData: Partial<Merchant>): Promise<ApiResponse<Merchant>> => {
    try {
      const response = await api.put<ApiResponse<Merchant>>('/merchants/profile', profileData);
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

// Dashboard API calls
export const dashboardApi = {
  getStats: async (): Promise<ApiResponse<DashboardStats>> => {
    try {
      const response = await api.get<ApiResponse<DashboardStats>>('/dashboard/stats');
      return response.data;
    } catch (error) {
      return handleApiError(error);
    }
  },
};

export default api; 
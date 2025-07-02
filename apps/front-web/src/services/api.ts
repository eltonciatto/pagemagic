import axios from 'axios';

// Configuração base do axios
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para adicionar token de autenticação
apiClient.interceptors.request.use(
  (config) => {
    // Pegar token do localStorage ou cookie
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Adicionar User-ID se disponível
    const userId = localStorage.getItem('userId');
    if (userId) {
      config.headers['X-User-ID'] = userId;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor para tratar respostas e erros
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token expirado ou inválido
      localStorage.removeItem('authToken');
      localStorage.removeItem('userId');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Funções específicas para cada serviço
export const authApi = {
  login: (email: string, password: string) =>
    apiClient.post('/api/v1/auth/login', { email, password }),
  
  register: (email: string, password: string, name: string) =>
    apiClient.post('/api/v1/auth/register', { email, password, name }),
  
  logout: () =>
    apiClient.post('/api/v1/auth/logout'),
  
  refreshToken: () =>
    apiClient.post('/api/v1/auth/refresh'),
  
  profile: () =>
    apiClient.get('/api/v1/auth/profile'),
  
  updateProfile: (data: any) =>
    apiClient.put('/api/v1/auth/profile', data),
  
  magicLink: (email: string) =>
    apiClient.post('/api/v1/auth/magic-link', { email }),
  
  verifyMagicLink: (token: string) =>
    apiClient.post('/api/v1/auth/verify-magic-link', { token })
};

export const sitesApi = {
  list: () =>
    apiClient.get('/api/v1/sites'),
  
  get: (siteId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}`),
  
  create: (data: any) =>
    apiClient.post('/api/v1/sites', data),
  
  update: (siteId: string, data: any) =>
    apiClient.put(`/api/v1/sites/${siteId}`, data),
  
  delete: (siteId: string) =>
    apiClient.delete(`/api/v1/sites/${siteId}`),
  
  status: (siteId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}/status`),
  
  analytics: (siteId: string, params?: any) =>
    apiClient.get(`/api/v1/sites/${siteId}/analytics`, { params }),
  
  stats: (siteId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}/stats`)
};

export const deploymentsApi = {
  list: (siteId: string, params?: any) =>
    apiClient.get(`/api/v1/sites/${siteId}/deployments`, { params }),
  
  get: (siteId: string, deployId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}/deployments/${deployId}`),
  
  create: (siteId: string, data: any) =>
    apiClient.post(`/api/v1/sites/${siteId}/deployments`, data),
  
  rollback: (siteId: string, deployId: string) =>
    apiClient.post(`/api/v1/sites/${siteId}/deployments/${deployId}/rollback`),
  
  purgeCache: (siteId: string) =>
    apiClient.post(`/api/v1/sites/${siteId}/cache/purge`),
  
  cacheStatus: (siteId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}/cache/status`)
};

export const domainsApi = {
  list: (siteId: string) =>
    apiClient.get(`/api/v1/sites/${siteId}/domains`),
  
  add: (siteId: string, data: any) =>
    apiClient.post(`/api/v1/sites/${siteId}/domains`, data),
  
  update: (siteId: string, domainId: string, data: any) =>
    apiClient.put(`/api/v1/sites/${siteId}/domains/${domainId}`, data),
  
  delete: (siteId: string, domainId: string) =>
    apiClient.delete(`/api/v1/sites/${siteId}/domains/${domainId}`),
  
  verify: (siteId: string, domainId: string) =>
    apiClient.post(`/api/v1/sites/${siteId}/domains/${domainId}/verify`)
};

export const buildsApi = {
  list: (siteId: string, params?: any) =>
    apiClient.get(`/api/v1/sites/${siteId}/builds`, { params }),
  
  get: (jobId: string) =>
    apiClient.get(`/api/v1/builds/${jobId}`),
  
  create: (data: any) =>
    apiClient.post('/api/v1/builds', data),
  
  cancel: (jobId: string) =>
    apiClient.post(`/api/v1/builds/${jobId}/cancel`),
  
  retry: (jobId: string) =>
    apiClient.post(`/api/v1/builds/${jobId}/retry`),
  
  logs: (jobId: string) =>
    apiClient.get(`/api/v1/builds/${jobId}/logs`)
};

export const templatesApi = {
  list: (params?: any) =>
    apiClient.get('/api/v1/templates', { params }),
  
  get: (templateId: string) =>
    apiClient.get(`/api/v1/templates/${templateId}`),
  
  featured: (params?: any) =>
    apiClient.get('/api/v1/templates/featured', { params }),
  
  search: (params?: any) =>
    apiClient.get('/api/v1/templates/search', { params }),
  
  categories: () =>
    apiClient.get('/api/v1/templates/categories'),
  
  frameworks: () =>
    apiClient.get('/api/v1/templates/frameworks'),
  
  download: (templateId: string) =>
    apiClient.post(`/api/v1/templates/${templateId}/download`)
};

export const billingApi = {
  subscription: () =>
    apiClient.get('/api/v1/billing/subscription'),
  
  createSubscription: (planId: string) =>
    apiClient.post('/api/v1/billing/subscription', { plan_id: planId }),
  
  cancelSubscription: () =>
    apiClient.delete('/api/v1/billing/subscription'),
  
  invoices: (params?: any) =>
    apiClient.get('/api/v1/billing/invoices', { params }),
  
  paymentMethods: () =>
    apiClient.get('/api/v1/billing/payment-methods'),
  
  addPaymentMethod: (data: any) =>
    apiClient.post('/api/v1/billing/payment-methods', data),
  
  deletePaymentMethod: (methodId: string) =>
    apiClient.delete(`/api/v1/billing/payment-methods/${methodId}`),
  
  usage: (params?: any) =>
    apiClient.get('/api/v1/billing/usage', { params })
};

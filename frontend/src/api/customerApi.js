import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8090/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export function setAuthToken(token) {
  api.defaults.headers.common.Authorization = token ? `Bearer ${token}` : '';
}

export const fetchCustomers = () => api.get('/customers');
export const fetchCustomerById = (id) => api.get(`/customers/${id}`);
export const createCustomer = (payload) => api.post('/customers', payload);
export const updateCustomer = (id, payload) => api.put(`/customers/${id}`, payload);
export const deleteCustomer = (id) => api.delete(`/customers/${id}`);

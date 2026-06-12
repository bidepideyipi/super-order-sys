import axios from 'axios';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
});

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

const tauriAPI = {
  getVersion: () => Promise.resolve('1.0.0-web'),
  
  openFile: async () => {
    console.warn('openFile not available in web version');
    return [];
  },
  
  saveFile: async () => {
    console.warn('saveFile not available in web version');
    return null;
  },
  
  openExternal: (url) => {
    window.open(url, '_blank');
  },

  sku: {
    list: () => api.get('/sku/list').then(res => res.data.data),

    listPaginated: (page, pageSize) =>
      api.get('/sku/list-paginated', { params: { page, page_size: pageSize } })
        .then(res => res.data.data),

    get: (id) => api.get(`/sku/${id}`).then(res => res.data.data),

    getImage: (skuCode) =>
      api.get(`/common/image/${skuCode}`).then(res => res.data.data),

    create: (data, imageBase64) =>
      api.post('/sku', { ...data, image_base64: imageBase64 }).then(res => res.data.data),

    update: (id, data, imageBase64) =>
      api.put(`/sku/${id}`, { ...data, image_base64: imageBase64 }).then(res => res.data.data),

    delete: (id) => api.delete(`/sku/${id}`).then(res => res.data.data),

    search: (keyword) =>
      api.get('/sku/search', { params: { keyword } }).then(res => res.data.data),

    searchPaginated: (keyword, page, pageSize) =>
      api.get('/sku/search-paginated', { params: { keyword, page, page_size: pageSize } })
        .then(res => res.data.data),

    searchWithCategory: (keyword, categoryId) =>
      api.get('/sku/search-with-category', { params: { keyword, category_id: categoryId } })
        .then(res => res.data.data)
  },

  category: {
    list: () => api.get('/category/list').then(res => res.data.data)
  },

  customer: {
    list: () => api.get('/customer/list').then(res => res.data.data)
  },

  order: {
    list: () => api.get('/order/list').then(res => res.data.data),

    get: (id) => api.get(`/order/${id}`).then(res => res.data.data),

    create: (data) => api.post('/order', data).then(res => res.data.data),

    update: (id, data) => api.put(`/order/${id}`, data).then(res => res.data.data),

    delete: (id) => api.delete(`/order/${id}`).then(res => res.data.data)
  },

  purchase: {
    getProcessingOrders: () =>
      api.get('/order/processing').then(res => res.data.data),

    getUnsettledOrders: () =>
      api.get('/order/unsettled').then(res => res.data.data),

    getOrderItems: (orderId) =>
      api.get(`/order/${orderId}/items`).then(res => res.data.data),

    searchSkuByCode: (keyword) =>
      api.get('/sku/search', { params: { keyword } }).then(res => res.data.data),

    createOrderItem: (data) =>
      api.post('/order-item', data).then(res => res.data.data),

    updateOrderItem: (id, data) =>
      api.put(`/order-item/${id}`, data).then(res => res.data.data),

    deleteOrderItem: (id) =>
      api.delete(`/order-item/${id}`).then(res => res.data.data)
  },

  financial: {
    list: () => api.get('/financial/list').then(res => res.data.data),

    get: (id) => api.get(`/financial/${id}`).then(res => res.data.data),

    create: (data) => api.post('/financial', data).then(res => res.data.data),

    update: (id, data) => api.put(`/financial/${id}`, data).then(res => res.data.data),

    delete: (id) => api.delete(`/financial/${id}`).then(res => res.data.data),

    getBalance: () => api.get('/financial/balance').then(res => res.data.data)
  }
};

if (typeof window !== 'undefined') {
  window.tauriAPI = tauriAPI;
}

export default tauriAPI;

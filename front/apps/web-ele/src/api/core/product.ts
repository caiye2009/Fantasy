import { requestClient } from '../request';
import { search } from './search';

export interface Product {
  id: number;
  productCode: string;
  productName: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductListResponse {
  total: number;
  list: Product[];
}

export async function createProduct(data: {
  productCode: string;
  productName: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Product>('/products', data);
  return response;
}

export async function getProduct(id: number) {
  const response = await requestClient.get<Product>(`/products/${id}`);
  return response;
}

export async function updateProduct(
  id: number,
  data: {
    productCode?: string;
    productName?: string;
  }
) {
  const response = await requestClient.post<Product>(`/products/${id}`, data);
  return response;
}

export async function deleteProduct(id: number) {
  const response = await requestClient.delete(`/products/${id}`);
  return response;
}

// List queries use search API
export async function getProducts(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'product',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.hits || [],
  };
}

// Product price API (stub for now)
export interface ProductPriceResponse {
  current_price: number;
  historical_high: number;
  historical_low: number;
}

export async function getProductPriceApi(_productId: number): Promise<ProductPriceResponse> {
  // TODO: Implement actual product pricing logic
  return {
    current_price: 0,
    historical_high: 0,
    historical_low: 0,
  };
}

// Backward compatibility aliases
export const getProductListApi = getProducts;
export const getProductDetailApi = getProduct;
export const createProductApi = createProduct;
export const updateProductApi = updateProduct;
export const deleteProductApi = deleteProduct;

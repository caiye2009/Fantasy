import { requestClient } from '../request';

export interface ProductMaterial {
  id: number;
  productCode: string;
  materialCode: string;
  requiredQty: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductMaterialListResponse {
  total: number;
  list: ProductMaterial[];
}

export async function createProductMaterial(data: {
  productCode: string;
  materialCode: string;
  requiredQty: number;
  createdBy: number;
}) {
  const response = await requestClient.post<ProductMaterial>('/product_materials', data);
  return response;
}

export async function getProductMaterials(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<ProductMaterialListResponse>('/product_materials', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getProductMaterial(id: number) {
  const response = await requestClient.get<ProductMaterial>(`/product_materials/${id}`);
  return response;
}

export async function updateProductMaterial(
  id: number,
  data: {
    productCode?: string;
    materialCode?: string;
    requiredQty?: number;
  }
) {
  const response = await requestClient.post<ProductMaterial>(`/product_materials/${id}`, data);
  return response;
}

export async function deleteProductMaterial(id: number) {
  const response = await requestClient.delete(`/product_materials/${id}`);
  return response;
}

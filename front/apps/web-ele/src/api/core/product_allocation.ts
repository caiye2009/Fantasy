import { requestClient } from '../request';

export interface ProductAllocation {
  id: number;
  allocationCode: string;
  orderCode: string;
  productCode: string;
  allocationQty: number;
  actualQty: number;
  allocationStatus: string;
  allocatedBy: number;
  allocatedAt: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductAllocationListResponse {
  total: number;
  list: ProductAllocation[];
}

export async function createProductAllocation(data: {
  allocationCode: string;
  orderCode: string;
  productCode: string;
  allocationQty: number;
  actualQty?: number;
  allocationStatus: string;
  allocatedBy: number;
  allocatedAt?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<ProductAllocation>('/product_allocations', data);
  return response;
}

export async function getProductAllocations(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<ProductAllocationListResponse>('/product_allocations', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getProductAllocation(id: number) {
  const response = await requestClient.get<ProductAllocation>(`/product_allocations/${id}`);
  return response;
}

export async function updateProductAllocation(
  id: number,
  data: {
    allocationCode?: string;
    orderCode?: string;
    productCode?: string;
    allocationQty?: number;
    actualQty?: number;
    allocationStatus?: string;
    allocatedBy?: number;
    allocatedAt?: string;
  }
) {
  const response = await requestClient.post<ProductAllocation>(`/product_allocations/${id}`, data);
  return response;
}

export async function deleteProductAllocation(id: number) {
  const response = await requestClient.delete(`/product_allocations/${id}`);
  return response;
}

// 获取产品库存信息
export interface InventoryResponse {
  batchId: string;
  category: string;
  quantity: number;
  unit: string;
  unitCost: number;
  totalCost: number;
  remark: string;
}

export async function getInventoriesByProductIdApi(productId: number): Promise<{ inventories: InventoryResponse[] }> {
  // TODO: 实现获取产品库存的逻辑
  // 可能需要后端添加按product_id筛选的接口
  console.warn('getInventoriesByProductIdApi not implemented yet, returning empty inventories');
  return { inventories: [] };
}

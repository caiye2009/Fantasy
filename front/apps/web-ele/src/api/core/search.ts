import { requestClient } from '../request';

export interface SearchFilter {
  field: string;
  operator: string;
  value: any;
}

export interface SearchAggregation {
  field: string;
  type: string;
}

export interface SearchRequest {
  index: string;
  query?: string;
  filters?: SearchFilter[];
  aggregations?: SearchAggregation[];
  from?: number;
  size?: number;
  sortField?: string;
  sortOrder?: string;
}

export interface SearchResponse {
  total: number;
  items: any[];
  aggregations?: Record<string, any>;
}

export interface SearchIndicesResponse {
  indices: string[];
}

export interface ListResponse<T> {
  total: number;
  list: T[];
}

export interface AggOption {
  label: string;
  value: any;
}

export interface FetchAggOptionsParams {
  index: string;
  field: string;
  size?: number;
  labelFormatter?: (key: any) => string;
}


export const search = async (data: SearchRequest) => {
  // 转换为后端期望的格式
  const backendRequest: any = {
    index: data.index,
    query: data.query,
    pagination: {
      offset: data.from || 0,
      size: data.size || 20,
    },
  };

  // 转换 aggregations 数组为 aggRequests 对象
  if (data.aggregations && data.aggregations.length > 0) {
    backendRequest.aggRequests = {};
    data.aggregations.forEach((agg) => {
      backendRequest.aggRequests[agg.field] = {
        size: 20,
      };
    });
  }

  // 转换 filters 数组为对象
  if (data.filters && data.filters.length > 0) {
    backendRequest.filters = {};
    data.filters.forEach((filter) => {
      backendRequest.filters[filter.field] = filter.value;
    });
  }

  // 转换 sort
  if (data.sortField) {
    backendRequest.sort = [{
      field: data.sortField,
      order: data.sortOrder || 'asc',
    }];
  }

  const response = await requestClient.post<SearchResponse>('/search', backendRequest);
  return response;
};

export const getSearchIndices = async () => {
  const response =
    await requestClient.get<SearchIndicesResponse>('/search/indices');
  return response;
};

export const fetchAggOptions = async (
  params: FetchAggOptionsParams,
): Promise<AggOption[]> => {
  const { index, field, size = 20, labelFormatter } = params;

  try {
    const response = await search({
      index,
      from: 0,
      size: 0,
      aggregations: [
        {
          field,
          type: 'terms',
        },
      ],
    });

    const buckets =
      response.data.aggregations?.[field]?.buckets ?? [];
    return buckets.map((bucket: any) => ({
      label: labelFormatter
        ? labelFormatter(bucket.key)
        : String(bucket.key),
      value: bucket.key,
    }));
  } catch (error) {
    console.error(`fetchAggOptions failed: ${field}`, error);
    return [];
  }
};

export const a = () => {
  console.log("diu lei lou mou")
}
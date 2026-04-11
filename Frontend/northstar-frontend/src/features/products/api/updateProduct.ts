import api from "@/api/axiosInstance";
import type { Product } from "../types";
import type { ProductFormValues } from "../schemas/productSchema";

export const updateProduct = async (payload: ProductFormValues, id: number): Promise<Product> => {
  const { data } = await api.put<Product>(`/products/${id}`, payload);
  return data;
};

import api from "@/api/axiosInstance";
import type { Product } from "../types";
import type { ProductFormValues } from "../schemas/productSchema";

export const createProduct = async (payload: ProductFormValues): Promise<Product> => {
  const { data } = await api.post<Product>("/products", payload);
  return data;
};
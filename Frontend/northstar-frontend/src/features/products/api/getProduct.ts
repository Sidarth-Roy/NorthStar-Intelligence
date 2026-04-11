import api from "@/api/axiosInstance";
import type { Product } from "../types";

export const getProducts = async (): Promise<Product[]> => {
  try {
    const { data } = await api.get<Product[]>("/products");
    
    // Logic Gate: If data is null/undefined, return empty array to prevent crash
    if (!data) {
      console.warn("API returned empty or null data");
      return [];
    }
    console.log("API Response Data:", data); // Debug log to verify response structure

    return data;
  } catch (error) {
    // Re-throw so TanStack Query can catch it and move to 'error' state
    throw error;
  }
};
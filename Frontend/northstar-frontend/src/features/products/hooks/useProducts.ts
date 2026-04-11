import { useQuery } from "@tanstack/react-query";
import { getProducts } from "../api/getProduct";

export const useProducts = () => {
  return useQuery({
    queryKey: ["products"],
    queryFn: getProducts,
    placeholderData: [], // Provides an empty array while loading
    retry: 1,            // Industry standard: don't spam the backend on failure
  });
};
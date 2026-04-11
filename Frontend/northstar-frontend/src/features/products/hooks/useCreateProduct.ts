import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createProduct } from "../api/createProduct";
import { toast } from "sonner";

export const useCreateProduct = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createProduct,
    onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["products"] });
    toast.success("Product created successfully");
    },
    onError: (error: any) => {
    const message = error.response?.data?.message || error.message;
    toast.error(`Failed to create product: ${message}`);
    }
  });
};
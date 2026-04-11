import api from "@/api/axiosInstance";

export const deleteProduct = async (id: number) => {
  await api.delete(`/products/${id}`);
};
import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { productSchema, type ProductFormValues } from "../schemas/productSchema";
import { useCreateProduct } from "../hooks/useCreateProduct";
import { useUpdateProduct } from "../hooks/useUpdateProduct";
import type { Product } from "../types";

interface ProductFormProps {
  initialData?: Product | null;
  onSuccess: () => void;
}

export const ProductForm = ({ initialData, onSuccess }: ProductFormProps) => {
  const isEdit = !!initialData;
  const createMutation = useCreateProduct();
  const updateMutation = useUpdateProduct(initialData?.id ?? 0);

  const { register, handleSubmit, formState: { errors } } = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema) as any,
    defaultValues: initialData ? {
      productName: initialData.productName,
      unitPrice: initialData.unitPrice,
      quantityPerUnit: initialData.quantityPerUnit,
      categoryID: initialData.categoryID,
      discontinued: initialData.discontinued
    } : {
      productName: "",
      unitPrice: 0,
      quantityPerUnit: "",
      categoryID: 1,
      discontinued: 0
    }
  });

  const onSubmit: SubmitHandler<ProductFormValues> = (data) => {
    const mutation = isEdit ? updateMutation : createMutation;
    mutation.mutate(data, { onSuccess });
  };

  const isPending = createMutation.isPending || updateMutation.isPending;
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 p-4 border rounded-lg bg-white">
      <div>
        <label className="block text-sm font-medium">Product Name</label>
        <input {...register("productName")} className="w-full border p-2 rounded" />
        {errors.productName && <p className="text-red-500 text-xs">{errors.productName.message}</p>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium">Unit Price</label>
          <input 
            type="number" 
            step="0.01" 
            {...register("unitPrice")} 
            className="w-full border p-2 rounded" 
          />
          {errors.unitPrice && <p className="text-red-500 text-xs">{errors.unitPrice.message}</p>}
        </div>
        <div>
          <label className="block text-sm font-medium">Category ID</label>
          <input 
            type="number" 
            {...register("categoryID")} 
            className="w-full border p-2 rounded" 
          />
          {errors.categoryID && <p className="text-red-500 text-xs">{errors.categoryID.message}</p>}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium">Quantity Per Unit</label>
        <input {...register("quantityPerUnit")} className="w-full border p-2 rounded" />
        {errors.quantityPerUnit && <p className="text-red-500 text-xs">{errors.quantityPerUnit.message}</p>}
      </div>
    <div>
    <label className="flex items-center gap-2 cursor-pointer">
        <input 
        type="checkbox" 
        // RHF treats checkboxes as booleans, but your Go API wants 0 or 1
        // We handle this by mapping the boolean to a number on change
        className="w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500"
        onChange={(e) => {
            const value = e.target.checked ? 1 : 0;
            register("discontinued").onChange({ target: { value, name: "discontinued" } });
        }}
        defaultChecked={initialData?.discontinued === 1}
        />
        <span className="text-sm font-medium text-slate-700">Mark as Discontinued</span>
    </label>
    <p className="text-xs text-slate-500 mt-1">
        Discontinued products will be hidden from new orders but kept in history.
    </p>
    </div>
      <button 
        type="submit" 
        disabled={isPending}
        className="w-full bg-blue-600 text-white p-2 rounded hover:bg-blue-700 disabled:bg-slate-400 font-medium"
      >
        {isPending ? "Saving..." : isEdit ? "Update Product" : "Create Product"}
      </button>
    </form>
  );
};
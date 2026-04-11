import { z } from "zod";

export const productSchema = z.object({
  productName: z.string().min(3, "Name must be at least 3 characters"),
  // Pass the message as a simple string instead of an object to clear error 2353
  unitPrice: z.coerce.number().positive("Price must be greater than 0"),
  quantityPerUnit: z.string().min(1, "Quantity description is required"),
  categoryID: z.coerce.number().int().positive("Please select a category"),
  discontinued: z.coerce.number().default(0),
});

export type ProductFormValues = z.infer<typeof productSchema>;

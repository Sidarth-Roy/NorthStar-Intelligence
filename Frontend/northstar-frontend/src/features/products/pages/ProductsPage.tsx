import { useState } from "react";
import { ProductList } from "../components/ProductList";
import { ProductForm } from "../components/ProductForm";
import { Modal } from "../components/ui/Modal";
import { Plus } from "lucide-react";
import type { Product } from "../types";

const ProductsPage = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

  const openCreateModal = () => {
    setSelectedProduct(null);
    setIsModalOpen(true);
  };

  const openEditModal = (product: Product) => {
    setSelectedProduct(product);
    setIsModalOpen(true);
  };

  const handleClose = () => {
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-slate-900">Products</h1>
        <button 
          onClick={openCreateModal}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700"
        >
          <Plus className="w-4 h-4" /> Add Product
        </button>
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={handleClose} 
        title={selectedProduct ? "Edit Product" : "Add New Product"}
      >
        <ProductForm initialData={selectedProduct} onSuccess={handleClose} />
      </Modal>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200">
        <ProductList onEdit={openEditModal} />
      </div>
    </div>
  );
};

export default ProductsPage;
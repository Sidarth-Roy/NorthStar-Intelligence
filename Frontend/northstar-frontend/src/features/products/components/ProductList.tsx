import { useState, useMemo } from "react";
import { useProducts } from "../hooks/useProducts";
import { useDeleteProduct } from "../hooks/useDeleteProduct";
import { Edit2, Trash2, Search, ArrowUpDown, ChevronUp, ChevronDown, ChevronLeft, ChevronRight } from "lucide-react";
import type { Product } from "../types/index"; 
import styles from "./ProductList.module.css";

interface ProductListProps {
  onEdit: (product: Product) => void;
}

type SortConfig = {
  key: keyof Product;
  direction: "asc" | "desc";
} | null;

export const ProductList = ({ onEdit }: ProductListProps) => {
  const { data: products = [], isLoading } = useProducts();
  const { mutate: deleteProduct } = useDeleteProduct();
  
  // State for Search, Sort, and Pagination
  const [searchQuery, setSearchQuery] = useState("");
  const [sortConfig, setSortConfig] = useState<SortConfig>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  // 1. Sorting Logic
  const handleSort = (key: keyof Product) => {
    let direction: "asc" | "desc" = "asc";
    if (sortConfig?.key === key && sortConfig.direction === "asc") {
      direction = "desc";
    }
    setSortConfig({ key, direction });
  };

  // 2. Data Processing Pipeline (Filter -> Sort -> Paginate)
  const processedData = useMemo(() => {
    // Filter
    let result = products.filter((p) =>
      p.productName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.categoryName.toLowerCase().includes(searchQuery.toLowerCase())
    );

    // Sort
    if (sortConfig) {
      result.sort((a, b) => {
        const aValue = a[sortConfig.key];
        const bValue = b[sortConfig.key];
        if (aValue < bValue) return sortConfig.direction === "asc" ? -1 : 1;
        if (aValue > bValue) return sortConfig.direction === "asc" ? 1 : -1;
        return 0;
      });
    }

    const totalItems = result.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    
    // Paginate
    const startIndex = (currentPage - 1) * itemsPerPage;
    const paginatedItems = result.slice(startIndex, startIndex + itemsPerPage);

    return { items: paginatedItems, totalPages, totalItems };
  }, [products, searchQuery, sortConfig, currentPage]);

  const handleDelete = (id: number, name: string) => {
    if (window.confirm(`Are you sure you want to delete ${name}?`)) {
      deleteProduct(id);
    }
  };

  // Helper for Sort Icons
  const SortIcon = ({ column }: { column: keyof Product }) => {
    if (sortConfig?.key !== column) return <ArrowUpDown className="w-3 h-3 ml-2 opacity-30" />;
    return sortConfig.direction === "asc" ? 
      <ChevronUp className="w-3 h-3 ml-2 text-blue-600" /> : 
      <ChevronDown className="w-3 h-3 ml-2 text-blue-600" />;
  };

  if (isLoading) return <div className="p-8 text-center text-slate-500">Syncing inventory...</div>;

  return (
    <div className="space-y-4">
      {/* Search Bar */}
      <div className="relative px-4 pt-4">
        <Search className="absolute left-7 top-7 w-4 h-4 text-slate-400" />
        <input
          type="text"
          placeholder="Search products or categories..."
          className="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 transition-all"
          value={searchQuery}
          onChange={(e) => {
            setSearchQuery(e.target.value);
            setCurrentPage(1); // Reset to first page on search
          }}
        />
      </div>

      <div className={styles.tableWrapper}>
        <table className={styles.inventoryTable}>
          <thead className={styles.thead}>
            <tr>
              <th className={`${styles.th} cursor-pointer hover:bg-slate-50`} onClick={() => handleSort("id")}>
                <div className="flex items-center">ID <SortIcon column="id" /></div>
              </th>
              <th className={`${styles.th} cursor-pointer hover:bg-slate-50`} onClick={() => handleSort("productName")}>
                <div className="flex items-center">Product <SortIcon column="productName" /></div>
              </th>
              <th className={`${styles.th} cursor-pointer hover:bg-slate-50`} onClick={() => handleSort("categoryName")}>
                <div className="flex items-center">Category <SortIcon column="categoryName" /></div>
              </th>
              <th className={`${styles.th} cursor-pointer hover:bg-slate-50 text-right`} onClick={() => handleSort("unitPrice")}>
                <div className="flex items-center justify-end">Price <SortIcon column="unitPrice" /></div>
              </th>
              <th className={`${styles.th} cursor-pointer hover:bg-slate-50`} onClick={() => handleSort("discontinued")}>
                <div className="flex items-center">Status <SortIcon column="discontinued" /></div>
              </th>
              <th className={styles.th}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {processedData.items.map((product) => (
              <tr key={product.id} className={styles.tr}>
                <td className={`${styles.td} font-mono text-xs text-slate-400`}>#{product.id}</td>
                <td className={`${styles.td} font-medium text-slate-900`}>{product.productName}</td>
                <td className={styles.tdMuted}>{product.categoryName}</td>
                <td className={`${styles.tdPrice} text-right`}>${product.unitPrice.toFixed(2)}</td>
                <td className={styles.td}>
                  <span className={`${styles.badge} ${product.discontinued ? styles.badgeDiscontinued : styles.badgeActive}`}>
                    {product.discontinued ? "Discontinued" : "Active"}
                  </span>
                </td>
                <td className={styles.td}>
                  <div className="flex items-center gap-1">
                    <button onClick={() => onEdit(product)} className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-all">
                      <Edit2 className="w-4 h-4" />
                    </button>
                    <button onClick={() => handleDelete(product.id, product.productName)} className="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-all">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        
        {processedData.items.length === 0 && (
          <div className="p-12 text-center text-slate-500 text-sm italic">
            No items found matching "{searchQuery}"
          </div>
        )}
      </div>

      {/* Minimalist Pagination Controls */}
      <div className="flex items-center justify-between px-4 pb-4">
        <div className="text-xs text-slate-500">
          Showing {processedData.items.length} of {processedData.totalItems} products
        </div>
        <div className="flex items-center gap-2">
          <button 
            disabled={currentPage === 1}
            onClick={() => setCurrentPage(prev => prev - 1)}
            className="p-1 rounded border border-slate-200 disabled:opacity-30 hover:bg-slate-50"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <span className="text-xs font-medium px-2">Page {currentPage} of {processedData.totalPages || 1}</span>
          <button 
            disabled={currentPage === processedData.totalPages || processedData.totalPages === 0}
            onClick={() => setCurrentPage(prev => prev + 1)}
            className="p-1 rounded border border-slate-200 disabled:opacity-30 hover:bg-slate-50"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
};
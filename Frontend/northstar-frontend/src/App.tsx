import { lazy, Suspense } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { Toaster } from 'sonner';

// Industry Best Practice: Lazy Loading features
const ProductsPage = lazy(() => import("@/features/products/pages/ProductsPage"));

// Simple Loader for Suspense
const PageLoader = () => (
  <div className="h-full w-full flex items-center justify-center p-12">
    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
  </div>
);

function App() {
  return (
    <BrowserRouter>
      <Toaster position="top-right" richColors />
      <Routes>
        <Route path="/" element={<AppLayout />}>
          <Route index element={<div>Dashboard Home</div>} />
          
          <Route 
            path="products" 
            element={
              <Suspense fallback={<PageLoader />}>
                <ProductsPage />
              </Suspense>
            } 
          />

          {/* Future modules... */}
          <Route path="categories" element={<div>Categories Module</div>} />
          <Route path="orders" element={<div>Orders Module</div>} />
          <Route path="customers" element={<div>Customers Module</div>} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
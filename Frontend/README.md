```markdown
# NorthStar Intelligence - React Frontend

An enterprise-grade ERP Dashboard built with **React**, **TypeScript**, and **Vite**. This application serves as the primary interface for the NorthStar Intelligence ecosystem, communicating with the Go/Gin backend to manage products, orders, customers, and supply chain logistics.

---

## 🛠️ Tech Stack

* **Framework:** React 18 (Vite build tool)
* **Language:** TypeScript (Strict Mode)
* **State Management:** * **Server State:** TanStack Query (React Query) v5
    * **Client State:** Zustand (Lightweight, modular stores)
* **Styling:** Tailwind CSS + Shadcn/UI (Radix + Nova)
* **Forms:** React Hook Form + Zod (Schema-based validation)
* **Routing:** TanStack Router (Type-safe routing)
* **API Client:** Axios (Custom instance with Request-ID interceptors)

---

## 🏗️ Project Architecture (Feature-First)

```text
src/
├── api/              # Axios configuration & global interceptors
├── components/       # Shared UI components (Buttons, Modals, Tables)
├── hooks/            # Global reusable hooks
├── lib/              # Utility configurations (utils, shadcn helpers)
├── store/            # Global Zustand stores (Auth, Theme)
└── features/         # Domain-driven modules (CRUD)
```

---

## 🗺️ Roadmap & Implementation Stages

### Stage 1: Foundation & Core Setup
- [x] **Infrastructure:** Initialize Vite + TS and configure path aliases (`@/`).
- [x] **Theming:** Configure Tailwind v4 and Shadcn/UI primitives.
- [ ] **API Layer:** Set up Axios instance with interceptors to log the backend's `X-Request-ID`.
- [ ] **Types:** Define global Base Response interfaces to match Go Gin's JSON output.

### Stage 2: Layout & Navigation
- [ ] **Dashboard Shell:** Sidebar navigation, breadcrumbs, and user profile header.
- [ ] **Routing:** Initialize TanStack Router with type-safe route definitions.
- [ ] **Error Handling:** Create a global Error Boundary and a 404 page.

### Stage 3: The "Generic CRUD" Pattern
- [ ] **Data Table:** Create a reusable `<DataTable />` with sorting, filtering, and pagination.
- [ ] **Forms:** Build a dynamic form builder using `react-hook-form` and `zod`.
- [ ] **Feedback:** Implement Toast notifications for API success/failure.

---

## 💎 Industry Best Practices Applied

1.  **Type Safety:** End-to-end typing from API responses to form inputs.
2.  **Server State:** No manual `useEffect` for data fetching. All async state is managed by React Query.
3.  **Component Composition:** Using the "Compound Component" pattern.
4.  **Request Tracing:** Propagation of `X-Request-ID` across the stack.
```

---

## 🚀 Your Next 3 Technical Steps

Following best practices, you should not build UI yet. You must build the **Infrastructure Layer** first.

### Step 1: Define the Global API Response Types
In Go, you likely have a standard response wrapper. You need to mirror that in TypeScript so your frontend always knows the "Shape" of the data coming from Gin.

**File:** `src/api/types.ts`
```typescript
export interface BaseResponse<T> {
  data: T;
  message: string;
  status: string;
  timestamp: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total_records: number;
  total_pages: number;
}
```

### Step 2: Implement the Axios Interceptor (Middleware)
As a backend dev, you know logging is key. This interceptor ensures that every frontend request is linked to a backend log entry via the `X-Request-ID`.

**File:** `src/api/axiosInstance.ts`
```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Standardized error handling - log the Request ID for backend debugging
    const requestId = error.response?.headers['x-request-id'];
    console.error(`[API Error] ${error.config.url} | ID: ${requestId}`, error.response?.data);
    return Promise.reject(error);
  }
);

export default api;
```

### Step 3: Setup the TanStack Router
This is the "URL Dispatcher" for your frontend. It makes your links type-safe, meaning if you rename a route, the compiler will catch broken links.

**Run this command:**
```bash
npx shadcn@latest add separator breadcrumb
```

Then, you'll want to start building the **Layout Shell** (The Sidebar). 

**Are your Backend APIs currently allowing CORS for `localhost:5173`?** (This is the #1 reason backend devs get stuck on their first React project).

## 1. Scaffold the Project
Start by creating the Vite project and installing the core dependencies from your stack.

```bash
# Initialize Vite with React and TypeScript
npm create vite@latest northstar-frontend -- --template react-ts
cd northstar-frontend

# Install core libraries
npm install axios @tanstack/react-query @tanstack/router zustand lucide-react
npm install react-hook-form zod @hookform/resolvers
```

---

## 2. Setup Tailwind & Shadcn/UI
Shadcn is perfect for backend engineers because it gives you raw code rather than an opaque library.

```bash
# Install Tailwind and its dependencies
npm install -D tailwindcss/vite postcss autoprefixer @types/node
npx tailwindcss init -p

# Initialize Shadcn (Follow the prompts: use 'Slate' and 'Default')
npx shadcn@latest init
```

---

## 3. The "Backend-Friendly" Folder Structure
Run this command in your root directory to generate the architecture defined in your README:

```bash
mkdir -p src/{api,components,hooks,lib,store,features/{products,orders,customers}/{api,components,hooks,types}}
```

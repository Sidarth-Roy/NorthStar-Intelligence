
function App() {
  return (
    <div className="min-h-screen bg-slate-50 p-8">
      <h1 className="text-2xl font-bold text-slate-900">
        NorthStar Intelligence Dashboard
      </h1>
      <p className="text-slate-600">Backend API connected: {import.meta.env.VITE_API_BASE_URL}</p>
    </div>
  )
}

export default App
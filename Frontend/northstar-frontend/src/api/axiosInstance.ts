import axios, { AxiosError } from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Response Interceptor
api.interceptors.response.use(
  (response) => {
    // Optional: Log the request ID on success for traceability
    const requestId = response.headers["x-request-id"];
    if (requestId) {
      console.debug(`[SUCCESS] RequestID: ${requestId}`);
    }
    return response;
  },
  (error: AxiosError<any>) => {
    const requestId = error.response?.headers["x-request-id"];
    const errorMessage = error.response?.data?.message || "An unexpected error occurred";

    // Industry Practice: Log errors with their Trace ID
    console.error(
      `[API ERROR] | URL: ${error.config?.url} | ID: ${requestId || "N/A"} | Message: ${errorMessage}`
    );

    // You can also trigger a Global Toast notification here later
    return Promise.reject({
      ...error,
      request_id: requestId,
    });
  }
);

export default api;
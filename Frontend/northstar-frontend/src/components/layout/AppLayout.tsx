import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import styles from "./Sidebar.module.css";

export const AppLayout = () => {
  return (
    <div className={styles.mainLayout}>
      <Sidebar />
      <div className={styles.contentWrapper}>
        <header className={styles.header}>
          {/* This header can later hold global search, user profile, etc. */}
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <span>Workspace</span>
            <span>/</span>
            <span className="font-medium text-slate-900">Lab 37</span>
          </div>
        </header>
        <main className={styles.mainArea}>
          <Outlet /> {/* Child routes will render here */}
        </main>
      </div>
    </div>
  );
};
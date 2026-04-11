import { NavLink } from "react-router-dom";
import { 
  LayoutDashboard, 
  Package, 
  Tags, 
  Users, 
  ShoppingCart, 
  Settings 
} from "lucide-react";
import styles from "./Sidebar.module.css";

const navItems = [
  { name: "Dashboard", path: "/", icon: LayoutDashboard },
  { name: "Products", path: "/products", icon: Package },
  { name: "Categories", path: "/categories", icon: Tags },
  { name: "Orders", path: "/orders", icon: ShoppingCart },
  { name: "Customers", path: "/customers", icon: Users },
];

export const Sidebar = () => {
  return (
    <aside className={styles.aside}>
      <div className={styles.brand}>
        <h1 className={styles.brandTitle}>NorthStar</h1>
        <p className={styles.brandSubtitle}>Intelligence ERP</p>
      </div>

      <nav className={styles.nav}>
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => 
                isActive ? styles.navItemActive : styles.navItem
              }
            >
              <Icon className={styles.icon} />
              {item.name}
            </NavLink>
          );
        })}
      </nav>

      <div className="p-4 mt-auto border-t border-slate-800">
        <NavLink to="/settings" className={styles.navItem}>
          <Settings className={styles.icon} />
          Settings
        </NavLink>
      </div>
    </aside>
  );
};
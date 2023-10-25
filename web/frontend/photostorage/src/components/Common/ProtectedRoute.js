import React from 'react';
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedRoute = (props) => {
  const navigate = useNavigate();
  const { user, target, redirect, children } = props
  const adminRoleID = 0

  const isAdmin = (user) => {
    return user.role.id === adminRoleID
  }

  React.useEffect(() => {
    if (user === null) {
      navigate(redirect);
    }
    if (target === "admin" && !isAdmin(user)) { 
      navigate(redirect);
    }
  });

  return children ? children : <Outlet />;
};
export default ProtectedRoute;
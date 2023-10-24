import React from 'react';
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedRoute = (props) => {
  const navigate = useNavigate();
  const adminRoleID = 0

  const isAdmin = (user) => {
    return user.role.id === adminRoleID
  }

  React.useEffect(() => {
    if (props.user === null) {
      navigate(props.redirect);
    }
    if (props.target === "admin" && !isAdmin(props.user)) { 
      navigate(props.redirect);
    }
  });

  return props.children ? props.children : <Outlet />;
};
export default ProtectedRoute;
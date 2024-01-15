import React from 'react';
import PropTypes from "prop-types";
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedRoute = ({ user, target, redirect, children }) => {
  const navigate = useNavigate();
  const adminRoleID = 1

  const isAdmin = (user) => {
    return user.role.id === adminRoleID
  }

  React.useEffect(() => {
    if (user === null) {
      navigate(redirect);
      return
    }
    if (target === "admin" && !isAdmin(user)) { 
      navigate(redirect);
    }
  });

  return children ? children : <Outlet />;
};

ProtectedRoute.propTypes = {
  user: PropTypes.object.isRequired,
  redirect: PropTypes.string.isRequired,
  target: PropTypes.string
};

export default ProtectedRoute;
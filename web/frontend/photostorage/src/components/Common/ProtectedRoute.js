import React from 'react';
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedRoute = (props) => {
  const navigate = useNavigate();

  React.useEffect(() => {
    if (props.user === null) {
      navigate(props.redirect);
    }
  });

  return props.children ? props.children : <Outlet />;
};
export default ProtectedRoute;
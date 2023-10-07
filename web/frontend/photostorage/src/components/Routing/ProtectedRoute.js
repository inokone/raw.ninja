import React from 'react';
import { Outlet, useNavigate } from "react-router-dom";

const ProtectedRoute = ({
  user,
  setUser, 
  redirectPath = '/login',
  children,
}) => {
  const { REACT_APP_API_PREFIX } = process.env;
  const navigate = useNavigate();

  const updateLoggedinUser = async (setUser, redirectPath) => {

  fetch(REACT_APP_API_PREFIX + '/api/v1/auth/profile', {
    method: "GET",
    mode: "cors",
    credentials: "include",
    headers: {
        "Content-Type": "application/json",
    }
  })
  .then(response => {
    if (!response.ok) {
      if (response.status !== 200) {
        console.log(response.status + ": " + response.statusText);
        navigate(redirectPath)
      } else {
        response.json().then(content => console.log(content))
        navigate(redirectPath)
      }
    } else {
      response.json().then(content => {
        console.log("Collected user, setting state: " + content)
        setUser(content)
      })
    }
  })
  .catch(error => console.log(error));
}

  if (!user) {
    console.log("Missing user, trying to collect from profile...")
    updateLoggedinUser(setUser, redirectPath)
  }
  return children ? children : <Outlet />;
};
export default ProtectedRoute;
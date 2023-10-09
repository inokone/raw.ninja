import React, {useEffect} from 'react';
import { useNavigate } from "react-router-dom"
import { CircularProgress } from "@mui/material";


const { REACT_APP_API_PREFIX } = process.env;

const Logout = ({setUser}) => {
  const navigate = useNavigate();

  const logout = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/auth/logout', {
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
                } else {
                    response.json().then(content => console.log(content))
                }
            } else {
                setUser(null);
                navigate("/", { replace: true });
            }
        })
        .catch(error => console.log(error));
    }

  useEffect(logout);


  return <CircularProgress />;
};

export default Logout;
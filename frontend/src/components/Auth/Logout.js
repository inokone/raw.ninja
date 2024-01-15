import React, { useEffect } from 'react';
import { useNavigate } from "react-router-dom"
import ProgressDisplay from '../Common/ProgressDisplay';
import PropTypes from "prop-types";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Logout = ({ setUser }) => {
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
                    throw new Error(response.status + ": " + response.statusText);
                } else {
                    setUser(null);
                    navigate("/", { replace: true });
                }
            })
            .catch(error => console.log(error));
    }

    useEffect(logout);


    return <ProgressDisplay />;
};

Logout.propTypes = {
    setUser: PropTypes.func.isRequired
};

export default Logout;
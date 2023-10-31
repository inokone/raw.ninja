import React, { useState } from "react";
import { Alert, Box, Container } from "@mui/material";
import { useSearchParams, useNavigate } from 'react-router-dom'

const { REACT_APP_API_PREFIX } = process.env;

const EmailConfirmation = () => {
    const [error, setError] = useState(null)
    const [success, setSuccess] = useState(false)
    const [queryParameters] = useSearchParams()
    const navigate = useNavigate()

    const confirmEmail = () => {
        setError(null)
        fetch(REACT_APP_API_PREFIX + '/api/v1/account/confirm?token=' + queryParameters.get("token"), {
            method: "GET",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            }
        }).then(response => {
            if (!response.ok) {
                response.json().then(content => {
                    setError(content.message)
                })
            } else {
                setError(null)
                setSuccess(true)
                navigate("/")
            }
        }).catch(error => {
            setError(error.message)
            })
    }

    React.useEffect(() => {
        confirmEmail()
    })

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box sx={{ width: 500, m: 4 }}>
                    {success ? <Alert sx={{ mb: 4 }} onClose={setSuccess(null)} severity="success">Email confirmed!</Alert> : null}
                    {error ? <Alert sx={{ mb: 4 }} onClose={setError(null)} severity="error">{error}</Alert> : null}
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default EmailConfirmation;
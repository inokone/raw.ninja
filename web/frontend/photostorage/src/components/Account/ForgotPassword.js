import React, { useState } from "react";
import { TextField, Button, Alert, Box, Container } from "@mui/material";

const { REACT_APP_API_PREFIX } = process.env;

const ForgotPassword = () => {
    const [email, setEmail] = useState("")
    const [error, setError] = useState(null)
    const [emailError, setEmailError] = useState(null)
    const [success, setSuccess] = useState(false)

    const handleClick = (event) => {
        setEmailError(email === '')
        if (emailError) {
            return
        }

        setError(null)
        fetch(REACT_APP_API_PREFIX + '/api/v1/account/password/request', {
            method: "PUT",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "email": email
            })
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    })
                } else {
                    setError(null)
                    setSuccess(true)
                }
            })
            .catch(error => {
                setError(error.message)
            });
    }

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box sx={{ width: 500, m: 4 }}>
                    <TextField
                        label="Reset Password"
                        onChange={e => {
                            setEmail(e.target.value)
                            setError(null)
                        }}
                        required
                        variant="outlined"
                        color="primary"
                        type="email"
                        sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                        fullWidth
                        value={email}
                        error={emailError}
                    />
                    {success ? <Alert sx={{ mb: 4 }} onClose={setSuccess(null)} severity="success">Password reset email sent!</Alert> : null}
                    {error ? <Alert sx={{ mb: 4 }} onClose={setError(null)} severity="error">{error}</Alert> : null}
                    <Button sx={{ mb: 4 }} variant="contained" color="primary" onClick={handleClick}>Request Password Reset</Button>
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default ForgotPassword;
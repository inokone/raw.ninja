import React, { useState } from "react";
import { TextField, Button, Alert, Box, Container, Typography } from "@mui/material";

const { REACT_APP_API_PREFIX } = process.env;

const RecoverPassword = () => {
    const [email, setEmail] = useState("")
    const [error, setError] = useState(null)
    const [emailError, setEmailError] = useState(null)
    const [success, setSuccess] = useState(false)

    const handleClick = () => {
        setEmailError(email === '')
        if (emailError) {
            return
        }

        setError(null)
        fetch(REACT_APP_API_PREFIX + '/api/v1/account/recover', {
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
                if (response.status !== 202) {
                    response.json().then(content => {
                        setError(content.message)
                    })
                } else {
                    console.log("SUCCESS!")
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
                    <Typography sx={{mb: 2}}>Please enter your email address so we can send you an email to reset your password.</Typography>
                    <TextField
                        label="Email"
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
                    {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Password reset email sent, check your inbox!</Alert>}
                    {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                    <Button sx={{ mb: 4 }} variant="contained" color="primary" onClick={handleClick}>Request Password Reset</Button>
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default RecoverPassword;
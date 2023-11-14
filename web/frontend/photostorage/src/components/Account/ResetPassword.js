import React, { useState } from "react";
import { TextField, Button, Alert, Box, Container } from "@mui/material";
import { useSearchParams, useNavigate } from 'react-router-dom'

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const ResetPassword = () => {
    const [newPassword, setNewPassword] = useState("")
    const [confirmation, setConfirmation] = useState("")
    const [error, setError] = useState(null)
    const [newPasswordError, setNewPasswordError] = useState(false)
    const [confirmationError, setConfirmationError] = useState(false)
    const [success, setSuccess] = useState(false)
    const [queryParameters] = useSearchParams()
    const navigate = useNavigate()


    const handleClick = (event) => {
        setNewPasswordError(newPassword === '')

        event.preventDefault()
        setError(null)

        if (newPasswordError) {
            return
        }

        fetch(REACT_APP_API_PREFIX + '/api/v1/account/password/reset', {
            method: "PUT",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "token": queryParameters.get("token"),
                "password": newPassword
            })
        }).then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    setSuccess(true)
                    navigate('/')
                }
            })
            .catch(error => {
                setError(error.message)
            })
    }

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box sx={{ width: 500, m: 4 }}>
                    <TextField
                        label="New Password"
                        onChange={e => {
                            let pass = e.target.value
                            setNewPassword(pass)
                            setNewPasswordError(pass.length > 0 && pass.length < 8)
                            setConfirmationError(pass.length > 0 && confirmation !== pass)
                            setError(null)
                        }}
                        required
                        variant="outlined"
                        color="primary"
                        type="password"
                        sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                        fullWidth
                        value={newPassword}
                        error={newPasswordError}
                        helperText={newPasswordError && "New password has to be at least 8 characters long."}
                    />
                    <TextField
                        label="Confirm New Password"
                        onChange={e => {
                            let conf = e.target.value
                            setConfirmation(conf)
                            setConfirmationError(newPassword.length > 0 && conf !== newPassword)
                            setError(null)
                        }}
                        required
                        variant="outlined"
                        color="primary"
                        type="password"
                        value={confirmation}
                        error={confirmationError}
                        fullWidth
                        sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                        helperText={confirmationError && "New password and confirmation must match."}
                    />
                    {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Password changed!</Alert>}
                    {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                    <Button sx={{ mb: 4 }} variant="contained" color="primary" onClick={handleClick}>Change Password</Button>
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default ResetPassword;
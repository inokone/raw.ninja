import React, { useState } from "react";
import PropTypes from "prop-types";
import { TextField, Button, Alert, Box, Container } from "@mui/material";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UserDetails = ({ user }) => {
    const [firstName, setFirstName] = useState(user.first_name)
    const [lastName, setLastName] = useState(user.last_name)
    const [error, setError] = useState(null)
    const [firstNameError, setFirstNameError] = useState(false)
    const [lastNameError, setLastNameError] = useState(false)
    const [success, setSuccess] = useState(false)

    const handleClick = (event) => {
        event.preventDefault()
        setError(null)

        if (lastNameError || firstNameError) {
            return
        }
        user.last_name = lastName
        user.first_name = firstName
        fetch(REACT_APP_API_PREFIX + '/api/v1/users/' + user.id, {
            method: "PUT",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(user)
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    setSuccess(true)
                }
            })
            .catch(error => {
                setError(error.message)
            });
    }

    return (
        <React.Fragment>
            {user &&
                <Container maxWidth="sm">
                    <Box style={{ flex: 1 }} sx={{ m: 4 }}>
                        <TextField
                            label="Firstname"
                            onChange={e => {
                                let name = e.target.value
                                setFirstName(name)
                                setFirstNameError(name.length === 0)
                                setError(null)
                            }}
                            variant="outlined"
                            color="primary"
                            type="text"
                            value={firstName}
                            error={firstNameError}
                            fullWidth
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                            helperText={firstNameError && "Firstname can not be empty"}
                        />
                        <TextField
                            label="Lastname"
                            onChange={e => {
                                let name = e.target.value
                                setLastName(name)
                                setLastNameError(name.length === 0)
                                setError(null)
                            }}
                            variant="outlined"
                            color="primary"
                            type="text"
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                            fullWidth
                            value={lastName}
                            error={lastNameError}
                            helperText={lastNameError && "Lastname can not be empty"}
                        />
                        {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">User details changed!</Alert>}
                        {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 4 }} variant="contained" color="primary" onClick={handleClick}>Save</Button>
                    </Box>
                </Container>
            }
        </React.Fragment>
    );
}


UserDetails.propTypes = {
    user: PropTypes.object
};

export default UserDetails;
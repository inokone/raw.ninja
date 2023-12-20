import React, { useState } from 'react';
import { TextField, Button, Alert, Typography, Container, Box } from '@mui/material';
import { Link, useNavigate, useLocation } from "react-router-dom";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";


const CreateAlbum = () => {
    const navigate = useNavigate()
    const [name, setName] = useState('')
    const [nameError, setNameError] = useState(false)
    const [tags, setTags] = useState([])
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)
    const [loading, setLoading] = useState(false)
    const { state } = useLocation();

    const handleSubmit = (event) => {
        event.preventDefault();
        if (nameError) {
            return
        }
        setError(null)
        fetch(REACT_APP_API_PREFIX + '/api/v1/albums/', {
            method: "POST",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "name": name,
 //               "tags": tags,
                "photos": state.photos
            })
        }).then(response => {
            response.json().then(content => {
                if (!response.ok) {
                    setError(content.message)
                } else {
                    setSuccess(true)
                    navigate("/albums/" + content.id)
                }});
            })
            .catch(() => setError("Network communication error. Maybe backend is down?"))
            .finally(() => {
                setLoading(false)
            });
    }

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box style={{ flex: 1 }} sx={{ m: 4 }}>
                    <Typography pb={3} variant='h4'>New Album</Typography>
                    <form onSubmit={handleSubmit} action={<Link to="/login" />}>
                        <TextField
                            type="text"
                            name="album"
                            autoComplete="album"
                            variant='outlined'
                            color='primary'
                            label="Name"
                            disabled={loading}
                            value={name}
                            onChange={e => {
                                setNameError(name === '')
                                setName(e.target.value)
                                setError(null)
                            }}
                            fullWidth
                            required
                            sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        <TextField
                            type="text"
                            name="tags"
                            variant='outlined'
                            color='primary'
                            label="Tags"
                            disabled={loading}
                            value={tags}
                            onChange={e => {
                                setTags(e.target.value)
                                setError(null)
                            }}
                            fullWidth
                            sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                        />                         
                        {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Created successfully! Loading...</Alert>}
                        {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 4 }} variant="contained" color="primary" type="submit" disabled={loading}>Create</Button>
                    </form>
                </Box>
            </Container>
        </React.Fragment>
    )
}

export default CreateAlbum;
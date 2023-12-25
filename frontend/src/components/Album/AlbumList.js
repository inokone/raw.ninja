import * as React from 'react';
import { useNavigate } from "react-router-dom";
import { Box, Fab, Typography, Alert, Grid } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import AlbumCard from './AlbumCard';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumList = ({ user }) => {
    const navigate = useNavigate()
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [albums, setAlbums] = React.useState(null)

    const populate = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/albums/', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        }).then(response => {
            if (!response.ok) {
                response.json().then(content => {
                    setError(content.message)
                    setLoading(false)
                });
            } else {
                response.json().then(content => {
                    if (content === null) {
                        content = []
                    }
                    setAlbums(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }

    const onAlbumClick = (id) => {
        navigate("/albums/" + id)
    }

    const onFabClick = (id) => {
        navigate("/albums/create")
    }

    React.useEffect(() => {
        if (!loading && !albums && !error) {
            populate()
        }
    }, [albums, loading, error])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {albums !== null &&
            <>
                <Typography variant='h4' sx={{ marginBottom: 4, marginTop: 2 }} >Albums</Typography>
                <Grid container>
                        {albums.map((album) => {
                            return (<Grid item key={album.id} xs={6} md={4} lg={2} xl={2}><AlbumCard album={album} onClick={onAlbumClick} /></Grid>)
                        })}
                </Grid>
                <Box sx={{
                    '& > :not(style)': { m: 1 },
                    position: "fixed",
                    bottom: 16,
                    right: 16
                    }}>
                    <Fab onClick={onFabClick} color="primary" aria-label="add">
                        <AddIcon />
                    </Fab>
                </Box>
            </>}
        </>
    )
}

export default AlbumList;
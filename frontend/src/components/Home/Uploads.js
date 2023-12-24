import * as React from 'react';
import { useNavigate } from "react-router-dom";
import { Box, Alert, Grid } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';
import AlbumCard from '../Album/AlbumCard';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Uploads = ({ user, onDataLoaded }) => {
    const navigate = useNavigate()
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [uploads, setUploads] = React.useState(null)

    const populate = () => {
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/uploads/', {
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
                    setUploads(content)
                    if (onDataLoaded) {
                        onDataLoaded(content)
                    }
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }

    const onAlbumClick = (id) => {
        navigate("/uploads/" + id)
    }

    React.useEffect(() => {
        if (!loading && !uploads && !error) {
            populate()
        }
    }, [uploads, loading, error])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {uploads !== null &&
                <>
                    <Grid container>
                    {uploads.slice(0, 10).map((album) => {
                            return (<Grid item key={album.id} xs={6} md={4} lg={2} xl={2}><AlbumCard album={album} onClick={onAlbumClick} /></Grid>)
                        })}
                    </Grid>
                    <Box sx={{
                        '& > :not(style)': { m: 1 },
                        position: "fixed",
                        bottom: 16,
                        right: 16
                    }}>
                    </Box>
                </>}
        </>
    )
}

export default Uploads;
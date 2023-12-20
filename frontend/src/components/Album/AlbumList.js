import * as React from 'react';
import { Link } from "react-router-dom";
import { List, ListItem, Typography, Alert } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumList = ({ user }) => {
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [albums, setAlbums] = React.useState([])

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
                    setAlbums(content)
                    setLoading(false)
                })
            }
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }

    React.useEffect(() => {
        if (!loading && albums.length === 0 && !error) {
            populate()
        }
    }, [albums, loading, error])

    return (
        <>
            {error && <Alert sx={{ mb: 1 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {albums.length !== 0 &&
            <>
                <Typography variant='h4'>Albums</Typography>
                <List>
                    {albums.map((album) => {
                        return (<ListItem key={album.id}><Link to={'/albums/' + album.id}><Typography>{album.name}</Typography></Link></ListItem>)
                    })}
                </List>
            </>}
        </>
    )
}

export default AlbumList;
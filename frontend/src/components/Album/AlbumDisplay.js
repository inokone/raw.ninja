import * as React from 'react';
import { useLocation, useNavigate } from "react-router-dom";
import { Stack, Chip, Typography, Alert } from "@mui/material";
import EditIcon from '@mui/icons-material/Edit';
import PhotoGrid from '../Photos/PhotoGrid';
import EditAlbumDialog from './EditAlbumDialog';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumDisplay = ({ user }) => {
    const navigate = useNavigate()
    const location = useLocation()
    const path = location.pathname
    const [title, setTitle] = React.useState(null)
    const [data, setData] = React.useState(null)
    const [isEditAlbumDialogOpen, setIsEditAlbumDialogOpen] = React.useState(false);
    const [error, setError] = React.useState()
    const [success, setSuccess] = React.useState(false)
    const [isHovering, setIsHovering] = React.useState(false)

    const populate = () => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
    }

    const update = (name, tags) => {
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "PATCH",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({
                "name": name,
                // "tags": tags,
            })
        }).then(response => {
            response.json().then(content => {
                if (!response.ok) {
                    setError(content.message)
                } else {
                    navigate(0);
                }
            });
        }).catch(() => setError("Network communication error. Maybe backend is down?"))
    }

    const handleDataLoaded = (data) => {
        setTitle(data.name)
        setData(data)
    }

    const handleEditAlbumDialogOpen = React.useCallback(() => {
        setIsEditAlbumDialogOpen(true);
    }, [setIsEditAlbumDialogOpen]);

    const handleEditAlbumDialogClose = React.useCallback(() => {
        setIsEditAlbumDialogOpen(false);
    }, [setIsEditAlbumDialogOpen]);

    const handleEditAlbumDialogSave = React.useCallback((name, tags) => {
        setIsEditAlbumDialogOpen(false);
        update(name, tags)
    }, [setIsEditAlbumDialogOpen]);

    const handleMouseOver = () => {
        setIsHovering(true);
    };

    const handleMouseOut = () => {
        setIsHovering(false);
    };

    return (
        <>
            {title && 
                <Stack sx={{ marginBottom: 2, marginTop: 2 }} justifyContent={'center'} alignItems={'baseline'} direction="row" spacing={1} onClick={handleEditAlbumDialogOpen}
                    onMouseOver={handleMouseOver} onMouseOut={handleMouseOut}>
                    <Typography variant='h4'>{title}</Typography>
                    {isHovering && <EditIcon/>}
                </Stack>
            }
            {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Updated successfully!</Alert>}
            {data && data.tags && <Stack sx={{ marginBottom: 4 }} justifyContent={'center'} direction="row" spacing={1} onClick={handleEditAlbumDialogOpen}>
                {data.tags.map(tag => {
                    return (<Chip key={tag} label={tag} />)
                })} 
            </Stack>}
            {data && <EditAlbumDialog
                open={isEditAlbumDialogOpen}
                onSave={handleEditAlbumDialogSave}
                onCancel={handleEditAlbumDialogClose}
                input = {{name: data.name, tags: data.tags}}
            />}
            <PhotoGrid populator={populate} data={[]} onDataLoaded={handleDataLoaded} />
        </>
    )
}

export default AlbumDisplay;
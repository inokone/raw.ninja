import * as React from 'react';
import { useTheme } from '@mui/material/styles';
import { useLocation, useNavigate } from "react-router-dom";
import { Stack, Chip, Typography, Alert, IconButton, Tooltip, Box, SpeedDial, SpeedDialAction } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import StarIcon from '@mui/icons-material/Star';
import SpeedDialIcon from '@mui/material/SpeedDialIcon';
import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';
import PhotoGrid from '../Photos/PhotoGrid';
import EditAlbumDialog from './EditAlbumDialog';
import DeleteDialog from '../Common/DeleteDialog';
import PropTypes from "prop-types";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumDisplay = ({ user }) => {
    const theme = useTheme();
    const navigate = useNavigate()
    const location = useLocation()
    const path = location.pathname
    const [title, setTitle] = React.useState(null)
    const [data, setData] = React.useState(null)
    const [isEditAlbumDialogOpen, setIsEditAlbumDialogOpen] = React.useState(false);
    const [error, setError] = React.useState()
    const [success, setSuccess] = React.useState(false)
    const [isHovering, setIsHovering] = React.useState(false)
    const [isDeleteDialogOpen, setDeleteDialogOpen] = React.useState(false)
    const [isDeleteAlbumDialogOpen, setDeleteAlbumDialogOpen] = React.useState(false)
    const [deleteItems, setDeleteItems] = React.useState(null)


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

    const update = React.useCallback((name, tags, photos) => {
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "PATCH",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({
                "name": name,
                "tags": tags,
                "photos": photos,
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
    }, [navigate, path])

    const handleDataLoaded = (data) => {
        setTitle(data.name)
        if (data.tags === null) {
            data.tags = []
            setData(data)
        }
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
    }, [setIsEditAlbumDialogOpen, update]);

    const handleMouseOver = () => {
        setIsHovering(true);
    };

    const handleMouseOut = () => {
        setIsHovering(false);
    };

    const handleDeleteAlbumClick = React.useCallback((items) => {
        setDeleteAlbumDialogOpen(true);
    }, [setDeleteAlbumDialogOpen]);

    const handleDeleteAlbumDialogClose = React.useCallback(() => {
        setDeleteAlbumDialogOpen(false);
    }, [setDeleteAlbumDialogOpen]);

    const deleteAlbum = React.useCallback(() => {
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "DELETE",
            mode: "cors",
            credentials: "include"
        }).then(response => {
            response.json().then(content => {
                if (!response.ok) {
                    setError(content.message)
                } else {
                    navigate('/albums')
                }
            });
        }).catch(() => setError("Network communication error. Maybe backend is down?"))
    }, [navigate, path])

    const handleDeleteAlbumDialogAccept = React.useCallback(() => {
        deleteAlbum()
        setDeleteAlbumDialogOpen(false);
    }, [setDeleteAlbumDialogOpen, deleteAlbum]);

    const handleAddClick = React.useCallback(() => {
        navigate(path + '/add', { 
            state: {
                album: data,
                sourcePath: path,
            }
        })
    }, [navigate, path, data]);

    const handleRateClick = React.useCallback(() => {
        navigate(path + '/ratings')
    }, [navigate, path]);

    const handleDeleteClick = React.useCallback((items) => {
        setDeleteDialogOpen(true);
        setDeleteItems(items)
    }, [setDeleteDialogOpen, setDeleteItems]);

    const handleDeleteDialogClose = React.useCallback(() => {
        setDeleteDialogOpen(false);
        setDeleteItems(null);
    }, [setDeleteDialogOpen, setDeleteItems]);

    const deleteFromAlbum = React.useCallback((items) => {
        let keptItems = items.filter(photo => !photo.selected).map(item => ({id: item.id}))
        update(null, null, keptItems)
    }, [update])

    const handleDeleteDialogAccept = React.useCallback(() => {
        deleteFromAlbum(deleteItems)
        setDeleteDialogOpen(false);
        setDeleteItems(null);
    }, [setDeleteDialogOpen, deleteFromAlbum, setDeleteItems, deleteItems]);

    const selectionActions = [
        {
            icon: <DeleteIcon sx={{ color: theme.palette.background.paper }} />,
            tooltip: "Remove photos from album",
            action: handleDeleteClick
        }
    ]

    const albumActions = [
        {
            icon: <EditIcon />,
            tooltip: "Edit album properties",
            action: handleEditAlbumDialogOpen
        },
        {
            icon: <AddPhotoAlternateIcon />,
            tooltip: "Add photos to the album",
            action: handleAddClick
        },
        {
            icon: <DeleteIcon />,
            tooltip: "Delete album - photos are not deleted",
            action: handleDeleteAlbumClick
        }
    ]

    const actions = [
        { icon: <AddIcon />, name: 'Add photos', action: handleAddClick },
        { icon: <StarIcon />, name: 'Rate photos', action: handleRateClick },
        { icon: <EditIcon />, name: 'Edit properties', action: handleEditAlbumDialogOpen },
        { icon: <DeleteIcon />, name: 'Delete album', action: handleDeleteAlbumClick },
    ];

    return (
        <>
            {title &&
                <Stack sx={{ marginBottom: 2, marginTop: 2 }} justifyContent={'center'} alignItems={'baseline'} direction="row" spacing={1}
                    onMouseOver={handleMouseOver} onMouseOut={handleMouseOut}>
                    <Typography variant='h4'>{title}</Typography>
                    {isHovering && 
                        <>
                            {albumActions.map((action => (
                                <Tooltip key={action.tooltip} title={action.tooltip}>
                                    <IconButton onClick={action.action}>
                                        {action.icon}
                                    </IconButton>
                                </Tooltip>
                            )))}
                        </>}
                    <DeleteDialog open={isDeleteAlbumDialogOpen} onCancel={handleDeleteAlbumDialogClose} onDelete={handleDeleteAlbumDialogAccept} name="the album" />
                </Stack>
            }
            {error && <Alert sx={{ mb: 4, maxWidth: "sm", mx: "auto" }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Updated successfully!</Alert>}
            {data && data.tags && <Stack sx={{ marginBottom: 4 }} justifyContent={'center'} direction="row" spacing={1} onClick={handleEditAlbumDialogOpen}>
                {data.tags.map(tag => {
                    return (<Chip key={tag} label={tag} />)
                })}
            </Stack>}
            {data && 
            <>
                <EditAlbumDialog
                        open={isEditAlbumDialogOpen}
                        onSave={handleEditAlbumDialogSave}
                        onCancel={handleEditAlbumDialogClose}
                        input={{ name: data.name, tags: data.tags }}
                     />
                <DeleteDialog open={isDeleteDialogOpen} onCancel={handleDeleteDialogClose} onDelete={handleDeleteDialogAccept} name="the selected photos from the album" />
            </>}
            <PhotoGrid populator={populate} onDataLoaded={handleDataLoaded} selectionActionOverride={selectionActions} />
            <Box sx={{
                '& > :not(style)': { m: 1 },
                position: "fixed",
                bottom: 16,
                right: 16
            }}>
                <SpeedDial
                    ariaLabel="Album actions"
                    sx={{ position: 'absolute', bottom: 16, right: 16 }}
                    icon={<SpeedDialIcon />}
                >
                    {actions.map((action) => (
                        <SpeedDialAction
                            key={action.name}
                            icon={action.icon}
                            tooltipTitle={action.name}
                            onClick={action.action}
                        />
                    ))}
                </SpeedDial>
            </Box>
        </>)
}

AlbumDisplay.propTypes = {
    user: PropTypes.object
};

export default AlbumDisplay;
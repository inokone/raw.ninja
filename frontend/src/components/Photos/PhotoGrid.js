import * as React from 'react';
import PropTypes from "prop-types";
import { Alert, Box } from "@mui/material";
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';
import ProgressDisplay from '../Common/ProgressDisplay';
import DeleteDialog from '../Common/DeleteDialog';
import CollectionsIcon from '@mui/icons-material/Collections';
import DeleteIcon from '@mui/icons-material/Delete';
import SelectableGallery from './SelectableGallery';
import { convertPhotos, convertPhoto } from './PhotoConverter';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoGrid = ({ populator, onDataLoaded, selectionActionOverride, config }) => {
    const theme = useTheme()
    const navigate = useNavigate()
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [images, setImages] = React.useState(null)
    const [isDeleteDialogOpen, setDeleteDialogOpen] = React.useState(false)

    const updateImage = (photo) => {
        let image = photo.base
        image.descriptor.favorite = photo.favorite
        image.descriptor.rating = photo.rating
        fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id, {
            method: "PUT",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(image)
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    let newImages = images.slice()
                    for (let i = 0; i < images.length; i++) {
                        if (images[i].id === image.id) {
                            newImages[i] = convertPhoto(image)
                            setImages(newImages)
                            return
                        }
                    }
                }
            })
            .catch(error => {
                setError(error.message)
            });
    }

    const deletePhoto = React.useCallback((id) => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + id, {
            method: "DELETE",
            mode: "cors",
            credentials: "include"
        }).then(response => {
            if (!response.ok) {
                return new Promise((resolve, reject) => {
                    reject(response.status + ":" + response.statusText)
                })
            }
            navigate(0)
        });
    }, [navigate])

    const batchDelete = React.useCallback(() => {
        let selectedIDs = images.filter((img) => img.selected).map((img) => img.id);
        // hacky way, batch delete backend option would be much better
        selectedIDs.forEach(id => {
            deletePhoto(id)
        })
    },[deletePhoto, images])

    const createAlbum = () => {
        let selectedIDs = images.filter((img) => img.selected).map((img) => img.id);
        navigate('/albums/create', { state: { photos: selectedIDs } })
    }

    const processPhotos = React.useCallback((content) => {
        if (onDataLoaded) {
            onDataLoaded(content)
        }
        let imgs = convertPhotos(content)
        setImages(imgs)
    }, [setImages, onDataLoaded])

    const handleDeleteClick = React.useCallback((photos) => {
        setDeleteDialogOpen(true);
    }, [setDeleteDialogOpen]);

    const handleDeleteDialogClose = React.useCallback(() => {
        setDeleteDialogOpen(false);
    }, [setDeleteDialogOpen]);

    const handleDeleteDialogAccept = React.useCallback(() => {
        batchDelete()
        setDeleteDialogOpen(false);
    }, [setDeleteDialogOpen, batchDelete]);

    const defaultSelectionActions = [
        {
            icon: <CollectionsIcon sx={{ color: theme.palette.background.paper }} />,
            tooltip: "Create album from selection",
            action: createAlbum
        },
        {
            icon: <DeleteIcon sx={{ color: theme.palette.background.paper }} />,
            tooltip: "Delete selected photos",
            action: handleDeleteClick
        }
    ]

    const selectionActions = selectionActionOverride ? selectionActionOverride : defaultSelectionActions

    React.useEffect(() => {
        const loadImages = () => {
            populator()
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            setError(content.message)
                        });
                    } else {
                        response.json().then(content => {
                            setLoading(false)
                            processPhotos(content)
                        })
                    }
                })
                .catch(error => {
                    setError(error.message)
                    setLoading(false)
                });
        }

        if (!images && !error && !loading) {
            loadImages()
        }
    }, [populator, images, error, loading, processPhotos])

    return (
        <Box sx={{ display: 'flex' }}>
            <DeleteDialog open={isDeleteDialogOpen} onCancel={handleDeleteDialogClose} onDelete={handleDeleteDialogAccept} name="the selected photos" />
            {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {images && <SelectableGallery images={images} setImages={setImages} updateImage={updateImage} selectionActionOverride={selectionActions} config={config} />}
        </Box >
    );
}

PhotoGrid.propTypes = {
    populator: PropTypes.func.isRequired,
    onDataLoaded: PropTypes.func,
    selectionActionOverride: PropTypes.array,
    config: PropTypes.object
};

export default PhotoGrid;
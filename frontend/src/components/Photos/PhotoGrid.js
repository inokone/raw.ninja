import * as React from 'react';
import { Alert, Box, Fab } from "@mui/material";
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';
import AddIcon from '@mui/icons-material/Add';
import ProgressDisplay from '../Common/ProgressDisplay';
import DeleteDialog from '../Common/DeleteDialog';
import CollectionsIcon from '@mui/icons-material/Collections';
import DeleteIcon from '@mui/icons-material/Delete';
import SelectableGallery from './SelectableGallery';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoGrid = ({ populator, data, fabAction, onDataLoaded, selectionActionOverride }) => {
    const theme = useTheme()
    const navigate = useNavigate()
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [images, setImages] = React.useState(null)
    const [isDeleteDialogOpen, setDeleteDialogOpen] = React.useState(false)

    const dateOf = (data) => {
        return new Date(data).toLocaleDateString()
    }

    const updateImage = (photo) => {
        let image = photo.base
        image.descriptor.favorite = photo.favorite
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
                            newImages[i] = asPhoto(image)
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

    const formatShutterSpeed = (shutterSpeed) => {
        let validDividers = [2, 4, 8, 15, 30, 60, 125, 250, 500, 1000, 2000, 4000, 8000]
        if (shutterSpeed < 1) {
            let fraction = 1 / shutterSpeed
            let lastDivider = 2
            for (let i = 0; i < validDividers.length; i++) {
                let divider = validDividers[i]
                if (fraction < divider) {
                    if (fraction - lastDivider > divider - fraction) {
                        return "1/" + divider
                    } else {
                        return "1/" + lastDivider
                    }
                }
            }
        } else {
            return shutterSpeed.toFixed(1) + "s";
        }
    }

    const formatBytes = (bytes, decimals = 2) => {
        let negative = (bytes < 0)
        if (negative) {
            bytes = -bytes
        }
        if (!+bytes) return '0 Bytes'

        const k = 1024
        const dm = decimals < 0 ? 0 : decimals
        const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

        const i = Math.floor(Math.log(bytes) / Math.log(k))

        let displayText = `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
        return negative ? "-" + displayText : displayText
    }

    const validTime = (time) => {
        return !time.startsWith("1970-01-01T")
    }

    const description = React.useCallback((image) => {
        let date = dateOf(image.descriptor.metadata.timestamp)
        let timestamp = validTime(image.descriptor.metadata.timestamp) ? ("Taken on " + date + "\n") : ""
        let aperture = image.descriptor.metadata.aperture !== 0 ? "ƒ/" + Math.round((image.descriptor.metadata.aperture + Number.EPSILON) * 100) / 100 + "  " : ""
        let shutter = image.descriptor.metadata.shutter !== 0 ? formatShutterSpeed(image.descriptor.metadata.shutter) + " sec  " : ""
        let iso = image.descriptor.metadata.ISO !== 0 ? "ISO " + image.descriptor.metadata.ISO + "  " : ""
        let dim = image.descriptor.metadata.width + " x " + image.descriptor.metadata.height + " px  "
        let size = formatBytes(image.descriptor.metadata.data_size)
        return timestamp + aperture + shutter + iso + dim + size
    }, [])

    const asPhoto = React.useCallback((image) => {
        return {
            src: image.thumbnail.url,
            original: image.thumbnail.url,
            width: image.descriptor.thumbnail_width ? image.descriptor.thumbnail_width : image.descriptor.metadata.width,
            height: image.descriptor.thumbnail_height ? image.descriptor.thumbnail_height : image.descriptor.metadata.height,
            caption: image.descriptor.filename,
            title: image.descriptor.filename,
            description: description(image),
            favorite: image.descriptor.favorite,
            id: image.id,
            format: image.descriptor.format,
            base: image,
            selected: false
        }
    }
        , [description])

    const processImages = React.useCallback((content) => {
        if (onDataLoaded) {
            onDataLoaded(content)
        }
        if (!Array.isArray(content)) {
            content = content.photos
        }
        let result = content.map(image => asPhoto(image))
        setImages(result)
    }, [onDataLoaded, asPhoto])


    const onDeleteClick = React.useCallback((photos) => {
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
            tooltip: "Create collection from selection",
            action: createAlbum
        },
        {
            icon: <DeleteIcon sx={{ color: theme.palette.background.paper }} />,
            tooltip: "Delete selected photos",
            action: onDeleteClick
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
                            processImages(content)
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
    }, [data, populator, images, error, loading, processImages])

    return (
        <Box sx={{ display: 'flex' }}>
            <DeleteDialog open={isDeleteDialogOpen} onCancel={handleDeleteDialogClose} onDelete={handleDeleteDialogAccept} name="the selected photos" />
            {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {images && <SelectableGallery images={images} setImages={setImages} updateImage={updateImage} selectionActionOverride={selectionActions}/>}
            {fabAction &&
                <Box onClick={fabAction} sx={{
                    '& > :not(style)': { m: 1 },
                    position: "fixed",
                    bottom: 16,
                    right: 16
                }}>
                    <Fab color="primary" aria-label="add">
                        <AddIcon />
                    </Fab>
                </Box>
            }
        </Box >
    );
}

export default PhotoGrid;
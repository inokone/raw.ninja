import * as React from 'react';
import { useLocation, useNavigate } from "react-router-dom";
import { Typography, Alert, Button, Box } from "@mui/material";
import SelectableGallery from '../Photos/SelectableGallery';
import ProgressDisplay from '../Common/ProgressDisplay';
import { convertPhoto, convertPhotos } from '../Photos/PhotoConverter';


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AddPhotos = ({ user }) => {
    const navigate = useNavigate()
    const { pathname, state } = useLocation();
    const path = pathname.slice(0, -4)
    const [images, setImages] = React.useState(null)
    const [error, setError] = React.useState()
    const [loading, setLoading] = React.useState(false)
    const [success, setSuccess] = React.useState(false)

    const processPhotos = React.useCallback((content) => {
        let imgs = convertPhotos(content)
        let selectedIDs = state.album.photos.map(photo => photo.id)
        imgs.forEach((img, idx) => {
            if (selectedIDs.indexOf(img.id) >= 0) {
                imgs[idx].selected = true
            }
        })
        setImages(imgs)
    }, [setImages, state])

    const loadImages = React.useCallback(() => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        }).then(response => {
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
        }).catch(error => {
            setError(error.message)
            setLoading(false)
        });
    }, [user, setLoading, setError, processPhotos])

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

    const updateAlbum = React.useCallback(() => {
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "PATCH",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({
                "photos": images.filter(image => image.selected),
            })
        }).then(response => {
            response.json().then(content => {
                if (!response.ok) {
                    setError(content.message)
                } else {
                    navigate(state.sourcePath);
                }
            });
        }).catch(() => setError("Network communication error. Maybe backend is down?"))
    }, [navigate, path, images, state])

    const handleCancelClick = React.useCallback(() => {
        navigate(state.sourcePath)
    }, [navigate, state])

    const handleAcceptClick = React.useCallback(() => {
        updateAlbum()
    }, [updateAlbum])

    React.useEffect(() => {
        loadImages()
    }, [loadImages])

    return (
        <>
            <Typography variant='h4' sx={{mt: 2, mb: 4}}>Select photos to add</Typography>
            <Button variant='contained' color='primary' onClick={handleAcceptClick} sx={{ mr: 2 }}>Finish</Button>
            <Button variant='contained' color='secondary' onClick={handleCancelClick}>Cancel</Button>
            <Box sx={{ display: 'flex' }}>
            {loading && <ProgressDisplay />}
            {error && <Alert sx={{ mb: 4, maxWidth: "sm", marginLeft: "auto", marginRight: "auto" }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Updated successfully!</Alert>}
            {images && <SelectableGallery images={images} setImages={setImages} updateImage={updateImage} />}
            </Box>
        </>)
}

export default AddPhotos;
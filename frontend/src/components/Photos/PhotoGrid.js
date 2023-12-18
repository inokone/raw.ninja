import * as React from 'react';
import { Alert } from "@mui/material";
import ProgressDisplay from '../Common/ProgressDisplay';
import PhotoGallery from './PhotoGallery';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoGrid = ({ populator, data }) => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [images, setImages] = React.useState(null)

    const dateOf = (data) => {
        return new Date(data).toLocaleDateString()
    }

    const setPhoto = (photo) => {
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

    const setSelected = (photo) => {
        let newImages = images.slice()
        newImages.forEach(i => {
            if (i.id === photo.id) {
                i.selected = photo.selected
            } 
        });
        setImages(newImages)
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

    const description = (image) => {
        let timestamp = "Taken on " + dateOf(image.descriptor.metadata.timestamp)
        let aperture = image.descriptor.metadata.aperture !== 0 ? "\nƒ/" + Math.round((image.descriptor.metadata.aperture + Number.EPSILON) * 100) / 100 : ""
        let shutter = image.descriptor.metadata.shutter !== 0 ? " - " + formatShutterSpeed(image.descriptor.metadata.shutter) + " sec" : ""
        let iso = image.descriptor.metadata.ISO !== 0 ? " - ISO " + image.descriptor.metadata.ISO : ""
        return timestamp + aperture + shutter + iso
    }

    const processImages = (content) => {
        if (!Array.isArray(content)) {
            content = content.photos
        }
        let result = content.map(image => asPhoto(image))
        setImages(result)
    }

    const asPhoto = (image) => {
        return {
            src: image.thumbnail.url,
            original: image.thumbnail.url,
            width: image.descriptor.metadata.width,
            height: image.descriptor.metadata.height,
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
    }, [data, populator, images, error, loading])

    return (
        <React.Fragment>
            {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay />}
            {images && <PhotoGallery photos={images} setPhoto={setPhoto} setSelected={setSelected}/>}
        </React.Fragment>
    );
}

export default PhotoGrid;
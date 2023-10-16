import * as React from 'react';
import { CircularProgress, Alert, Grid, Typography } from "@mui/material";
import PhotoCard from '../Photos/PhotoCard';


const { REACT_APP_API_PREFIX } = process.env;

const PhotoGrid = (props) => {
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(true)
    const [images, setImages] = React.useState(null)
    const gridPopulator = props.populator

    React.useEffect(() => {
        const loadImages = () => {
            gridPopulator()
                .then(response => {
                    if (!response.ok) {
                        if (response.status !== 200) {
                            setError(response.status + ": " + response.statusText);
                        } else {
                            response.json().then(content => setError(content.message))
                        }
                        setLoading(false)
                    } else {
                        response.json().then(content => {
                            setLoading(false)
                            setImages(content)
                        })
                    }
                })
                .catch(error => {
                    setError(error)
                    setLoading(false)
                });
        }
        loadImages()
    }, [gridPopulator, props.data])

    const setImage = (image) => {
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
                    if (response.status !== 200) {
                        setError(response.status + ": " + response.statusText);
                    } else {
                        response.json().then(content => setError(content.message))
                    }
                } else {
                    let newImages = images.slice()
                    for (let i = 0; i < images.length; i++) {
                        if (images[i].id === image.id) {
                            newImages[i] = image
                            setImages(newImages)
                            return
                        }
                    }
                }
            })
            .catch(error => {
                setError(error)
            });
    }

    const dateOf = (image) => {
        return new Date(image.descriptor.uploaded).toLocaleDateString()
    }

    return (
        <>
            {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
            {loading ? <CircularProgress /> :
                <Grid container spacing={1} sx={{ flexGrow: 1, pl: 2, pt: 3 }} >
                    {images.map((image, index) => {
                        if (index === 0 || dateOf(image) !== dateOf(images[index - 1])) {
                            return (
                                <>
                                    <Grid item key={dateOf(image)} lg={12} sx={{textAlign: 'left'}}>
                                        <Typography variant='h5' mt={3}>{dateOf(image)}</Typography>
                                    </Grid>
                                    <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                                        <PhotoCard image={image} setImage={setImage} />
                                    </Grid>
                                </>
                            );
                        }
                        return (
                            <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                                <PhotoCard image={image} setImage={setImage} />
                            </Grid>
                        );
                    })}
                </Grid>}
        </>
    );
}

export default PhotoGrid;
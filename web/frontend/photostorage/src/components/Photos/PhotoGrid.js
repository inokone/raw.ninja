import * as React from 'react';
import { Alert, Grid, Typography } from "@mui/material";
import PhotoCard from '../Photos/PhotoCard';
import DetailedPhotoCard from './DetailedPhotoCard';
import useMediaQuery from '@mui/material/useMediaQuery';
import { useTheme } from '@mui/material/styles';
import { useNavigate } from "react-router-dom"
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const PhotoGrid = ({ populator, data }) => {
    const theme = useTheme();
    const navigate = useNavigate();
    const [error, setError] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [images, setImages] = React.useState(null)
    const [selected, setSelected] = React.useState(null)
    const isSmScreen = useMediaQuery(theme.breakpoints.down('sm'));
    const isMdScreen = useMediaQuery(theme.breakpoints.between('sm', 'md'));
    const isLgScreen = useMediaQuery(theme.breakpoints.between('md', 'lg'));
    const isXlScreen = useMediaQuery(theme.breakpoints.up('lg'));

    React.useEffect(() => {
        const loadImages = () => {
            populator()
                .then(response => {
                    if (!response.ok) {
                        throw new Error(response.status + ": " + response.statusText);
                    } else {
                        response.json().then(content => {
                            setLoading(false)
                            setImages(content)
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
    }, [populator, data, images, error, loading])


    const handleImageClick = (index) => {
        if (isLgScreen || isXlScreen) {
            setSelected(index)
        } else {
            navigate("/photos/" + images[index].id)
        }
    }

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
                    throw new Error(response.status + ": " + response.statusText);
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
                setError(error.message)
            });
    }

    const dateOf = (image) => {
        return new Date(image.descriptor.uploaded).toLocaleDateString()
    }

    const dateHeader = (image) => {
        return (
            <Grid item key={dateOf(image)} xs={72} sm={48} md={24} lg={12} sx={{ textAlign: 'left' }}>
                <Typography variant='h5' mt={2}>{dateOf(image)}</Typography>
            </Grid>
        )
    };

    const card = (image, index, selected) => {
        return (
            <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                <PhotoCard image={image} setImage={setImage} selected={selected} onClick={() => handleImageClick(index)} />
            </Grid>
        )
    };

    const cardWithHeader = (images, image, index, isSelected) => {
        if (index === 0 || dateOf(image) !== dateOf(images[index - 1])) {
            return (
                <React.Fragment>
                    {dateHeader(image)}
                    {card(image, index, isSelected)}
                </React.Fragment>
            );
        }
        return card(image, index, isSelected);
    };

    const selection = (image) => {
        return (
            <React.Fragment>            
                {selected !== null ?
                <Grid item key={"selection"} xs={72} sm={48} md={24} lg={12}>
                    <DetailedPhotoCard image={image} onClose={() => setSelected(null)} setImage={setImage} closable={true} />
                </Grid> : null
                }
            </React.Fragment>)
    }

    const isPlaceForSelection = (index) => {
        // if no selection, or we are before selection 
        if (!selected || index < selected) {
            return false
        }
        let dividers = rowDividers()
        for (let i = 0; i < dividers.length; i++) {
            // skip for dividers before selection
            if (index !== dividers[i] || dividers[i] <= selected) {
                continue
            }
            // if we are on a row break and selected was in the previous row, we have found the place 
            if (i === 0 || dividers[i - 1] <= selected) {
                return true
            }
        }
        return false
    }

    const rowDividers = () => {
        if (images.length < 1) {
            return []
        }
        let result = []
        let date = dateOf(images[0])
        let rowCounter = 0
        let rowLength = thumbnailsInARow()
        for (let i = 0; i < images.length; i++) {
            if (rowCounter === rowLength) {
                rowCounter = 1
                result.push(i)
            } else {
                rowCounter++
                if (date !== dateOf(images[i])) {
                    date = dateOf(images[i])
                    rowCounter = 1
                    result.push(i)
                }
            }
        }
        return result
    }

    const isInLastRow = (index) => {
        return rowDividers().slice(-1) <= index
    }

    const thumbnailsInARow = () => {
        if (isSmScreen) {
            return 2
        }
        if (isMdScreen) {
            return 3
        }
        if (isLgScreen) {
            return 4
        }
        if (isXlScreen) {
            return 6
        }
    }

    return (
        <React.Fragment>
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {loading && <ProgressDisplay /> }
            {images && <Grid container spacing={1} sx={{ flexGrow: 1, pl: 2, pt: 3 }} >
                    {images.map((image, index) => {
                        if (isPlaceForSelection(index)) {
                            return (
                                <React.Fragment key={index}>
                                    {selection(images[selected])}
                                    {cardWithHeader(images, image, index, selected === index)}
                                </React.Fragment>
                            )
                        }
                        if (index === images.length - 1 && isInLastRow(selected)) {
                            return (
                                <React.Fragment key={index}>
                                    {cardWithHeader(images, image, index, selected === index)}
                                    {selection(images[selected])}
                                </React.Fragment>
                            )
                        }
                        return (
                            <React.Fragment key={index}>
                                {cardWithHeader(images, image, index, selected === index)}
                            </React.Fragment>
                        )
                    })}
                </Grid>}
        </React.Fragment>
    );
}

export default PhotoGrid;
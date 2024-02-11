import * as React from 'react';
import { Rating } from "@mui/material";
import PropTypes from "prop-types";

import Lightbox from "yet-another-react-lightbox";
import { useLightboxState } from "yet-another-react-lightbox";
import "yet-another-react-lightbox/styles.css";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Captions from "yet-another-react-lightbox/plugins/captions";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import "yet-another-react-lightbox/plugins/captions.css";
import StarIcon from "@mui/icons-material/Star";

function RatingComponent({ setIndex, onRatingChange }) {
    const { currentSlide, currentIndex } = useLightboxState();
    return (
        <>
            <Rating
                key="rating"
                sx={{ pt: 1, mr: 3 }}
                size="large"
                value={currentSlide.rating}
                onChange={(_, rating) => {
                    onRatingChange(currentSlide, rating)
                    setIndex(currentIndex)
                }}
                emptyIcon={<StarIcon style={{ opacity: 0.55, color: 'white' }} fontSize="inherit" />}
            />
        </>);
}

const RatingLightbox = ({ photos, index, setIndex, onRatingChange }) => {
    return (<Lightbox
        slides={photos}
        open={index >= 0}
        onRatingChange={onRatingChange}
        index={index}
        setIndex={setIndex}
        close={() => setIndex(-1)}
        toolbar={{
            buttons: [
                <RatingComponent onRatingChange={onRatingChange} setIndex={setIndex} />,
                "zoom", 
                "close"]
        }}
        plugins={[Captions, Thumbnails, Zoom]}
    />)
}

RatingLightbox.propTypes = {
    photos: PropTypes.array.isRequired,
    index: PropTypes.number.isRequired,
    setIndex: PropTypes.func.isRequired,
    onRatingChange: PropTypes.func.isRequired,
    onDeleteClick: PropTypes.func.isRequired,
    onEditClick: PropTypes.func.isRequired
};

export default RatingLightbox;
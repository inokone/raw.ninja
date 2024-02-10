import * as React from 'react';
import { Rating } from "@mui/material";
import PropTypes from "prop-types";

import Lightbox from "yet-another-react-lightbox";
import "yet-another-react-lightbox/styles.css";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Captions from "yet-another-react-lightbox/plugins/captions";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import "yet-another-react-lightbox/plugins/captions.css";
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import StarIcon from "@mui/icons-material/Star";

const RatingLightbox = ({ photos, config, index, setIndex, onRatingChange, onDeleteClick, onEditClick}) => {
    const [buttons, setButtons] = React.useState([])
    
    const createButtons = React.useCallback(() => {
        const defaultConfig = {
            editingEnabled: true,
            ratingEnabled: false,
            deletingEnabled: true,
        }
        let settings = { ...defaultConfig, ...config }
        console.log(settings)
        let buttons = []
        if (settings.ratingEnabled) {
            buttons.push(
                <Rating
                    key="rating"
                    sx={{ pt: 1, mr: 3 }}
                    size="large"
                    value={photos && photos[index] ? photos[index].rating : 0}
                    onChange={(_, rating) => onRatingChange(photos[index], rating)}
                    emptyIcon={<StarIcon style={{ opacity: 0.55, color: 'white' }} fontSize="inherit" />}
                />
            )
        }
        buttons.push("zoom")
        if (settings.editingEnabled) {
            buttons.push(
                <button key="edit" type="button" className="yarl__button" onClick={() => onEditClick(photos[index])}>
                    <EditIcon />
                </button>)
        }
        if (settings.deletingEnabled) {
            buttons.push(<button key="delete" type="button" className="yarl__button" onClick={() => onDeleteClick(photos[index])}>
                <DeleteIcon />
            </button>)
        }
        buttons.push("close")
        return buttons
    },[config, index, onDeleteClick, onEditClick, onRatingChange, photos])


    React.useEffect(() => {
        setButtons(createButtons())
    },[createButtons]);


    return (<Lightbox
        slides={photos}
        open={index >= 0}
        index={index}
        close={() => setIndex(-1)}
        toolbar={{
            buttons: buttons
        }}
        plugins={[Captions, Thumbnails, Zoom]}
    />)
}

RatingLightbox.propTypes = {
    photos: PropTypes.array.isRequired,
    config: PropTypes.object,
    index: PropTypes.number.isRequired, 
    setIndex: PropTypes.func.isRequired,
    onRatingChange: PropTypes.func.isRequired,
    onDeleteClick: PropTypes.func.isRequired,
    onEditClick: PropTypes.func.isRequired
};

export default RatingLightbox;
import * as React from 'react';
import PropTypes from "prop-types";

import YALightbox from "yet-another-react-lightbox";
import { useLightboxState } from "yet-another-react-lightbox";
import "yet-another-react-lightbox/styles.css";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Captions from "yet-another-react-lightbox/plugins/captions";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import "yet-another-react-lightbox/plugins/captions.css";
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

const EditButton = ({onEditClick}) => {
    const { currentSlide } = useLightboxState();

    return (
        <button key="edit" type="button" className="yarl__button" onClick={() => onEditClick(currentSlide)}>
            <EditIcon />
        </button>
    )
}

const DeleteButton = ({ onDeleteClick }) => {
    const { currentSlide } = useLightboxState();

    return (
        <button key="delete" type="button" className="yarl__button" onClick={() => onDeleteClick(currentSlide)}>
            <DeleteIcon />
        </button>
    )
}

const Lightbox = ({ photos, config, index, setIndex, onDeleteClick, onEditClick}) => {
    const [buttons, setButtons] = React.useState([])
    
    const createButtons = React.useCallback(() => {
        const defaultConfig = {
            editingEnabled: true,
            ratingEnabled: false,
            deletingEnabled: true,
        }
        let settings = { ...defaultConfig, ...config }
        let buttons = ["zoom"]
        if (settings.editingEnabled) {
            buttons.push(<EditButton onEditClick={onEditClick}/>)
        }
        if (settings.deletingEnabled) {
            buttons.push(<DeleteButton onDeleteClick={onDeleteClick}/>)
        }
        buttons.push("close")
        return buttons
    },[config, onDeleteClick, onEditClick])


    React.useEffect(() => {
        setButtons(createButtons())
    },[createButtons]);


    return (<YALightbox
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

Lightbox.propTypes = {
    photos: PropTypes.array.isRequired,
    config: PropTypes.object,
    index: PropTypes.number.isRequired, 
    setIndex: PropTypes.func.isRequired,
    onDeleteClick: PropTypes.func.isRequired,
    onEditClick: PropTypes.func.isRequired
};

export default Lightbox;
import * as React from 'react';
import PropTypes from "prop-types";
import { styled } from '@mui/material/styles';
import { useTheme } from '@mui/material/styles';
import PhotoGallery from './PhotoGallery';
import SelectionActionBar from './SelectionActionBar';
import ClearIcon from '@mui/icons-material/Clear';

const drawerWidth = 60;

const Main = styled('main', { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
        flexGrow: 1,
        padding: theme.spacing(1),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: `-${drawerWidth}px`,
        ...(open && {
            transition: theme.transitions.create('margin', {
                easing: theme.transitions.easing.easeOut,
                duration: theme.transitions.duration.enteringScreen,
            }),
            marginLeft: 0,
        }),
    }),
);

const SelectableGallery = ({ images, setImages, updateImage, selectionActionOverride }) => {
    const theme = useTheme()
    const [isSelectionBarOpen, setSelectionBarOpen] = React.useState(false)

    const setSelected = (image) => {
        let newImages = images.slice()
        let selectedCount = 0
        newImages.forEach(i => {
            if (i.id === image.id) {
                i.selected = image.selected
            }
            if (i.selected) {
                selectedCount++
            }
        });
        setImages(newImages)
        setSelectionBarOpen(selectedCount > 0)
    }

    const clearSelection = () => {
        let newImages = images.slice()
        newImages.forEach(i => { i.selected = false });
        setImages(newImages)
        setSelectionBarOpen(false)
    }

    const clearAction = {
        icon: <ClearIcon sx={{ color: theme.palette.background.paper }} />,
        tooltip: "Clear selection",
        action: clearSelection
    }

    const selectionActions = selectionActionOverride ? selectionActionOverride.concat([clearAction]) : [clearAction]

    return (
        <>
            <SelectionActionBar open={isSelectionBarOpen} items={images} actions={selectionActions} setSelected={setSelected}/>
            <Main open={isSelectionBarOpen}>
                {images && <PhotoGallery photos={images} updatePhoto={updateImage} setSelected={setSelected} />}
            </Main>
        </>
    );
}

SelectableGallery.propTypes = {
    images: PropTypes.object.isRequired,
    setImages: PropTypes.func.isRequired,
    updateImage: PropTypes.func.isRequired,
    selectionActionOverride: PropTypes.array
};

export default SelectableGallery;
import * as React from 'react';
import PropTypes from "prop-types";
import StarIcon from '@mui/icons-material/Star';
import EditIcon from '@mui/icons-material/Edit';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import Brightness1Icon from '@mui/icons-material/Brightness1';
import { Box, IconButton, Tooltip, Typography, Rating } from "@mui/material";
import { makeStyles } from '@mui/styles';
import { useNavigate } from 'react-router-dom';

/*
Schema of "config":
{
  displayTitle: true/false,
  selectionEnabled: true/false,
  favoriteEnabled: true/false,
  editingEnabled: true/false,
  ratingEnabled: true/false,
}
*/

const useStyles = makeStyles(() => ({
  favorite: {
    color: 'white',
    cursor: "pointer",
  },
  nonfavorite: {
    color: 'lightgray',
    cursor: "pointer",
  },
  selected: {
    position: 'absolute',
    fill: "#06befa",
    cursor: "pointer",
    overflow: "hidden",
    width: "32px",
    height: "32px",
  },
  unselected: {
    position: 'absolute',
    color: "lightgray",
    cursor: "pointer",
  },
  selected_bg: {
    position: 'absolute',
    fill: "white",
    cursor: "pointer",
  },
  unselected_bg: {
    display: 'none',
  },
  selected_box: {
    marginTop: '4px',
    marginLeft: '4px',
    width: "32px",
    height: "32px",
    position: 'relative',
  },
  unselected_box: {
    marginTop: '8px',
    marginLeft: '8px',
    width: "24px",
    height: "24px",
    color: 'lightgray',
    position: 'relative'
  },

}));

const PhotoCard = ({ photo, updatePhoto, setSelected, onClick, config, imageProps: { src, alt, style, width, height, ...restImageProps } }) => {
  const defaultConfig = {
    displayTitle: false,
    selectionEnabled: true,
    favoriteEnabled: true,
    editingEnabled: true,
    ratingEnabled: false,
  }
  const settings = { ...defaultConfig, ...config }
  const navigate = useNavigate();
  const classes = useStyles();
  const [isHovering, setIsHovering] = React.useState(false);

  const handleMouseOver = () => {
    setIsHovering(true);
  };

  const handleMouseOut = () => {
    setIsHovering(false);
  };

  const handleFavoriteClick = (photo) => {
    const updatedPhoto = { ...photo };
    updatedPhoto.favorite = !photo.favorite
    updatePhoto(updatedPhoto)
  }

  const handleRatingChanged = (rating) => {
    const updatedPhoto = { ...photo };
    updatedPhoto.rating = rating
    updatePhoto(updatedPhoto)
  }

  const handleSelectClick = (photo) => {
    photo.selected = !photo.selected
    setSelected(photo)
  }

  const imgStyle = {
    transition: "transform .135s cubic-bezier(0.0,0.0,0.2,1),opacity linear .15s"
  };
  
  const selectedImgStyle = {
    transform: "translateZ(0px) scale3d(0.9, 0.9, 1)",
    transition: "transform .135s cubic-bezier(0.0,0.0,0.2,1),opacity linear .15s"
  };

  const handleClick = (id) => {
    onClick(id)
  }

  const handleEditClick = () => {
    navigate('/editor/' + photo.id + '?format=' + photo.format, {
      state: {
        photo_id: photo.id,
        photo_format: photo.format,
        photo_name: photo.title
      }
    })
  }

  return (
    <Box sx={{ position: 'relative' }} style={style} onMouseOver={handleMouseOver} onMouseOut={handleMouseOut}>
      <img
        src={src}
        alt={alt}
        width='100%'
        height='100%'
        style={
          photo.selected ? { ...imgStyle, ...selectedImgStyle } : { ...imgStyle }
        }
        {...restImageProps}
        onClick={() => handleClick(photo)}
      />
      {settings.displayTitle &&
        <Box sx={{
          width: '100%',
          background:
            'linear-gradient(to top, rgba(0,0,0, 1) 0%, ' +
            'rgba(0,0,0, 1) 30%, rgba(0,0,0, .35) 100%)',
          position: 'absolute',
          bottom: 0,
          right: 0,
        }}>
          <Typography sx={{ color: 'white', marginBottom: '2px', marginTop: '2px' }}>{photo.title}</Typography>
        </Box>
      }
      {((settings.favoriteEnabled && photo.favorite) || ((settings.favoriteEnabled || settings.editingEnabled) && isHovering) || settings.ratingEnabled) && <Box sx={{
        background:
          'linear-gradient(to bottom, rgba(0,0,0,0.7) 0%, ' +
          'rgba(0,0,0,0.3) 70%, rgba(0,0,0,0) 100%)',
        position: 'absolute',
        width: '100%',
        textAlign: 'right',
        top: 0,
        right: 0,
      }}>
        {settings.editingEnabled && isHovering &&
          <Tooltip title="Edit photo">
            <IconButton aria-label="Edit this photo" onClick={() => handleEditClick(photo)} sx={{ color: 'lightgray' }}>
              <EditIcon />
            </IconButton>
          </Tooltip>
        }
        {settings.favoriteEnabled &&
          <Tooltip title="Mark as favorite">
            <IconButton aria-label="Add to favorites" onClick={() => handleFavoriteClick(photo)}>
              <StarIcon className={photo.favorite ? classes.favorite : classes.nonfavorite} />
            </IconButton>
          </Tooltip>
        }
        {settings.ratingEnabled &&
          <Box sx={{ width: '100%', display: 'flex' }}>
            <Tooltip title="Rate photo">
              <Rating
                key="rating"
                sx={{ pt: 1, mx: 'auto'}}
                value={photo.rating}
                onChange={(_, rating) => handleRatingChanged(rating)}
                emptyIcon={<StarIcon style={{ opacity: 0.55, color: 'white' }} fontSize="inherit" />}
              />
            </Tooltip>
          </Box>
        }
      </Box>}
      {settings.selectionEnabled && (isHovering || photo.selected) && <Box sx={{
        position: 'absolute',
        textAlign: 'left',
        top: 0,
        left: 0,
      }}>
        <Tooltip title="Select photo">
          <IconButton aria-label="Select photo" onClick={() => handleSelectClick(photo)} className={photo.selected ? classes.selected_box : classes.unselected_box}>
            <span // Stacking 2 Icons over each other so it looks better
              style={{
                display: "flex",
                flexDirection: "column",
                justifyContent: "Center",
                alignItems: "Center"
              }}
            >
              <Brightness1Icon className={photo.selected ? classes.selected_bg : classes.unselected_bg} />
              <CheckCircleIcon className={photo.selected ? classes.selected : classes.unselected} />
            </span>
          </IconButton>
        </Tooltip>
      </Box>}
    </Box>
  );
}

PhotoCard.propTypes = {
  photo: PropTypes.object.isRequired,
  config: PropTypes.object,
  updatePhoto: PropTypes.func.isRequired,
  setSelected: PropTypes.func.isRequired,
  onClick: PropTypes.func.isRequired
};

export default PhotoCard;
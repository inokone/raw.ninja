import * as React from 'react';
import DownloadIcon from '@mui/icons-material/Download';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { Card, CardMedia, Box, CardActions, CardActionArea, IconButton, Tooltip } from "@mui/material";
import { makeStyles } from '@mui/styles';

const { REACT_APP_API_PREFIX } = process.env;

const useStyles = makeStyles((theme) => ({
  favorite: {
    color: 'red',
  },
  nonfavorite: {
    color: 'lightgray',
  },
}));

const PhotoCard = ({ image, setImage, selected, onClick }) => {
  const classes = useStyles();

  const handleFavoriteClick = (image) => {
    const updatedImage = { ...image };
    updatedImage.descriptor.favorite = !image.descriptor.favorite
    setImage(updatedImage)
  }

  const handleDownloadClick = (image) => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id + '/download', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
      .then(response => response.blob())
      .then(blob => {
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.setAttribute(
          'download',
          image.descriptor.filename,
        );
        document.body.appendChild(a);
        a.click();
        a.remove();
      });
  }

  const handleClick = (id) => {
    onClick(id)
  }

  return (
    <Card sx={{ maxWidth: 250, position: 'relative', cursor: "pointer" }}>
      <Box>
        <CardMedia
          component="img"
          height="200px"
          image={image.descriptor.thumbnail}
          loading="lazy"
          alt={image.descriptor.filename}
          onClick={() => handleClick(image.id)}
        />
      </Box>
      <CardActionArea component="div" sx={{
        ...(!selected && {
        background:
          'linear-gradient(to top, rgba(0,0,0,0.7) 0%, ' +
          'rgba(0,0,0,0.3) 70%, rgba(0,0,0,0) 100%)',
        }),
        ...(selected && {
          background:
            'linear-gradient(to top, rgba(0,255,0,0.9) 0%, ' +
            'rgba(0,255,0,0.9) 10%, rgba(0,0,0,0.65) 11%, rgba(0,0,0,0.3) 70%, rgba(0,0,0,0) 100%)',
        }),
        position: 'absolute',
        bottom: 0,
        right: 0,
        width: '100%',
      }}>
        <CardActions disableSpacing>
          <Tooltip title="Prevent file from being deleted by life-cycle rules">
            <IconButton aria-label="Add to favorites" onClick={() => handleFavoriteClick(image)}>
              <FavoriteIcon className={image.descriptor.favorite ? classes.favorite : classes.nonfavorite} />
            </IconButton>
          </Tooltip>
          <Tooltip title="Download RAW file">
            <IconButton aria-label="Download RAW image" onClick={() => handleDownloadClick(image)} sx={{ color: 'lightgray' }}>
              <DownloadIcon />
            </IconButton>
          </Tooltip>
        </CardActions>
      </CardActionArea>
    </Card>
  );
}
export default PhotoCard;
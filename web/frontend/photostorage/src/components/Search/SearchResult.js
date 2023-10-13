import * as React from 'react';
import { CircularProgress, Alert, Grid } from "@mui/material";
import PhotoCard from '../Photos/PhotoCard';


const { REACT_APP_API_PREFIX } = process.env;

const SearchResult = (props) => {
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(true)
  const [images, setImages] = React.useState(null) 

  const loadImages = () => {
    let url = REACT_APP_API_PREFIX + '/api/v1/search?query=' + props.query;
    fetch(url, {
        method: "GET",
        mode: "cors",
        credentials: "include"
    })
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

  React.useEffect(() => {
      loadImages()
  }, [props.query])

  return (
    <>
      {error !== null ? <Alert sx={{mb: 4}} severity="error">{error}</Alert>:null}
      {loading ? <CircularProgress /> : 
        <Grid container spacing={1} sx={{ flexGrow: 1,  pl: 2, pt: 3 }} >
          {images.map((image) => {
            return (
              <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                <PhotoCard id={image.id} source={image.descriptor.thumbnail} filename={image.descriptor.filename} date={image.descriptor.uploaded}/>
              </Grid>
            );
          })}
        </Grid>}
    </>
  );
}

export default SearchResult;
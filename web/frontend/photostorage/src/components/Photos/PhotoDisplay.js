import React from 'react';
import Image from 'mui-image'
import { CircularProgress, Alert, Button, Tooltip } from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import { useLocation } from "react-router-dom"
import MetadataDisplay from './MetadataDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const PhotoDisplay = () => {
  const location = useLocation()
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(true)
  const [source, setSource] = React.useState("dummy.png")
  const [metadata, setMetadata] = React.useState({})

  const handleDownloadClick = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1' + location.pathname + '/download', {
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
          metadata.filename,
        );
        document.body.appendChild(a);
        a.click();
        a.remove();
      });
  }

  React.useEffect(() => {
    const loadImage = () => {
      fetch(REACT_APP_API_PREFIX + '/api/v1' + location.pathname, {
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
              setSource(content.descriptor.thumbnail)
              setMetadata(content.descriptor)
            })
          }
        })
        .catch(error => {
          setError(error)
          setLoading(false)
        });
    }

    loadImage()
  }, [])

  return (
    <>
      {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
      {loading ? <CircularProgress /> :
        <Grid container spacing={2} padding={3}>
          <Grid xs={9}>
            <Image src={source} height="80vh" />
          </Grid>
          <Grid xs={3}>
            <MetadataDisplay metadata={metadata} />
            <Tooltip title="Download RAW file">
              <Button variant='contained' color='primary' onClick={handleDownloadClick} mt={2}>Download</Button>
            </Tooltip>
          </Grid>
        </Grid>}
    </>
  );
}
export default PhotoDisplay; 
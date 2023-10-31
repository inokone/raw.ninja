import * as React from "react";
import "./Upload.css";
import { DropzoneArea } from "react-mui-dropzone";
import { Alert, Container, Box, Button } from "@mui/material";
import { useNavigate } from "react-router-dom"
import { createStyles, makeStyles } from '@mui/styles';
import ProgressDisplay from "../Common/ProgressDisplay";


const { REACT_APP_API_PREFIX } = process.env;

const useStyles = makeStyles(theme => createStyles({
  previewChip: {
    minWidth: 160,
    maxWidth: 210,
    borderWidth: '2px',
    background: '#bbb',
    color: '#0056b2',
    fontWeight: 'bold',
  },
}));

const Upload = () => {
  const navigate = useNavigate()
  const [stage, setStage] = React.useState(0)
  const [error, setError] = React.useState()
  const [files, setFiles] = React.useState([])
  const classes = useStyles();

  const handleChange = (files) => {
    setFiles(files)
  }

  const handleClick = () => {
    if (files.length === 0) {
      return
    }
    setStage(1)
    var data = new FormData()

    for (const file of files) {
      data.append('files[]', file, file.name);
    }

    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
      method: "POST",
      mode: "cors",
      credentials: "include",
      body: data
    })
      .then(response => {
        console.log(response)
        if (!response.ok) {
          response.json().then(content => {
            setError(content.message)
          });
        } else {
          response.json().then(content => {
            let first = content.photo_ids[0]
            setStage(2)
            setError(null)
            navigate("/photos/" + first)
          })
        }
      })
      .catch(error => {
        setError(error.message)
        setStage(3)
      });
  }

  return (
    <React.Fragment>
      {stage === 0 ?
        <Container>
          <Box m={5}>
            <DropzoneArea 
              m={5} 
              variant="l" 
              filesLimit={20}
              onChange={handleChange}
              acceptedFiles={[".dng, .arw, .cr2, .crw, .nef, .orf, .raf, .jpg, .jpeg, .png"]}
              maxFileSize={100000000} sx={{ flexGrow: 1 }}
              showPreviews={true}
              showPreviewsInDropzone={false}
              useChipsForPreview
              previewGridProps={{ container: { spacing: 1, direction: 'row' } }}
              previewChipProps={{ classes: { root: classes.previewChip } }}
              previewText=""
            />
          </Box>
          {files.length > 0 ?
          <Button onClick={handleClick} variant="contained">Upload</Button>: null}
        </Container> : null}
      {stage === 1 ? <ProgressDisplay text="Please wait for processing the uploaded photos..." /> : null}
      {stage === 2 ? <Alert sx={{ mb: 4 }} severity="success">Upload successful!</Alert> : null}
      {stage === 3 ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
    </React.Fragment>
  );
}

export default Upload;
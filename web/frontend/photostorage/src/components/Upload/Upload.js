import * as React from "react";
import { DropzoneArea } from "react-mui-dropzone";
import { Alert, Container, Box, Button } from "@mui/material";
import { useNavigate } from "react-router-dom"
import { createStyles, makeStyles } from '@mui/styles';
import CloudCircleIcon from '@mui/icons-material/CloudCircle'
import ProgressDisplay from "../Common/ProgressDisplay";


const { REACT_APP_API_PREFIX } = process.env;

const useStyles = makeStyles(theme => createStyles({
  previewChip: {
    minWidth: 160,
    maxWidth: 210,
  },
  dropZone: {
    color: theme.palette.secondary.main,
    width: '80%',
    minHeight: '600px',
    boxSizing: 'border-box',
  },
  uploadIcon: {
    color: theme.palette.secondary.main,
  }
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
      {stage === 0 &&
        <Container>
          <Box m={5}>
            <DropzoneArea 
              m={5}  
              filesLimit={20}
              onChange={handleChange}
              acceptedFiles={[".dng, .arw, .cr2, .crw, .nef, .orf, .raf, .jpg, .jpeg, .png, .gif"]}
              maxFileSize={100000000} sx={{ flexGrow: 1 }}
              showPreviews={true}
              showPreviewsInDropzone={false}
              useChipsForPreview
              previewGridProps={{ container: { spacing: 1, direction: 'row' } }}
              previewChipProps={{ classes: { root: classes.previewChip } }}
              previewText=""
              showAlerts={false}
              Icon={CloudCircleIcon}
            />
          </Box>
          {files.length > 0 && <Button onClick={handleClick} variant="contained">Upload</Button>}
        </Container>
      }
      {stage === 1 && <ProgressDisplay text="Please wait for processing the uploaded photos..." />}
      {stage === 2 && <Alert sx={{ mb: 4 }} severity="success">Upload successful!</Alert>}
      {stage === 3 && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
    </React.Fragment>
  );
}

export default Upload;
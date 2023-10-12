import * as React from "react";
import { DropzoneArea } from "react-mui-dropzone";
import { CircularProgress, Alert } from "@mui/material";
import { useNavigate } from "react-router-dom"


const { REACT_APP_API_PREFIX } = process.env;

const Upload = () => {
  const navigate = useNavigate()
  const [stage, setStage] = React.useState(0)
  const [error, setError] = React.useState()

  const handleChange = (files) => {
    if (files.length === 0) {
      return
    }
    setStage(1)
    console.log(files)
    var data = new FormData()
    data.append('file', files[0])
    data.append('path', files[0].path)

    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
                method: "POST",
                mode: "cors",
                credentials: "include",
                body: data
            })
    .then(response => {
        if (!response.ok) {
            if (response.status !== 200) {
              setStage(3)
              setError(response.status + ": " + response.statusText);
            } else {
              setStage(3)
              response.json().then(content => setError(content.message))
            }
        } else {
            response.json().then(content => {
              let photoId = content.photoId
              setStage(2)
              setError(null)
              navigate("/photos/" + photoId)
            })
        }
    })
    .catch(error => {
      setError(error)
      setStage(3)
    });
  }

  return (
    <>
    {stage === 0 ?
      <DropzoneArea
        onChange={handleChange}
        acceptedFiles={[".dng, .arw, .cr2, .crw, .nef, .orf, .jpg, .jpeg, .png"]} 
        maxFileSize={100000000} sx={{ flexGrow: 1 }}
      />: null }
    {stage === 1 ? <CircularProgress /> : null }
    {stage === 2 ? <Alert sx={{mb: 4}} severity="success">Upload successful!</Alert>:null}
    {stage === 3 ? <Alert sx={{mb: 4}} severity="error">{error}</Alert>:null}
    </>
  );
}

export default Upload;
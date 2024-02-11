

import * as React from 'react';
import PropTypes from "prop-types";
import { useNavigate } from 'react-router-dom';
import { SpeedDial, SpeedDialAction } from "@mui/material";
import PhotoGrid from '../Photos/PhotoGrid';
import PhotoDocs from '../Docs/PhotoDocs';
import UploadIcon from '@mui/icons-material/Upload';
import StarIcon from '@mui/icons-material/Star';
import SpeedDialIcon from '@mui/material/SpeedDialIcon';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoList = ({ user }) => {
  const navigate = useNavigate();
  const [empty, setEmpty] = React.useState(false)

  const navigateToUpload = () => {
    navigate("/upload")
  }

  const navigateToRatings = () => {
    navigate("/ratings")
  }

  const actions = [
    { icon: <UploadIcon />, name: 'Upload photo', action: navigateToUpload },
    { icon: <StarIcon />, name: 'Rate photos', action: navigateToRatings },
  ];


  const populate = () => {
    if (!user) {
      return new Promise((resolve, reject) => {
        reject("User not set!")
      })
    }
    return fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  const handleLoaded = (data) => {
    if (!data || data.length === 0) {
      setEmpty(true)
    }
  }

  return (
    <>
      {empty && <PhotoDocs />}
      <PhotoGrid populator={populate} onDataLoaded={handleLoaded} />
      <SpeedDial
        ariaLabel="Photo actions"
        sx={{ position: 'fixed', bottom: 16, right: 16 }}
        icon={<SpeedDialIcon />}
      >
        {actions.map((action) => (
          <SpeedDialAction
            key={action.name}
            icon={action.icon}
            tooltipTitle={action.name}
            onClick={action.action}
          />
        ))}
      </SpeedDial>
    </>
  )
}

PhotoList.propTypes = {
  user: PropTypes.object
};

export default PhotoList;
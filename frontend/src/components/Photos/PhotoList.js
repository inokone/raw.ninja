

import * as React from 'react';
import PropTypes from "prop-types";
import { useNavigate } from 'react-router-dom';
import PhotoGrid from '../Photos/PhotoGrid';
import PhotoDocs from '../Docs/PhotoDocs';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoList = ({ user }) => {
  const navigate = useNavigate();
  const [empty, setEmpty] = React.useState(false)

  const onFabClick = () => {
    navigate("/upload")
  }

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
      <PhotoGrid populator={populate} onDataLoaded={handleLoaded} fabAction={onFabClick}></PhotoGrid>
    </>
  )
}

PhotoList.propTypes = {
  user: PropTypes.object.isRequired
};

export default PhotoList;
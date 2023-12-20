

import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoList = ({user}) => {
  const populate = () => {
    if(!user) {
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

  return (
    <PhotoGrid populator={populate} data={[]}></PhotoGrid>
  )
}

export default PhotoList;
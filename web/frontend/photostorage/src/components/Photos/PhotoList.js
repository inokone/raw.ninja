

import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';


const { REACT_APP_API_PREFIX } = process.env;

const SearchResult = (props) => {

  const populate = () => {
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

export default SearchResult;
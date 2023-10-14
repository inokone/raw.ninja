import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';


const { REACT_APP_API_PREFIX } = process.env;

const SearchResult = (props) => {

  const populate = () => {
    let url = REACT_APP_API_PREFIX + '/api/v1/search?query=' + props.query;
    return fetch(url, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  return (
    <PhotoGrid populator={populate} data={props.query}></PhotoGrid>
  )
}

export default SearchResult;

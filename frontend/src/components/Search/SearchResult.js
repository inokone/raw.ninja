import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation } from 'react-router-dom'; 

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const SearchResult = () => {
  const { state } = useLocation()

  const populate = () => {
    let url = REACT_APP_API_PREFIX + '/api/v1/search?query=' + state.query;
    return fetch(url, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  return (
    <PhotoGrid key={state.query} populator={populate}/>
  )
}

export default SearchResult;

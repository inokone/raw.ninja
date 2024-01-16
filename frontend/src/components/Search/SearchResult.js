import * as React from 'react';
import PropTypes from "prop-types";
import PhotoGrid from '../Photos/PhotoGrid';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const SearchResult = ({ query }) => {

  const populate = () => {
    let url = REACT_APP_API_PREFIX + '/api/v1/search?query=' + query;
    return fetch(url, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  return (
    <PhotoGrid populator={populate}></PhotoGrid>
  )
}

SearchResult.propTypes = {
  query: PropTypes.string.isRequired,
};

export default SearchResult;

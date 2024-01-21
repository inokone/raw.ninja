import * as React from 'react';
import { useNavigate } from "react-router-dom";
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation } from 'react-router-dom'; 
import { Typography } from '@mui/material';
import Collections from '../Album/Collections'; 
import { Marker } from "react-mark.js";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const SearchResult = () => {
  const navigate = useNavigate()
  const { state } = useLocation()
  const [ data, setData ] = React.useState(null)

  const populate = () => {
    if (state.query.length < 2) {
      return
    }  
    let url = REACT_APP_API_PREFIX + '/api/v1/search?query=' + state.query;
    return fetch(url, {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  const handleDataLoaded = (e) => {
    setData(e)
  }

  const handleAlbumClick = (id) => {
    navigate("/albums/" + id)
  }

  const handleUploadClick = (id) => {
    navigate("/uploads/" + id)
  }

  return (
    <>
      {state.query.length > 1 && 
        <>
          {data && 
            <>
              <Typography variant='h4' sx={{ marginBottom: 2, marginTop: 2 }}>Search results for '{data.query}'</Typography>
              {data.photos.length > 0 &&
                <Typography align='left' variant='h6' sx={{ marginLeft: 1 }}>Photos</Typography>
              }
            </>
          }
          <Marker mark={state.query}>
            <PhotoGrid key={state.query} populator={populate} onDataLoaded={handleDataLoaded} displayTitle={true}/>
          </Marker>
          {data &&
            <>
              {data.albums.length > 0 &&
                <>
                  <Typography align='left' variant='h6' sx={{ marginLeft: 1 }}>Albums</Typography>
                  <Marker mark={data.query}>
                    <Collections collections={data.albums} onClick={handleAlbumClick}/>
                  </Marker>
                </>
              }
              {data.uploads.length > 0 &&
                <>
                  <Typography align='left' variant='h6' sx={{ marginLeft: 1 }}>Uploads</Typography>
                  <Marker mark={data.query}>
                    <Collections collections={data.uploads} onClick={handleUploadClick}/>
                  </Marker>
                </>
              }
            </>
        }
      </>}
    </>
  )
}

export default SearchResult;

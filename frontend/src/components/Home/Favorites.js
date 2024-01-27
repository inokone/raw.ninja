import * as React from 'react';
import PropTypes from "prop-types";
import PhotoGrid from '../Photos/PhotoGrid';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Favorites = ({ user, onDataLoaded }) => {

    const populate = () => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1/search/favorites', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
    }

    const handleDataLoaded = (images) => {
        if (onDataLoaded) {
            onDataLoaded(images)
        }
    }

    return (
        <React.Fragment>
            <PhotoGrid populator={populate} onDataLoaded={handleDataLoaded}/>
        </React.Fragment>
    )
}

Favorites.propTypes = {
    user: PropTypes.object,
    onDataLoaded: PropTypes.func
};

export default Favorites;
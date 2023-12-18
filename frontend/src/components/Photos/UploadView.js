import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation } from "react-router-dom"

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UploadView = ({user}) => {
    const location = useLocation()
    const path = location.pathname
    console.log(path)

    const populate = () => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
                method: "GET",
                mode: "cors",
                credentials: "include"
            }) 
    }

    return (
        <PhotoGrid populator={populate} data={[]}></PhotoGrid>
    )
}

export default UploadView;
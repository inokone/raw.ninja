import React from 'react';
import { Alert } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';
import { useLocation, useNavigate } from "react-router-dom"

const { REACT_APP_API_PREFIX } = process.env;

const format = "jpg"

const settings = () => {
    return encodeURIComponent(JSON.stringify({
        files: [],
        environment: {
            theme: 1,
            vmode: 1,
            intro: false,
            lang: "en",
            localsave: false,
            phrases: [[1, 2], "Save As " + format.toUpperCase()],
            menus: [[0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0 , 0, 1], 1, 1, 1, 1, 1, 1],
            customIO: { save: "app.activeDocument.saveToOE(\"" + format + "\");" },
        },
        server: {
            version: 1
        }
    }))
}

const setEditorImage = (image) => {
    let pp = document.getElementById("pp");
    if (pp) {
        var wnd = document.getElementById("pp").contentWindow;
        wnd.postMessage(image, "*");
    } else {
        console.log("PhotoPea is not initialized, can not set image")
    }
}

const Photopea = () => {
    const navigate = useNavigate()
    const location = useLocation()
    const [image, setImage] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [initialzing, setInitialzing] = React.useState(false)
    const [saving, setSaving] = React.useState(false)
    const [counter, setCounter] = React.useState(0)

    const saveImage = (data) => {
        setSaving(true)
        const formData = new FormData();
        formData.append(
            'files[]',
            new Blob([data], { type: getMimeTypeFromArrayBuffer(data) }),
            "edited-" + new Date().toISOString().split('T')[0] + "." + format);
        const requestOptions = {
            method: 'POST',
            mode: "cors",
            credentials: "include",
            body: formData,
        };

        const url = REACT_APP_API_PREFIX + '/api/v1/photos/';
        fetch(url, requestOptions)
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Request failed');
                }
            })
            .then(data => {
                console.log(data)
                setSaving(false)
                navigate('/photos/' + data.photo_ids[0])
            })
            .catch(error => {
                console.error('Error:', error);
                setSaving(false)
            });
    }

    function getMimeTypeFromArrayBuffer(arrayBuffer) {
        const uint8arr = new Uint8Array(arrayBuffer)

        const len = 4
        if (uint8arr.length >= len) {
            let signatureArr = new Array(len)
            for (let i = 0; i < len; i++)
                signatureArr[i] = (new Uint8Array(arrayBuffer))[i].toString(16)
            const signature = signatureArr.join('').toUpperCase()

            switch (signature) {
                case '89504E47':
                    return 'image/png'
                case '47494638':
                    return 'image/gif'
                case 'FFD8FFDB':
                case 'FFD8FFE0':
                    return 'image/jpeg'
                default:
                    return null
            }
        }
        return null
    }

    const handleEditorMessage = (e) => {
        if (e.data.source === "react-devtools-content-script" || e.data.source === "react-devtools-bridge") {
            return
        }
        if (e.data === "done") {
            setCounter(counter + 1)
            if (counter === 1) {
                setInitialzing(false)
            } else if (counter === 2) {
                // image loaded, fit to screen
                var wnd = document.getElementById("pp").contentWindow;
                wnd.postMessage("app.UI.fitTheArea()", "*")
                // turn off progress
            }
        } else {
            let size = e.data.byteLength
            if (!saving && size) {
                setSaving(true)
                console.log("Saving file of size " + size) // e.data is an arrayBuffer, we need to save it
                saveImage(e.data)
            }
        }
    }

    const loadImage = () => {
        setCounter(0)
        let id = location.pathname.split('/').slice(-1)
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + id + '/download', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                        setLoading(false)
                    })
                } else {
                    response.blob().then(content => {
                        content.arrayBuffer().then(img => {
                            setLoading(false)
                            setInitialzing(true)
                            setImage(img)
                            setEditorImage(img)
                        })
                    })
                }
            })
    }

    React.useEffect(() => {
        window.addEventListener("message", handleEditorMessage);
        if (!loading && !error && !image) {
            loadImage()
        }
        return () => {
            window.removeEventListener('message', handleEditorMessage);
        }
    }, [loading, error, image])

    return (
        <div className="iframe-container">
            {loading && <ProgressDisplay />}
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {initialzing && <Alert sx={{ mb: 1 }} onClose={() => { setInitialzing(false) }}>Loading image...</Alert>}
            {saving && <Alert sx={{ mb: 1 }} onClose={() => { setSaving(false) }}>Saving modifications...</Alert>}
            <iframe title="Editor" width="100%" id="pp" src={"https://photopea.com#" + settings()}
                frameBorder="no" border="0" scrolling="no" height='800' />
        </div>
    );
}

export default Photopea; 
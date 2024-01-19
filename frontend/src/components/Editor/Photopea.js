import React from 'react';
import { Alert } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';
import { useLocation, useNavigate } from "react-router-dom"

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const format = "jpg"

const settings = (imageURL) => {
    return encodeURIComponent(JSON.stringify({
        files: [imageURL],
        environment: {
            theme: 1,
            vmode: 1,
            intro: false,
            lang: "en",
            localsave: false,
            phrases: [[1, 2], "Save As " + format.toUpperCase()],
            menus: [[0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1], 1, 1, 1, 1, 1, 1],
            customIO: { save: "app.activeDocument.saveToOE(\"" + format + "\");" },
        },
        server: {
            version: 1
        }
    }))
}

const Photopea = () => {
    const { state } = useLocation();
    const navigate = useNavigate()
    const [image, setImage] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [saving, setSaving] = React.useState(false)
    const [counter, setCounter] = React.useState(0)

    const saveImage = React.useCallback((data) => {
        setSaving(true)
        const formData = new FormData();
        formData.append(
            'files[]',
            new Blob([data], { type: getMimeTypeFromArrayBuffer(data) }),
            state.photo_name + "-edited." + format);
        const requestOptions = {
            method: 'POST',
            mode: "cors",
            credentials: "include",
            body: formData,
        };

        console.log(formData)
        const url = REACT_APP_API_PREFIX + '/api/v1/uploads/';
        fetch(url, requestOptions)
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Request failed');
                }
            })
            .then(data => {
                setSaving(false)
                navigate('/uploads/' + data)
            })
            .catch(error => {
                console.error('Error:', error);
                setSaving(false)
            });
    }, [navigate, state])

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

    const handleEditorMessage = React.useCallback((e) => {
        if (["react-devtools-content-script", "react-devtools-bridge", "react-devtools-backend-manager"].indexOf(e.data.source) >= 0) {
            return
        }
        if (e.data === "done") {
            let cnt = counter + 1
            setCounter(cnt)
            if (cnt === 2) {
                // image loaded, fit to screen
                var wnd = document.getElementById("pp").contentWindow;
                wnd.postMessage("app.UI.fitTheArea()", "*")
                // turn off progress
            }
        } else {
            let size = e.data.byteLength
            if (!saving && size) {
                setSaving(true)
                saveImage(e.data)
            }
        }
    }, [counter, saveImage, saving])

    const loadImage = React.useCallback(() => {
        setCounter(0)
        setLoading(true)
        fetch(REACT_APP_API_PREFIX + '/api/v1/onetime/', {
            method: "POST",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify({
                original_id: state.photo_id,
                one_time: false
            })
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                        setLoading(false)
                    })
                } else {
                    response.json().then(content => {
                        setImage(REACT_APP_API_PREFIX + '/api/public/v1/onetime/raw/' + content.id)
                        setLoading(false)
                    })
                }
            })
    }, [state])

    React.useEffect(() => {
        window.addEventListener("message", handleEditorMessage);
        if (!loading && !error && !image) {
            loadImage()
        }
        return () => {
            window.removeEventListener('message', handleEditorMessage);
        }
    }, [loading, error, image, handleEditorMessage, loadImage])

    return (
        <div className="iframe-container">
            {loading && <ProgressDisplay />}
            {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
            {saving && <Alert sx={{ mb: 1 }} onClose={() => { setSaving(false) }}>Saving modifications...</Alert>}
            <iframe title="Editor" width="100%" id="pp" src={"https://photopea.com#" + settings(image)}
                frameBorder="no" border="0" scrolling="no" height='800' />
        </div>
    );
}

export default Photopea; 
import React from 'react';
import { Alert } from '@mui/material';
import ProgressDisplay from '../Common/ProgressDisplay';
import { useLocation, useSearchParams } from "react-router-dom"

const { REACT_APP_API_PREFIX } = process.env;

const settings = (format) => {
    return encodeURIComponent(JSON.stringify({
        files: [],
        environment: {
            theme: 1,
            vmode: 1,
            intro: false,
            lang: "en",
            localsave: false,
            phrases: [[1, 2], "Save As PNG"],
            menus: [[0,0,0,0,0,1], 1, 1, 1, 1, 1, 1],
            customIO: { save: "app.activeDocument.saveToOE(\"png\");" },
        },
        server: {
            version: 1
        }
    }))
}

const setEditorImage = (image) => {
    var wnd = document.getElementById("pp").contentWindow;
    wnd.postMessage(image, "*");
}

const Photopea = () => {
    const location = useLocation()
    const [searchParams, setSearchParams] = useSearchParams();
    const [image, setImage] = React.useState(null)
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState(null)
    const [initialzed, setInitialized] = React.useState(false)
    const [saving, setSaving] = React.useState(false)

    React.useEffect(() => {
        let ppDoneCounter = 0

        const handleEditorMessage = (e) => {
            if (e.data.source === "react-devtools-content-script" || e.data.source === "react-devtools-bridge") {
                return
            }
            if (e.data === "done") {
                setInitialized(true)
            } else if (e.data === "done") {
                ppDoneCounter++
                if (ppDoneCounter === 1) {
                    console.log("PP: loaded")
                    // pp loaded
                } else if (ppDoneCounter === 2) {
                    // image loaded, fit to screen
                    console.log("PP: fitting to area")
                    var wnd = document.getElementById("pp").contentWindow;
                    wnd.postMessage("app.UI.fitTheArea()", "*")
                    // turn off progress
                } else {
                    console.log("PP: done " + ppDoneCounter)
                    // save image finished
                    setSaving(false)
                }
            } else {
                let size = e.data.byteLength
                if (!saving && size) {
                    console.log("Saving file of size " + size)
                    setSaving(true)
                }
            }
        }

        const getImage = () => {
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
                                setImage(img)
                                setEditorImage(img)
                            })
                        })
                    }
                })
        }

        window.addEventListener("message", handleEditorMessage);
        if(!loading && !error && !image){
            getImage()
        }
    }, [loading, error, image, initialzed, location.pathname])

    return (
        <div className="iframe-container">
            {loading && <ProgressDisplay />}
            {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
            {saving && <Alert sx={{ mb: 4 }}>Saving modifications...</Alert>}
            <iframe title="Editor" width="100%" id="pp" src={"https://photopea.com#" + settings(searchParams.get("format"))}
                    frameBorder="no" border="0" scrolling="no" height='800'/>
        </div>
    );
}

export default Photopea; 
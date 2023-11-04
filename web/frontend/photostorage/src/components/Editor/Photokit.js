import React from 'react';

function Photokit({photoB64, setPhotoB64}) {

    React.useEffect(() => {
        window.addEventListener("message", function (event) {
            if (event.data.type === "photokitsdk" && event.data.name === "editorLoaded") {
                this.window.photokit.contentWindow.postMessage({ type: 'photokitsdk', name: 'openimage', data: photoB64, opentype: 0 }, '*');
            }
        });

        window.addEventListener("message", function (event) {
            if (event.data.type === "photokitsdk" && event.data.name === "saveimage") {
                setPhotoB64(event.data.imagedata);
                // setName(event.data.imagename);
            }
        });
    }, [photoB64, setPhotoB64]);


    return (
        <div>
            <iframe id="photokit" title="Editor" width="100%" src="https://photokit.com/editor/?lang=en" frameborder="no" border="0" scrolling="no" 
                height='800' />
        </div>
    );
}


export default Photokit;
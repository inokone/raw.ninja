const dateOf = (data) => {
    return new Date(data).toLocaleDateString()
}

const formatShutterSpeed = (shutterSpeed) => {
    let validDividers = [2, 4, 8, 15, 30, 60, 125, 250, 500, 1000, 2000, 4000, 8000]
    if (shutterSpeed < 1) {
        let fraction = 1 / shutterSpeed
        let lastDivider = 2
        for (let i = 0; i < validDividers.length; i++) {
            let divider = validDividers[i]
            if (fraction < divider) {
                if (fraction - lastDivider > divider - fraction) {
                    return "1/" + divider
                } else {
                    return "1/" + lastDivider
                }
            }
        }
    } else {
        return shutterSpeed.toFixed(1) + "s";
    }
}

const formatBytes = (bytes, decimals = 2) => {
    let negative = (bytes < 0)
    if (negative) {
        bytes = -bytes
    }
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    let displayText = `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    return negative ? "-" + displayText : displayText
}

const validTime = (time) => {
    return !time.startsWith("1970-01-01T")
}

const description = (photo) => {
    let date = dateOf(photo.descriptor.metadata.timestamp)
    let timestamp = validTime(photo.descriptor.metadata.timestamp) ? ("Taken on " + date + "\n") : ""
    let aperture = photo.descriptor.metadata.aperture !== 0 ? "ƒ/" + Math.round((photo.descriptor.metadata.aperture + Number.EPSILON) * 100) / 100 + "  " : ""
    let shutter = photo.descriptor.metadata.shutter !== 0 ? formatShutterSpeed(photo.descriptor.metadata.shutter) + " sec  " : ""
    let iso = photo.descriptor.metadata.ISO !== 0 ? "ISO " + photo.descriptor.metadata.ISO + "  " : ""
    let dim = photo.descriptor.metadata.width + " x " + photo.descriptor.metadata.height + " px  "
    let size = formatBytes(photo.descriptor.metadata.data_size)
    return timestamp + aperture + shutter + iso + dim + size
}

export function convertPhoto(photo) {
    return {
        src: photo.thumbnail.url,
        original: photo.thumbnail.url,
        width: photo.descriptor.thumbnail_width ? photo.descriptor.thumbnail_width : photo.descriptor.metadata.width,
        height: photo.descriptor.thumbnail_height ? photo.descriptor.thumbnail_height : photo.descriptor.metadata.height,
        caption: photo.descriptor.filename,
        title: photo.descriptor.filename,
        description: description(photo),
        favorite: photo.descriptor.favorite,
        rating: photo.descriptor.rating,
        id: photo.id,
        format: photo.descriptor.format,
        base: photo,
        selected: false
    }
}

export function convertPhotos(content) {
    if (!Array.isArray(content)) {
        content = content.photos
    }
    return content.map(photo => convertPhoto(photo))
}

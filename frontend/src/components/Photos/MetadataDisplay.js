import React from 'react';
import Typography from '@mui/material/Typography';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import { createStyles, makeStyles } from '@mui/styles';
import { Grid, Box } from "@mui/material";



const useStyles = makeStyles(theme => createStyles({
  tableDataCell: {
    color: theme.palette.common.bluegray,
    align: 'left'
  },
  tableGroup: {
    color: theme.palette.common.white,
  }
}))

const MetadataDisplay = ({ metadata }) => {
  const classes = useStyles();

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

  const hasCameraData = () => {
    return metadata.metadata.camera_make || metadata.metadata.camera_model || metadata.metadata.camera_sw
  }

  const hasLensData = () => {
    return metadata.metadata.lens_make || metadata.metadata.lens_model
  }

  const generalData = {
    title: "General",
    values: [
      {
        label: "Format", 
        value: metadata.format},
      {
        label: "Uploaded", 
        value: new Date(metadata.uploaded).toLocaleString()},
      {
        label: "Taken",
        value: metadata.metadata.timestamp === 0 ? "N/A" : new Date(metadata.metadata.timestamp).toLocaleString()},
      {
        label: "Size", 
        value: formatBytes(metadata.metadata.data_size)
      }
    ]
  }

  const imageData = {
    title: "Image",
    values: [
      {
        label: "Width",
        value: metadata.metadata.width + " px"
      },
      {
        label: "Height",
        value: metadata.metadata.height + " px"
      },
      {
        label: "ISO",
        value: metadata.metadata.ISO !== 0 ? metadata.metadata.ISO : "N/A"
      },
      {
        label: "Aperture",
        value: metadata.metadata.aperture !== 0 ? "ƒ/" + Math.round((metadata.metadata.aperture + Number.EPSILON) * 100) / 100 : "N/A"
      },
      {
        label: "Shutter Speed",
        value: metadata.metadata.shutter !== 0 ? formatShutterSpeed(metadata.metadata.shutter) : "N/A"
      }
    ]
  }

  const cameraData = {
    title: "Camera",
    values: [
      {
        label: "Manufacturer",
        value: metadata.metadata.camera_make
      },
      {
        label: "Model",
        value: metadata.metadata.camera_model
      },
      {
        label: "Software Version",
        value: metadata.metadata.camera_sw
      }
    ]
  }

  const lensData = {
    title: "Lens",
    values: [
      {
        label: "Manufacturer",
        value: metadata.metadata.lens_make
      },
      {
        label: "Model",
        value: metadata.metadata.lens_model
      }
    ]
  }

  const displayDataOld = (data) => {
    return (
      <React.Fragment>
        <Typography variant="h5" className={classes.tableGroup}>{data.title}</Typography>
        <Table size="small">
          <TableBody>
            {
              data.values.map((value, index) => {
                return (
                  <TableRow key={data.title + index}>
                    <TableCell className={classes.tableDataCell}>{value.label}</TableCell>
                    <TableCell className={classes.tableDataCell}>{value.value}</TableCell>
                  </TableRow>
                )
              })
            }
          </TableBody>
        </Table>
      </React.Fragment>
    )
  }

  const displayData = (data) => {
    return (
      <React.Fragment>
        <Box sx={{ display: 'flex', justifyContent: 'center', borderRadius: '4px', pb: 1 }}>
          <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.28)', color: 'white', mb: 1, borderRadius: '4px', width: '500px' }}>
            <Grid container>
              <Grid item xs={12}><Typography variant='h6' sx={{ borderRadius: '4px', bgcolor: 'rgba(0, 0, 0, 0.34)' }}>{data.title}</Typography></Grid>
            </Grid>
            <Grid container>
              {data.values.map((value, index) => {
                return (
                  <React.Fragment key={data.title + index}>
                    <Grid item xs={5} textAlign={'left'} pl='5px' pb='8px' ><Typography>{value.label}</Typography></Grid>
                    <Grid item xs={7}><Typography>{value.value}</Typography></Grid>
                  </React.Fragment>
                )
              })}
            </Grid>
          </Box>
        </Box>
        </React.Fragment>
    )
  }

  return (
    <React.Fragment>
      {displayData(generalData)}
      {displayData(imageData)}
      {hasCameraData() && displayData(cameraData)}
      {hasLensData() && displayData(lensData)}
    </React.Fragment>
  );
}
export default MetadataDisplay; 
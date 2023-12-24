import { Doughnut } from "react-chartjs-2";
import { Chart, ArcElement, Title, Tooltip } from "chart.js";
import { Box, Typography, Stack } from '@mui/material';
import { makeStyles, useTheme } from '@mui/styles';

const useStyles = makeStyles((theme) => ({
    percentage: {
        color: theme.palette.primary.main,
        fontWeight: 1000,
        fontSize: '3em'
    },
    percentageText: {
        color: theme.palette.primary.main,
        fontWeight: 800,
        fontSize: '2em',
        marginLeft: 5
    },
}));

Chart.register(ArcElement, Title, Tooltip);

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

const data = (usedSpace, quota, theme) => {
    return {
        labels: [
            formatBytes(usedSpace),
            quota <= 0 ? 'Unlimited' : formatBytes(quota)
        ],
        datasets: [
            {
                data: [
                    usedSpace,  
                    quota <= 0 ? 9 * usedSpace : quota - usedSpace],
                backgroundColor: [
                    theme.palette.primary.light,
                    theme.palette.secondary.light,
                ],
                borderColor: "#D1D6DC"
            }
        ]
    }
};

const options = {
    plugins: {
        title: {
            display: true,
            text: 'Storage',
        },
        tooltip: {
            callbacks: {
                label: function (context) {
                    return context.dataset.label || '';
                }
            }
        }
    },
    rotation: 270,
    circumference: 180,
    cutout: "80%",
    responsive: true,
    maintainAspectRatio: true
}

const SpaceChart = ({usedSpace, quota}) => {
    const classes = useStyles();
    const theme = useTheme();

    return (
        <Box sx={{position: 'relative'}}>
            <Doughnut
                data={data(usedSpace, quota, theme)}
                options={options}
            />
            <Stack
                direction={"row"}
                alignItems={"baseline"}
                style={{
                    position: "absolute",
                    top: "70%",
                    left: "50%",
                    transform: "translate(-50%, -50%)",
                    textAlign: "center"
                }}
            >
                <Typography className={classes.percentage}>{quota > 0 ? Math.round(100 * usedSpace / quota) + "%" : "Unlimited"}</Typography>
                <Typography className={classes.percentageText}>used</Typography>
            </Stack>
        </Box>
    );
};

export default SpaceChart;

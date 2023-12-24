import { Doughnut } from "react-chartjs-2";
import { Chart, ArcElement, Title, Tooltip } from "chart.js";
import { useTheme } from '@mui/styles';

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
    aspectRatio: 2,
    responsive: true,
    maintainAspectRatio: true
}

const textCenter = (text, color) => {
    const t = text
    return {
        id: 'textCenter',
        beforeDatasetsDraw(chart) {
            const { ctx, data, width, height } = chart;
            ctx.save();
            const fontsize = (height / 160).toFixed(2);
            ctx.font = `bolder ${fontsize}em Poppins`;
            ctx.fillStyle = color
            ctx.textBaseLine = 'middle';
            ctx.textAlign = 'center'
            ctx.fillText(t, chart.getDatasetMeta(0).data[0].x, chart.getDatasetMeta(0).data[0].y);
            ctx.save();
        }
    }
}

const SpaceChart = ({usedSpace, quota}) => {
    const theme = useTheme();

    return (
        <Doughnut 
            data={data(usedSpace, quota, theme)}
            options={options}
            plugins={[textCenter(
                quota > 0 ? Math.round(100 * usedSpace / quota) + "% used" : "Unlimited",
                theme.palette.primary.light)]}
        />
    );
};

export default SpaceChart;

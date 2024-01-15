import PropTypes from "prop-types";
import { Line } from "react-chartjs-2";
import { Chart, CategoryScale, LinearScale, LineElement, Title, Tooltip, PointElement } from "chart.js";
import { useTheme } from '@mui/styles';

Chart.register(CategoryScale, LinearScale, LineElement, Title, Tooltip, PointElement);

const data = (uploads) => {
    return {
        labels: Object.keys(uploads).map(timestamp => dateOf(timestamp)),
        datasets: [
            {
                data: Object.values(uploads),
            }
        ]
    }
};

const options = (theme) => {
    return {
        tension: 0.2,
        plugins: {
            title: {
                display: true,
                text: 'Upload history',
            },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        return context.parsed.y + (context.parsed.y !== 1 ? " photos" : " photo");
                    }
                }
            }
        },
        scales: {
            x: {
                ticks: false
            }
        },
        backgroundColor: theme.palette.secondary.main,
        border: {
            color: '#673AB7',
            width: 5,
            capStyle: 'butt',
            dash: [],
            dashOffset: 0.1,
            joinStyle: 'miter',
        },
        point: {
            borderColor: '#673AB7',
            backgroundColor: '#fff',
            borderWidth: 2,
            hoverRadius: 5,
            hoverBackgroundColor: '#673AB7',
            hoverBorderColor: 'rgba(220,220,220,1)',
            hoverBorderWidth: 2,
            radius: 1,
            hitRadius: 10,
        },
        borderColor: theme.palette.primary.light,
        responsive: true
    }
};

const dateOf = (data) => {
    return new Date(data).toLocaleDateString()
};

const UploadChart = ({ uploads }) => {
    const theme = useTheme();
    return (
        <Line
            data={data(uploads)}
            options={options(theme)}
        />
    );
};

UploadChart.propTypes = {
    uploads: PropTypes.array
};

export default UploadChart;

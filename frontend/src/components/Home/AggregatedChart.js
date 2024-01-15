import { Bar } from "react-chartjs-2";
import { Chart, CategoryScale, LinearScale, BarElement, Title, Tooltip } from "chart.js";
import { useTheme } from '@mui/styles';
import PropTypes from "prop-types";

Chart.register(CategoryScale, LinearScale, BarElement, Title, Tooltip);

const data = (photos, favorites, albums, theme) => {
    return {
        labels: [
            "Photos",
            "Favorites",
            "Albums"
        ],
        datasets: [
            {
                data: [
                    photos,
                    favorites,
                    albums
                ],
                backgroundColor: [
                    theme.palette.primary.light,
                    theme.palette.primary.main,
                    theme.palette.primary.dark,
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
            text: 'Stored entities',
        }
    },
    scales: {
        y: {
            ticks: {
                display: false,
            },
        }
    },
    responsive: true
};

const AggregatedChart = ({ photos, favorites, albums }) => {
    const theme = useTheme();
    return (
        <Bar
            data={data(photos, favorites, albums, theme)}
            options={options}
        />
    );
};

AggregatedChart.propTypes = {
    photos: PropTypes.number,
    favorites: PropTypes.number,
    albums: PropTypes.number
};

export default AggregatedChart;

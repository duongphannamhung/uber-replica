import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar } from 'react-chartjs-2';
import TitleCard from '../../../components/Cards/TitleCard';
import axios from 'axios'
import { useState, useEffect } from 'react'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

function getMonthNames(monthNumbers) {
  const months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
  return monthNumbers.map(number => months[number - 1]);
}

function BarChartCustom(){
    const [month_arr, set_month] = useState([]);
    const [revenue_arr, set_revenue] = useState([]);
    const [profit_arr, set_profit] = useState([]);

    const fetchRevenueYear = () => {
      axios.get('/api/dashboard/calculate-revenue-year')
      .then(response => {
          const data = response.data;
          set_month(prevStats => data.month);
          set_revenue(prevStats => data.revenue);
          let _profit_arr = [];
          for (let i = 0; i < data.month.length; i++) {
            _profit_arr.push(data.revenue[i] * Math.random());
          }
          set_profit(prevStats => _profit_arr);
      })
      .catch(error => console.error('Error:', error));
    }

    useEffect(() => {
      set_month([]);
      set_revenue([]);
      set_profit([]);
      fetchRevenueYear();
    }, []);


    const options = {
        responsive: true,
        plugins: {
          legend: {
            position: 'top',
          }
        },
      };

      const labels = getMonthNames(month_arr);

      const data = {
        labels,
        datasets: [
          {
            label: 'profit',
            data: profit_arr,
            backgroundColor: 'rgba(255, 99, 132, 1)',
          },
          {
            label: 'revenue',
            data: revenue_arr,
            backgroundColor: 'rgba(53, 162, 235, 1)',
          },
        ],
      };

    return(
      <TitleCard title={"Revenue (in VND)"}>
            <Bar options={options} data={data} />
      </TitleCard>

    )
}


export default BarChartCustom